package leave

import (
	"context"
	"time"

	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type CreateLeaveRequestCommand struct {
	EmployeeID  string
	FechaInicio time.Time
	FechaFin    time.Time
	Tipo        string
	Motivo      string
}

type CreateLeaveRequestHandler struct {
	repo ports.LeaveRequestRepository
}

func NewCreateLeaveRequestHandler(repo ports.LeaveRequestRepository) *CreateLeaveRequestHandler {
	return &CreateLeaveRequestHandler{repo: repo}
}

func (h *CreateLeaveRequestHandler) Handle(ctx context.Context, cmd CreateLeaveRequestCommand) (*ausencia.LeaveRequest, error) {
	tipo, err := ausencia.ParseTipoAusencia(cmd.Tipo)
	if err != nil {
		return nil, err
	}

	lr, err := ausencia.NewLeaveRequest(ausencia.NewLeaveRequestParams{
		EmployeeID:  cmd.EmployeeID,
		FechaInicio: cmd.FechaInicio,
		FechaFin:    cmd.FechaFin,
		Tipo:        tipo,
		Motivo:      cmd.Motivo,
	})
	if err != nil {
		return nil, err
	}

	if err := h.repo.Create(ctx, lr); err != nil {
		return nil, err
	}

	return lr, nil
}
