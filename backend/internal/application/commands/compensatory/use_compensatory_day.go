package compensatory

import (
	"context"
	"fmt"
	"time"

	"github.com/tuusuario/nursery-portal/internal/ports"
)

type UseCompensatoryDayCommand struct {
	ID       string
	FechaUso time.Time
}

type UseCompensatoryDayHandler struct {
	repo ports.CompensatoryDayRepository
}

func NewUseCompensatoryDayHandler(repo ports.CompensatoryDayRepository) *UseCompensatoryDayHandler {
	return &UseCompensatoryDayHandler{repo: repo}
}

func (h *UseCompensatoryDayHandler) Handle(ctx context.Context, cmd UseCompensatoryDayCommand) error {
	cd, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("compensatory day not found: %w", err)
	}

	if err := cd.MarkAsUsed(cmd.FechaUso); err != nil {
		return err
	}

	return h.repo.Update(ctx, cd)
}
