package intercambio

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type ApproveSwapRequestCommand struct {
	ID       string
	ActorID  string
}

type ApproveSwapRequestHandler struct {
	repo      ports.ShiftSwapRequestRepository
	turnoRepo ports.TurnoRepository
}

func NewApproveSwapRequestHandler(repo ports.ShiftSwapRequestRepository, turnoRepo ports.TurnoRepository) *ApproveSwapRequestHandler {
	return &ApproveSwapRequestHandler{repo: repo, turnoRepo: turnoRepo}
}

func (h *ApproveSwapRequestHandler) Handle(ctx context.Context, cmd ApproveSwapRequestCommand) error {
	req, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("swap request not found: %w", err)
	}

	if err := req.Approve(cmd.ActorID); err != nil {
		return err
	}

	turnoSol, err := h.turnoRepo.FindByTurnoID(ctx, req.TurnoSolicitanteID)
	if err != nil {
		return fmt.Errorf("turno solicitante not found: %w", err)
	}

	turnoDest, err := h.turnoRepo.FindByTurnoID(ctx, req.TurnoDestinoID)
	if err != nil {
		return fmt.Errorf("turno destino not found: %w", err)
	}

	solicitanteEmpleadoID := turnoSol.EmpleadoID
	destinoEmpleadoID := turnoDest.EmpleadoID

	turnoSol.EmpleadoID = destinoEmpleadoID
	turnoDest.EmpleadoID = solicitanteEmpleadoID

	if err := h.turnoRepo.Update(ctx, turnoSol); err != nil {
		return fmt.Errorf("failed to update turno solicitante: %w", err)
	}
	if err := h.turnoRepo.Update(ctx, turnoDest); err != nil {
		return fmt.Errorf("failed to update turno destino: %w", err)
	}

	if err := h.repo.Update(ctx, req); err != nil {
		return err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionEjecutado, cmd.ActorID, nil)
	if err := h.repo.AddHistoryEntry(ctx, entry); err != nil {
		return err
	}

	return nil
}
