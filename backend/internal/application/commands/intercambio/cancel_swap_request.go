package intercambio

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type CancelSwapRequestCommand struct {
	ID      string
	ActorID string
}

type CancelSwapRequestHandler struct {
	repo ports.ShiftSwapRequestRepository
}

func NewCancelSwapRequestHandler(repo ports.ShiftSwapRequestRepository) *CancelSwapRequestHandler {
	return &CancelSwapRequestHandler{repo: repo}
}

func (h *CancelSwapRequestHandler) Handle(ctx context.Context, cmd CancelSwapRequestCommand) error {
	req, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("swap request not found: %w", err)
	}

	if req.SolicitanteID != cmd.ActorID {
		return fmt.Errorf("only the requesting employee can cancel this request")
	}

	if err := req.Cancel(); err != nil {
		return err
	}

	if err := h.repo.Update(ctx, req); err != nil {
		return err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionCancelado, cmd.ActorID, nil)
	return h.repo.AddHistoryEntry(ctx, entry)
}
