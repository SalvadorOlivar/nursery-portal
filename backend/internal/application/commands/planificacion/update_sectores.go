package planificacion

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/ports"
)

type UpdateSectoresCommand struct {
	PlanificacionID string
	Sectores        []string
}

type UpdateSectoresHandler struct {
	planifRepo ports.PlanificacionRepository
	sectorRepo ports.SectorRepository
}

func NewUpdateSectoresHandler(planifRepo ports.PlanificacionRepository, sectorRepo ports.SectorRepository) *UpdateSectoresHandler {
	return &UpdateSectoresHandler{
		planifRepo: planifRepo,
		sectorRepo: sectorRepo,
	}
}

func (h *UpdateSectoresHandler) Handle(ctx context.Context, cmd UpdateSectoresCommand) error {
	if _, err := h.planifRepo.FindByID(ctx, cmd.PlanificacionID); err != nil {
		return err
	}
	return h.sectorRepo.SaveSectores(ctx, cmd.PlanificacionID, cmd.Sectores)
}
