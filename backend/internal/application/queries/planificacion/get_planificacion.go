package planificacion

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/domain/turno"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type GetPlanificacionQuery struct {
	ID string
}

type PlanificacionConTurnos struct {
	Planificacion *planificacion.Planificacion
	Turnos        []*turno.Turno
}

type GetPlanificacionHandler struct {
	planifRepo ports.PlanificacionRepository
	turnoRepo  ports.TurnoRepository
}

func NewGetPlanificacionHandler(planifRepo ports.PlanificacionRepository, turnoRepo ports.TurnoRepository) *GetPlanificacionHandler {
	return &GetPlanificacionHandler{
		planifRepo: planifRepo,
		turnoRepo:  turnoRepo,
	}
}

func (h *GetPlanificacionHandler) Handle(ctx context.Context, qry GetPlanificacionQuery) (*PlanificacionConTurnos, error) {
	p, err := h.planifRepo.FindByID(ctx, qry.ID)
	if err != nil {
		return nil, fmt.Errorf("planificacion not found: %w", err)
	}

	turnos, err := h.turnoRepo.FindByPlanificacion(ctx, qry.ID)
	if err != nil {
		return nil, err
	}

	return &PlanificacionConTurnos{
		Planificacion: p,
		Turnos:        turnos,
	}, nil
}
