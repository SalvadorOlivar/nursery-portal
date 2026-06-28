package intercambio

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EstadoSwap string

const (
	SwapPendienteRespuesta EstadoSwap = "PENDIENTE_RESPUESTA"
	SwapPendienteAprobacion EstadoSwap = "PENDIENTE_APROBACION"
	SwapAprobado           EstadoSwap = "APROBADO"
	SwapRechazado          EstadoSwap = "RECHAZADO"
	SwapCancelado          EstadoSwap = "CANCELADO"
)

func (e EstadoSwap) IsValid() bool {
	switch e {
	case SwapPendienteRespuesta, SwapPendienteAprobacion, SwapAprobado, SwapRechazado, SwapCancelado:
		return true
	}
	return false
}

type AccionHistorial string

const (
	AccionSolicitado   AccionHistorial = "SOLICITADO"
	AccionAceptado     AccionHistorial = "ACEPTADO"
	AccionRechazado    AccionHistorial = "RECHAZADO"
	AccionAprobado     AccionHistorial = "APROBADO"
	AccionCancelado    AccionHistorial = "CANCELADO"
	AccionEjecutado    AccionHistorial = "EJECUTADO"
)

type ShiftSwapRequest struct {
	ID                 string
	PlanificacionID    string
	TurnoSolicitanteID string
	TurnoDestinoID     string
	SolicitanteID      string
	DestinoID          string
	Estado             EstadoSwap
	AprobadoPor        *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type NewSwapRequestParams struct {
	PlanificacionID    string
	TurnoSolicitanteID string
	TurnoDestinoID     string
	SolicitanteID      string
	DestinoID          string
}

func NewShiftSwapRequest(params NewSwapRequestParams) (*ShiftSwapRequest, error) {
	if params.PlanificacionID == "" {
		return nil, fmt.Errorf("el id de planificación es requerido")
	}
	if params.TurnoSolicitanteID == "" {
		return nil, fmt.Errorf("el id del turno solicitante es requerido")
	}
	if params.TurnoDestinoID == "" {
		return nil, fmt.Errorf("el id del turno destino es requerido")
	}
	if params.SolicitanteID == "" {
		return nil, fmt.Errorf("el id del solicitante es requerido")
	}
	if params.DestinoID == "" {
		return nil, fmt.Errorf("el id del destino es requerido")
	}
	if params.SolicitanteID == params.DestinoID {
		return nil, fmt.Errorf("no puedes intercambiar turnos contigo mismo")
	}
	if params.TurnoSolicitanteID == params.TurnoDestinoID {
		return nil, fmt.Errorf("no puedes intercambiar un turno consigo mismo")
	}

	now := time.Now().UTC()
	return &ShiftSwapRequest{
		ID:                 uuid.New().String(),
		PlanificacionID:    params.PlanificacionID,
		TurnoSolicitanteID: params.TurnoSolicitanteID,
		TurnoDestinoID:     params.TurnoDestinoID,
		SolicitanteID:      params.SolicitanteID,
		DestinoID:          params.DestinoID,
		Estado:             SwapPendienteRespuesta,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func (s *ShiftSwapRequest) AcceptByDestino() error {
	if s.Estado != SwapPendienteRespuesta {
		return fmt.Errorf("solo solicitudes en estado PENDIENTE_RESPUESTA pueden ser aceptadas")
	}
	s.Estado = SwapPendienteAprobacion
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *ShiftSwapRequest) Reject() error {
	if s.Estado != SwapPendienteRespuesta {
		return fmt.Errorf("solo solicitudes en estado PENDIENTE_RESPUESTA pueden ser rechazadas")
	}
	s.Estado = SwapRechazado
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *ShiftSwapRequest) Approve(approvedBy string) error {
	if s.Estado != SwapPendienteAprobacion {
		return fmt.Errorf("solo solicitudes en estado PENDIENTE_APROBACION pueden ser aprobadas")
	}
	s.Estado = SwapAprobado
	s.AprobadoPor = &approvedBy
	s.UpdatedAt = time.Now().UTC()
	return nil
}

func (s *ShiftSwapRequest) Cancel() error {
	if s.Estado != SwapPendienteRespuesta && s.Estado != SwapPendienteAprobacion {
		return fmt.Errorf("solo se pueden cancelar solicitudes pendientes")
	}
	s.Estado = SwapCancelado
	s.UpdatedAt = time.Now().UTC()
	return nil
}

type ShiftSwapHistoryEntry struct {
	ID            string
	SwapRequestID string
	Accion        AccionHistorial
	ActorID       string
	Detalle       *string
	CreatedAt     time.Time
}

func NewHistoryEntry(swapRequestID string, accion AccionHistorial, actorID string, detalle *string) *ShiftSwapHistoryEntry {
	return &ShiftSwapHistoryEntry{
		ID:            uuid.New().String(),
		SwapRequestID: swapRequestID,
		Accion:        accion,
		ActorID:       actorID,
		Detalle:       detalle,
		CreatedAt:     time.Now().UTC(),
	}
}
