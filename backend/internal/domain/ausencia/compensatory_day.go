package ausencia

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CompensatoryDay struct {
	ID          string
	EmployeeID  string
	FechaOrigen time.Time
	Motivo      MotivoCompensatorio
	TurnoID     *string
	Descripcion string
	Utilizado   bool
	FechaUso    *time.Time
	CreatedAt   time.Time
}

type NewCompensatoryDayParams struct {
	EmployeeID  string
	FechaOrigen time.Time
	Motivo      MotivoCompensatorio
	TurnoID     *string
	Descripcion string
}

func NewCompensatoryDay(params NewCompensatoryDayParams) (*CompensatoryDay, error) {
	if params.EmployeeID == "" {
		return nil, fmt.Errorf("employee id is required")
	}
	if params.FechaOrigen.IsZero() {
		return nil, fmt.Errorf("fecha origen is required")
	}
	if !params.Motivo.IsValid() {
		return nil, fmt.Errorf("invalid motivo: %s", params.Motivo)
	}

	return &CompensatoryDay{
		ID:          uuid.New().String(),
		EmployeeID:  params.EmployeeID,
		FechaOrigen: params.FechaOrigen,
		Motivo:      params.Motivo,
		TurnoID:     params.TurnoID,
		Descripcion: params.Descripcion,
		Utilizado:   false,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (cd *CompensatoryDay) MarkAsUsed(fechaUso time.Time) error {
	if cd.Utilizado {
		return fmt.Errorf("compensatory day already used")
	}
	cd.Utilizado = true
	cd.FechaUso = &fechaUso
	return nil
}
