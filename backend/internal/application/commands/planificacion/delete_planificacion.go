package planificacion

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/ports"
)

type DeletePlanificacionHandler struct {
	planificacionRepo ports.PlanificacionRepository
}

func NewDeletePlanificacionHandler(planificacionRepo ports.PlanificacionRepository) *DeletePlanificacionHandler {
	return &DeletePlanificacionHandler{
		planificacionRepo: planificacionRepo,
	}
}

func (h *DeletePlanificacionHandler) Handle(ctx context.Context, id string) error {
	return h.planificacionRepo.DeleteCascade(ctx, id)
}
