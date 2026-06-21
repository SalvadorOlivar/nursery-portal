package turno

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TipoTurno string

const (
	Maniana    TipoTurno = "MANANA"
	Tarde      TipoTurno = "TARDE"
	Vespertino TipoTurno = "VESPERTINO"
	Noche      TipoTurno = "NOCHE"
)

func (t TipoTurno) IsValid() bool {
	switch t {
	case Maniana, Tarde, Vespertino, Noche:
		return true
	}
	return false
}

var AllTiposTurno = []TipoTurno{Maniana, Tarde, Vespertino, Noche}

type Turno struct {
	ID              string
	PlanificacionID string
	EmpleadoID      string
	DiaSemana       int
	Tipo            TipoTurno
	Sector          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type NewTurnoParams struct {
	PlanificacionID string
	EmpleadoID      string
	DiaSemana       int
	Tipo            TipoTurno
	Sector          string
}

func NewTurno(params NewTurnoParams) (*Turno, error) {
	if params.PlanificacionID == "" {
		return nil, fmt.Errorf("planificacion id is required")
	}
	if params.EmpleadoID == "" {
		return nil, fmt.Errorf("empleado id is required")
	}
	if params.DiaSemana < 1 || params.DiaSemana > 7 {
		return nil, fmt.Errorf("dia_semana must be between 1 and 7")
	}
	if !params.Tipo.IsValid() {
		return nil, fmt.Errorf("invalid turno type: %s", params.Tipo)
	}

	now := time.Now().UTC()
	return &Turno{
		ID:              uuid.New().String(),
		PlanificacionID: params.PlanificacionID,
		EmpleadoID:      params.EmpleadoID,
		DiaSemana:       params.DiaSemana,
		Tipo:            params.Tipo,
		Sector:          params.Sector,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}
