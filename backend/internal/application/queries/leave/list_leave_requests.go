package leave

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type ListLeaveRequestsQuery struct {
	EmployeeID string
}

type ListLeaveRequestsHandler struct {
	repo ports.LeaveRequestRepository
}

func NewListLeaveRequestsHandler(repo ports.LeaveRequestRepository) *ListLeaveRequestsHandler {
	return &ListLeaveRequestsHandler{repo: repo}
}

func (h *ListLeaveRequestsHandler) Handle(ctx context.Context, q ListLeaveRequestsQuery) ([]*ausencia.LeaveRequest, error) {
	if q.EmployeeID != "" {
		return h.repo.FindByEmployee(ctx, q.EmployeeID)
	}
	return h.repo.FindAll(ctx)
}
