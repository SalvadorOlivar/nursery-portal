package leave

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/ports"
)

type RejectLeaveRequestCommand struct {
	ID          string
	ApprovedBy  string
}

type RejectLeaveRequestHandler struct {
	repo ports.LeaveRequestRepository
}

func NewRejectLeaveRequestHandler(repo ports.LeaveRequestRepository) *RejectLeaveRequestHandler {
	return &RejectLeaveRequestHandler{repo: repo}
}

func (h *RejectLeaveRequestHandler) Handle(ctx context.Context, cmd RejectLeaveRequestCommand) error {
	lr, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("leave request not found: %w", err)
	}

	if err := lr.Reject(cmd.ApprovedBy); err != nil {
		return err
	}

	return h.repo.Update(ctx, lr)
}
