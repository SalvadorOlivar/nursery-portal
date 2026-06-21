package leave

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/ports"
)

type ApproveLeaveRequestCommand struct {
	ID          string
	ApprovedBy  string
}

type ApproveLeaveRequestHandler struct {
	repo ports.LeaveRequestRepository
}

func NewApproveLeaveRequestHandler(repo ports.LeaveRequestRepository) *ApproveLeaveRequestHandler {
	return &ApproveLeaveRequestHandler{repo: repo}
}

func (h *ApproveLeaveRequestHandler) Handle(ctx context.Context, cmd ApproveLeaveRequestCommand) error {
	lr, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("leave request not found: %w", err)
	}

	if err := lr.Approve(cmd.ApprovedBy); err != nil {
		return err
	}

	return h.repo.Update(ctx, lr)
}
