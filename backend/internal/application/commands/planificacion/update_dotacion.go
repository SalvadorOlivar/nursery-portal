package planificacion

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type UpdateDotacionCommand struct {
	PlanificacionID string
	Items           []DotacionItemInput
}

type DotacionItemInput struct {
	Sector         string
	TipoEmpleado   string
	Turno          string
	CantidadMinima int
}

type UpdateDotacionHandler struct {
	planifRepo ports.PlanificacionRepository
	dotRepo    ports.DotacionRepository
}

func NewUpdateDotacionHandler(planifRepo ports.PlanificacionRepository, dotRepo ports.DotacionRepository) *UpdateDotacionHandler {
	return &UpdateDotacionHandler{
		planifRepo: planifRepo,
		dotRepo:    dotRepo,
	}
}

func (h *UpdateDotacionHandler) Handle(ctx context.Context, cmd UpdateDotacionCommand) error {
	if _, err := h.planifRepo.FindByID(ctx, cmd.PlanificacionID); err != nil {
		return err
	}

	items := make([]*planificacion.DotacionPlanificacion, len(cmd.Items))
	for i, in := range cmd.Items {
		items[i] = &planificacion.DotacionPlanificacion{
			PlanificacionID: cmd.PlanificacionID,
			Sector:          in.Sector,
			TipoEmpleado:    in.TipoEmpleado,
			Turno:           in.Turno,
			CantidadMinima:  in.CantidadMinima,
		}
	}

	return h.dotRepo.SaveDotacion(ctx, items)
}
