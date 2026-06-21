package compensatory

import (
	"context"
	"time"

	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type CreateCompensatoryDayCommand struct {
	EmployeeID  string
	FechaOrigen time.Time
	Motivo      string
	TurnoID     *string
	Descripcion string
}

type CreateCompensatoryDayHandler struct {
	repo ports.CompensatoryDayRepository
}

func NewCreateCompensatoryDayHandler(repo ports.CompensatoryDayRepository) *CreateCompensatoryDayHandler {
	return &CreateCompensatoryDayHandler{repo: repo}
}

func (h *CreateCompensatoryDayHandler) Handle(ctx context.Context, cmd CreateCompensatoryDayCommand) (*ausencia.CompensatoryDay, error) {
	cd, err := ausencia.NewCompensatoryDay(ausencia.NewCompensatoryDayParams{
		EmployeeID:  cmd.EmployeeID,
		FechaOrigen: cmd.FechaOrigen,
		Motivo:      ausencia.MotivoCompensatorio(cmd.Motivo),
		TurnoID:     cmd.TurnoID,
		Descripcion: cmd.Descripcion,
	})
	if err != nil {
		return nil, err
	}

	if err := h.repo.Create(ctx, cd); err != nil {
		return nil, err
	}

	return cd, nil
}
