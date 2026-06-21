package intercambio

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type CreateSwapRequestCommand struct {
	PlanificacionID    string
	TurnoSolicitanteID string
	TurnoDestinoID     string
	SolicitanteID      string
	DestinoID          string
	ActorID            string
}

type CreateSwapRequestHandler struct {
	repo      ports.ShiftSwapRequestRepository
	turnoRepo ports.TurnoRepository
}

func NewCreateSwapRequestHandler(repo ports.ShiftSwapRequestRepository, turnoRepo ports.TurnoRepository) *CreateSwapRequestHandler {
	return &CreateSwapRequestHandler{repo: repo, turnoRepo: turnoRepo}
}

func (h *CreateSwapRequestHandler) Handle(ctx context.Context, cmd CreateSwapRequestCommand) (*intercambio.ShiftSwapRequest, error) {
	turnoSol, err := h.turnoRepo.FindByTurnoID(ctx, cmd.TurnoSolicitanteID)
	if err != nil {
		return nil, fmt.Errorf("turno solicitante not found: %w", err)
	}
	turnoDest, err := h.turnoRepo.FindByTurnoID(ctx, cmd.TurnoDestinoID)
	if err != nil {
		return nil, fmt.Errorf("turno destino not found: %w", err)
	}

	if turnoSol.PlanificacionID != cmd.PlanificacionID || turnoDest.PlanificacionID != cmd.PlanificacionID {
		return nil, fmt.Errorf("both shifts must belong to the same planificacion")
	}
	if turnoSol.EmpleadoID != cmd.SolicitanteID {
		return nil, fmt.Errorf("turno solicitante does not belong to the requesting employee")
	}
	if turnoDest.EmpleadoID != cmd.DestinoID {
		return nil, fmt.Errorf("turno destino does not belong to the target employee")
	}

	req, err := intercambio.NewShiftSwapRequest(intercambio.NewSwapRequestParams{
		PlanificacionID:    cmd.PlanificacionID,
		TurnoSolicitanteID: cmd.TurnoSolicitanteID,
		TurnoDestinoID:     cmd.TurnoDestinoID,
		SolicitanteID:      cmd.SolicitanteID,
		DestinoID:          cmd.DestinoID,
	})
	if err != nil {
		return nil, err
	}

	if err := h.repo.Create(ctx, req); err != nil {
		return nil, err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionSolicitado, cmd.ActorID, nil)
	if err := h.repo.AddHistoryEntry(ctx, entry); err != nil {
		return nil, err
	}

	return req, nil
}
