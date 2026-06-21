package compensatory

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type ListCompensatoryDaysQuery struct {
	EmployeeID string
}

type ListCompensatoryDaysHandler struct {
	repo ports.CompensatoryDayRepository
}

func NewListCompensatoryDaysHandler(repo ports.CompensatoryDayRepository) *ListCompensatoryDaysHandler {
	return &ListCompensatoryDaysHandler{repo: repo}
}

func (h *ListCompensatoryDaysHandler) Handle(ctx context.Context, q ListCompensatoryDaysQuery) (*CompensatoryDaysResult, error) {
	days, err := h.repo.FindByEmployee(ctx, q.EmployeeID)
	if err != nil {
		return nil, err
	}

	count, err := h.repo.CountAvailable(ctx, q.EmployeeID)
	if err != nil {
		return nil, err
	}

	return &CompensatoryDaysResult{
		Items:           days,
		AvailableCount:  count,
	}, nil
}

type CompensatoryDaysResult struct {
	Items          []*ausencia.CompensatoryDay
	AvailableCount int
}
