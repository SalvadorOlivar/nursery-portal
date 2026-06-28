package intercambio

import (
	"context"
	"fmt"

	"github.com/tuusuario/nurse-portal/internal/domain/intercambio"
	"github.com/tuusuario/nurse-portal/internal/domain/planificacion"
	"github.com/tuusuario/nurse-portal/internal/ports"
)

type ApproveSwapRequestCommand struct {
	ID      string
	ActorID string
}

type ApproveSwapRequestHandler struct {
	repo       ports.ShiftSwapRequestRepository
	turnoRepo  ports.TurnoRepository
	planifRepo ports.PlanificacionRepository
	leaveRepo  ports.LeaveRequestRepository
}

func NewApproveSwapRequestHandler(
	repo ports.ShiftSwapRequestRepository,
	turnoRepo ports.TurnoRepository,
	planifRepo ports.PlanificacionRepository,
	leaveRepo ports.LeaveRequestRepository,
) *ApproveSwapRequestHandler {
	return &ApproveSwapRequestHandler{
		repo:       repo,
		turnoRepo:  turnoRepo,
		planifRepo: planifRepo,
		leaveRepo:  leaveRepo,
	}
}

func (h *ApproveSwapRequestHandler) Handle(ctx context.Context, cmd ApproveSwapRequestCommand) error {
	req, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("solicitud de intercambio no encontrada: %w", err)
	}

	if err := req.Approve(cmd.ActorID); err != nil {
		return err
	}

	turnoSol, err := h.turnoRepo.FindByTurnoID(ctx, req.TurnoSolicitanteID)
	if err != nil {
		return fmt.Errorf("turno solicitante no encontrado: %w", err)
	}

	turnoDest, err := h.turnoRepo.FindByTurnoID(ctx, req.TurnoDestinoID)
	if err != nil {
		return fmt.Errorf("turno destino no encontrado: %w", err)
	}

	p, err := h.planifRepo.FindByID(ctx, req.PlanificacionID)
	if err != nil {
		return fmt.Errorf("planificación no encontrada: %w", err)
	}

	fechaTurnoSol := planificacion.TurnoDate(p.Semana, p.Anio, turnoSol.DiaSemana)
	fechaTurnoDest := planificacion.TurnoDate(p.Semana, p.Anio, turnoDest.DiaSemana)

	leavesSol, err := h.leaveRepo.FindApprovedByEmployeeAndDate(ctx, turnoSol.EmpleadoID, fechaTurnoSol)
	if err != nil {
		return fmt.Errorf("error al verificar licencia del solicitante: %w", err)
	}
	if len(leavesSol) > 0 {
		return fmt.Errorf(
			"el empleado solicitante tiene licencia aprobada el %s",
			fechaTurnoSol.Format("2006-01-02"),
		)
	}

	leavesDest, err := h.leaveRepo.FindApprovedByEmployeeAndDate(ctx, turnoDest.EmpleadoID, fechaTurnoDest)
	if err != nil {
		return fmt.Errorf("error al verificar licencia del destino: %w", err)
	}
	if len(leavesDest) > 0 {
		return fmt.Errorf(
			"el empleado destino tiene licencia aprobada el %s",
			fechaTurnoDest.Format("2006-01-02"),
		)
	}

	solicitanteEmpleadoID := turnoSol.EmpleadoID
	destinoEmpleadoID := turnoDest.EmpleadoID

	turnoSol.EmpleadoID = destinoEmpleadoID
	turnoDest.EmpleadoID = solicitanteEmpleadoID

	if err := h.turnoRepo.Update(ctx, turnoSol); err != nil {
		return fmt.Errorf("error al actualizar turno solicitante: %w", err)
	}
	if err := h.turnoRepo.Update(ctx, turnoDest); err != nil {
		return fmt.Errorf("error al actualizar turno destino: %w", err)
	}

	if err := h.repo.Update(ctx, req); err != nil {
		return err
	}

	entry := intercambio.NewHistoryEntry(req.ID, intercambio.AccionEjecutado, cmd.ActorID, nil)
	if err := h.repo.AddHistoryEntry(ctx, entry); err != nil {
		return err
	}

	return nil
}
