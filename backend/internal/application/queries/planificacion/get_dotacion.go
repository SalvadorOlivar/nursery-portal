package planificacion

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/domain/employee"
	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type GetDotacionQuery struct {
	PlanificacionID string
}

type GetDotacionHandler struct {
	dotRepo ports.DotacionRepository
}

func NewGetDotacionHandler(dotRepo ports.DotacionRepository) *GetDotacionHandler {
	return &GetDotacionHandler{dotRepo: dotRepo}
}

func (h *GetDotacionHandler) Handle(ctx context.Context, query GetDotacionQuery) ([]planificacion.DotacionItem, error) {
	items, err := h.dotRepo.GetDotacion(ctx, query.PlanificacionID)
	if err != nil {
		return nil, err
	}

	result := make([]planificacion.DotacionItem, len(items))
	for i, d := range items {
		result[i] = planificacion.DotacionItem{
			TipoEmpleado:   employee.Type(d.TipoEmpleado),
			Turno:          d.Turno,
			CantidadMinima: d.CantidadMinima,
			Sector:         d.Sector,
		}
	}
	return result, nil
}
