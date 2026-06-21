package planificacion

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type UpdatePlanificacionCommand struct {
	ID     string
	Nombre string
}

type UpdatePlanificacionHandler struct {
	repo ports.PlanificacionRepository
}

func NewUpdatePlanificacionHandler(repo ports.PlanificacionRepository) *UpdatePlanificacionHandler {
	return &UpdatePlanificacionHandler{repo: repo}
}

func (h *UpdatePlanificacionHandler) Handle(ctx context.Context, cmd UpdatePlanificacionCommand) error {
	p, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("planificacion not found: %w", err)
	}
	if p.Estado != planificacion.EstadoBorrador {
		return fmt.Errorf("only plans in BORRADOR can be updated")
	}
	p.Nombre = cmd.Nombre
	return h.repo.Update(ctx, p)
}
