package planificacion

import (
	"context"
	"fmt"

	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type CambiarEstadoHandler struct {
	repo       ports.PlanificacionRepository
	transicion func(*planificacion.Planificacion) error
}

func NewCambiarEstadoHandler(repo ports.PlanificacionRepository, transicion func(*planificacion.Planificacion) error) *CambiarEstadoHandler {
	return &CambiarEstadoHandler{repo: repo, transicion: transicion}
}

func (h *CambiarEstadoHandler) Handle(ctx context.Context, id string) error {
	p, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("planificacion not found: %w", err)
	}
	if err := h.transicion(p); err != nil {
		return err
	}
	return h.repo.Update(ctx, p)
}
