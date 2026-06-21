package ausencia

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LeaveRequest struct {
	ID          string
	EmployeeID  string
	FechaInicio time.Time
	FechaFin    time.Time
	Tipo        TipoAusencia
	Estado      EstadoAusencia
	Motivo      string
	AprobadoPor *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewLeaveRequestParams struct {
	EmployeeID  string
	FechaInicio time.Time
	FechaFin    time.Time
	Tipo        TipoAusencia
	Motivo      string
}

func NewLeaveRequest(params NewLeaveRequestParams) (*LeaveRequest, error) {
	if params.EmployeeID == "" {
		return nil, fmt.Errorf("employee id is required")
	}
	if params.FechaInicio.IsZero() {
		return nil, fmt.Errorf("fecha inicio is required")
	}
	if params.FechaFin.IsZero() {
		return nil, fmt.Errorf("fecha fin is required")
	}
	if params.FechaFin.Before(params.FechaInicio) {
		return nil, fmt.Errorf("fecha fin must be after or equal to fecha inicio")
	}
	if !params.Tipo.IsValid() {
		return nil, fmt.Errorf("invalid leave type: %s", params.Tipo)
	}

	now := time.Now().UTC()
	return &LeaveRequest{
		ID:          uuid.New().String(),
		EmployeeID:  params.EmployeeID,
		FechaInicio: params.FechaInicio,
		FechaFin:    params.FechaFin,
		Tipo:        params.Tipo,
		Estado:      Pendiente,
		Motivo:      params.Motivo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (lr *LeaveRequest) Approve(approvedBy string) error {
	if lr.Estado != Pendiente {
		return fmt.Errorf("only pending requests can be approved")
	}
	lr.Estado = Aprobado
	lr.AprobadoPor = &approvedBy
	lr.UpdatedAt = time.Now().UTC()
	return nil
}

func (lr *LeaveRequest) Reject(approvedBy string) error {
	if lr.Estado != Pendiente {
		return fmt.Errorf("only pending requests can be rejected")
	}
	lr.Estado = Rechazado
	lr.AprobadoPor = &approvedBy
	lr.UpdatedAt = time.Now().UTC()
	return nil
}

func (lr *LeaveRequest) DiasSolicitados() int {
	return int(lr.FechaFin.Sub(lr.FechaInicio).Hours()/24) + 1
}
