package intercambio

import (
	"context"
	"fmt"

	"github.com/tuusuario/nurse-portal/internal/domain/intercambio"
	"github.com/tuusuario/nurse-portal/internal/domain/planificacion"
	"github.com/tuusuario/nurse-portal/internal/ports"
)

type CreateSwapRequestCommand struct {
	PlanificacionID    string
	TurnoSolicitanteID string
	TurnoDestinoID     string
	SolicitanteID      string
	DestinoID          string
	ActorID            string
}

type CreateSwapRequestHandler struct {
	repo       ports.ShiftSwapRequestRepository
	turnoRepo  ports.TurnoRepository
	planifRepo ports.PlanificacionRepository
	leaveRepo  ports.LeaveRequestRepository
}

func NewCreateSwapRequestHandler(
	repo ports.ShiftSwapRequestRepository,
	turnoRepo ports.TurnoRepository,
	planifRepo ports.PlanificacionRepository,
	leaveRepo ports.LeaveRequestRepository,
) *CreateSwapRequestHandler {
	return &CreateSwapRequestHandler{
		repo:       repo,
		turnoRepo:  turnoRepo,
		planifRepo: planifRepo,
		leaveRepo:  leaveRepo,
	}
}

func (h *CreateSwapRequestHandler) Handle(ctx context.Context, cmd CreateSwapRequestCommand) (*intercambio.ShiftSwapRequest, error) {
	turnoSol, err := h.turnoRepo.FindByTurnoID(ctx, cmd.TurnoSolicitanteID)
	if err != nil {
		return nil, fmt.Errorf("turno solicitante no encontrado: %w", err)
	}
	turnoDest, err := h.turnoRepo.FindByTurnoID(ctx, cmd.TurnoDestinoID)
	if err != nil {
		return nil, fmt.Errorf("turno destino no encontrado: %w", err)
	}

	if turnoSol.PlanificacionID != cmd.PlanificacionID || turnoDest.PlanificacionID != cmd.PlanificacionID {
		return nil, fmt.Errorf("ambos turnos deben pertenecer a la misma planificación")
	}
	if turnoSol.EmpleadoID != cmd.SolicitanteID {
		return nil, fmt.Errorf("el turno solicitante no pertenece al empleado solicitante")
	}
	if turnoDest.EmpleadoID != cmd.DestinoID {
		return nil, fmt.Errorf("el turno destino no pertenece al empleado destino")
	}

	p, err := h.planifRepo.FindByID(ctx, cmd.PlanificacionID)
	if err != nil {
		return nil, fmt.Errorf("planificación no encontrada: %w", err)
	}

	fechaTurnoSol := planificacion.TurnoDate(p.Semana, p.Anio, turnoSol.DiaSemana)
	fechaTurnoDest := planificacion.TurnoDate(p.Semana, p.Anio, turnoDest.DiaSemana)

	leavesSol, err := h.leaveRepo.FindApprovedByEmployeeAndDate(ctx, cmd.SolicitanteID, fechaTurnoSol)
	if err != nil {
		return nil, fmt.Errorf("error al verificar licencia del solicitante: %w", err)
	}
	if len(leavesSol) > 0 {
		return nil, fmt.Errorf(
			"el empleado solicitante tiene licencia aprobada el %s",
			fechaTurnoSol.Format("2006-01-02"),
		)
	}

	leavesDestByDate, err := h.leaveRepo.FindApprovedByEmployeeAndDate(ctx, cmd.DestinoID, fechaTurnoDest)
	if err != nil {
		return nil, fmt.Errorf("error al verificar licencia del destino: %w", err)
	}
	if len(leavesDestByDate) > 0 {
		return nil, fmt.Errorf(
			"el empleado destino tiene licencia aprobada el %s",
			fechaTurnoDest.Format("2006-01-02"),
		)
	}

	req, err := intercambio.NewShiftSwapRequest(intercambio.NewSwapRequestParams{
		PlanificacionID:    cmd.PlanificacionID,
		TurnoSolicitanteID: cmd.TurnoSolicitanteID,
		TurnoDestinoID:     cmd.TurnoDestinoID,
		SolicitanteID:      cmd.SolicitanteID,
		DestinoID:          cmd.DestinoID,
	})
	if err != nil {
		return nil, err
	}

	if err := h.repo.Create(ctx, req); err != nil {
		return nil, err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionSolicitado, cmd.ActorID, nil)
	if err := h.repo.AddHistoryEntry(ctx, entry); err != nil {
		return nil, err
	}

	return req, nil
}
