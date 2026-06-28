package intercambio

import (
	"context"
	"fmt"

	"github.com/tuusuario/nurse-portal/internal/domain/intercambio"
	"github.com/tuusuario/nurse-portal/internal/ports"
)

type AcceptSwapRequestCommand struct {
	ID         string
	ActorID    string
	EmployeeID string
}

type RejectSwapRequestCommand struct {
	ID         string
	ActorID    string
	EmployeeID string
}

type AcceptSwapRequestHandler struct {
	repo ports.ShiftSwapRequestRepository
}

func NewAcceptSwapRequestHandler(repo ports.ShiftSwapRequestRepository) *AcceptSwapRequestHandler {
	return &AcceptSwapRequestHandler{repo: repo}
}

func (h *AcceptSwapRequestHandler) Handle(ctx context.Context, cmd AcceptSwapRequestCommand) error {
	req, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("solicitud de intercambio no encontrada: %w", err)
	}

	if req.DestinoID != cmd.EmployeeID {
		return fmt.Errorf("solo el empleado destino puede aceptar esta solicitud")
	}

	if err := req.AcceptByDestino(); err != nil {
		return err
	}

	if err := h.repo.Update(ctx, req); err != nil {
		return err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionAceptado, cmd.ActorID, nil)
	return h.repo.AddHistoryEntry(ctx, entry)
}

type RejectSwapRequestHandler struct {
	repo ports.ShiftSwapRequestRepository
}

func NewRejectSwapRequestHandler(repo ports.ShiftSwapRequestRepository) *RejectSwapRequestHandler {
	return &RejectSwapRequestHandler{repo: repo}
}

func (h *RejectSwapRequestHandler) Handle(ctx context.Context, cmd RejectSwapRequestCommand) error {
	req, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("solicitud de intercambio no encontrada: %w", err)
	}

	if req.DestinoID != cmd.EmployeeID {
		return fmt.Errorf("solo el empleado destino puede rechazar esta solicitud")
	}

	if err := req.Reject(); err != nil {
		return err
	}

	if err := h.repo.Update(ctx, req); err != nil {
		return err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionRechazado, cmd.ActorID, nil)
	return h.repo.AddHistoryEntry(ctx, entry)
}
