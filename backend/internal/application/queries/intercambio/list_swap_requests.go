package intercambio

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type ListSwapRequestsQuery struct {
	EmployeeID string
	Role       string
}

type ListSwapRequestsHandler struct {
	repo ports.ShiftSwapRequestRepository
}

func NewListSwapRequestsHandler(repo ports.ShiftSwapRequestRepository) *ListSwapRequestsHandler {
	return &ListSwapRequestsHandler{repo: repo}
}

func (h *ListSwapRequestsHandler) Handle(ctx context.Context, q ListSwapRequestsQuery) ([]*intercambio.ShiftSwapRequest, error) {
	all, err := h.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if q.Role == "ADMIN" || q.Role == "SUPERVISOR" {
		return all, nil
	}

	if q.EmployeeID == "" {
		return all, nil
	}

	var filtered []*intercambio.ShiftSwapRequest
	for _, req := range all {
		if req.SolicitanteID == q.EmployeeID || req.DestinoID == q.EmployeeID {
			filtered = append(filtered, req)
		}
	}
	return filtered, nil
}

type GetSwapRequestHistoryHandler struct {
	repo ports.ShiftSwapRequestRepository
}

func NewGetSwapRequestHistoryHandler(repo ports.ShiftSwapRequestRepository) *GetSwapRequestHistoryHandler {
	return &GetSwapRequestHistoryHandler{repo: repo}
}

func (h *GetSwapRequestHistoryHandler) Handle(ctx context.Context, swapRequestID string) ([]*intercambio.ShiftSwapHistoryEntry, error) {
	return h.repo.GetHistory(ctx, swapRequestID)
}
