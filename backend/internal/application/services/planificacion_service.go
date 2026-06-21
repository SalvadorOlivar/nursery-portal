package services

import (
	"context"
	"fmt"
	"time"

	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/domain/turno"
	"github.com/tuusuario/nursery-portal/internal/ports"

	cmdplanif "github.com/tuusuario/nursery-portal/internal/application/commands/planificacion"
	cmdturno "github.com/tuusuario/nursery-portal/internal/application/commands/turno"
	qryplanif "github.com/tuusuario/nursery-portal/internal/application/queries/planificacion"
)

type PlanificacionService struct {
	planifRepo            ports.PlanificacionRepository
	turnoRepo             ports.TurnoRepository
	sectorRepo            ports.SectorRepository
	dotRepo               ports.DotacionRepository
	employeeRepo          ports.EmployeeRepository
	leaveRepo             ports.LeaveRequestRepository
	compRepo              ports.CompensatoryDayRepository
	createHandler         *cmdplanif.CreatePlanificacionHandler
	updateHandler         *cmdplanif.UpdatePlanificacionHandler
	deleteHandler         *cmdplanif.DeletePlanificacionHandler
	publicarHandler       *cmdplanif.CambiarEstadoHandler
	cerrarHandler         *cmdplanif.CambiarEstadoHandler
	getByIDHandler        *qryplanif.GetPlanificacionHandler
	createTurnoHandler    *cmdturno.CreateTurnoHandler
	getDotacionHandler    *qryplanif.GetDotacionHandler
	updateSectoresHandler *cmdplanif.UpdateSectoresHandler
	updateDotacionHandler *cmdplanif.UpdateDotacionHandler
}

func NewPlanificacionService(
	planifRepo ports.PlanificacionRepository,
	turnoRepo ports.TurnoRepository,
	sectorRepo ports.SectorRepository,
	dotRepo ports.DotacionRepository,
	employeeRepo ports.EmployeeRepository,
	leaveRepo ports.LeaveRequestRepository,
	compRepo ports.CompensatoryDayRepository,
) *PlanificacionService {
	return &PlanificacionService{
		planifRepo:            planifRepo,
		turnoRepo:             turnoRepo,
		sectorRepo:            sectorRepo,
		dotRepo:               dotRepo,
		employeeRepo:          employeeRepo,
		leaveRepo:             leaveRepo,
		compRepo:              compRepo,
		createHandler:         cmdplanif.NewCreatePlanificacionHandler(planifRepo, sectorRepo, dotRepo),
		updateHandler:         cmdplanif.NewUpdatePlanificacionHandler(planifRepo),
		deleteHandler:         cmdplanif.NewDeletePlanificacionHandler(planifRepo),
		publicarHandler:       cmdplanif.NewCambiarEstadoHandler(planifRepo, (*planificacion.Planificacion).Publicar),
		cerrarHandler:         cmdplanif.NewCambiarEstadoHandler(planifRepo, (*planificacion.Planificacion).Cerrar),
		getByIDHandler:        qryplanif.NewGetPlanificacionHandler(planifRepo, turnoRepo),
		createTurnoHandler:    cmdturno.NewCreateTurnoHandler(turnoRepo),
		getDotacionHandler:    qryplanif.NewGetDotacionHandler(dotRepo),
		updateSectoresHandler: cmdplanif.NewUpdateSectoresHandler(planifRepo, sectorRepo),
		updateDotacionHandler: cmdplanif.NewUpdateDotacionHandler(planifRepo, dotRepo),
	}
}

func (s *PlanificacionService) Create(ctx context.Context, cmd cmdplanif.CreatePlanificacionCommand) (*planificacion.Planificacion, error) {
	return s.createHandler.Handle(ctx, cmd)
}

func (s *PlanificacionService) Update(ctx context.Context, cmd cmdplanif.UpdatePlanificacionCommand) error {
	return s.updateHandler.Handle(ctx, cmd)
}

func (s *PlanificacionService) Delete(ctx context.Context, id string) error {
	return s.deleteHandler.Handle(ctx, id)
}

func (s *PlanificacionService) Publicar(ctx context.Context, id string) error {
	return s.publicarHandler.Handle(ctx, id)
}

func (s *PlanificacionService) Cerrar(ctx context.Context, id string) error {
	return s.cerrarHandler.Handle(ctx, id)
}

func (s *PlanificacionService) List(ctx context.Context) ([]*planificacion.Planificacion, error) {
	return s.planifRepo.FindAll(ctx)
}

func (s *PlanificacionService) GetByID(ctx context.Context, id string) (*qryplanif.PlanificacionConTurnos, error) {
	return s.getByIDHandler.Handle(ctx, qryplanif.GetPlanificacionQuery{ID: id})
}

func (s *PlanificacionService) CreateTurno(ctx context.Context, cmd cmdturno.CreateTurnoCommand) (*turno.Turno, error) {
	planif, err := s.planifRepo.FindByID(ctx, cmd.PlanificacionID)
	if err != nil {
		return nil, fmt.Errorf("planificacion not found: %w", err)
	}

	fecha := isoWeekToDate(planif.Anio, planif.Semana, cmd.DiaSemana)

	leaves, err := s.leaveRepo.FindApprovedByEmployeeAndDate(ctx, cmd.EmpleadoID, fecha)
	if err != nil {
		return nil, fmt.Errorf("error checking leave requests: %w", err)
	}
	if len(leaves) > 0 {
		tipo := leaves[0].Tipo
		return nil, fmt.Errorf("el empleado tiene una licencia aprobada de tipo %s para el dia %s", tipo, fecha.Format("2006-01-02"))
	}

	t, err := s.createTurnoHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, err
	}

	existingTurnos, err := s.turnoRepo.FindByPlanificacionAndEmpleado(ctx, cmd.PlanificacionID, cmd.EmpleadoID)
	if err != nil {
		return nil, err
	}

	turnosEnMismoDia := 0
	for _, et := range existingTurnos {
		if et.DiaSemana == cmd.DiaSemana {
			turnosEnMismoDia++
		}
	}

	isRestDay := false
	emp, err := s.employeeRepo.FindByID(ctx, cmd.EmpleadoID)
	if err == nil && emp != nil {
		wp := emp.PatronTrabajo
		diasDesdeInicio := diasDesde(planif.Anio, planif.Semana, cmd.DiaSemana)
		if wp.IsRestDay(diasDesdeInicio) {
			isRestDay = true
		}
	}

	if turnosEnMismoDia >= 2 {
		desc := fmt.Sprintf("Turno doble el %s - %s", fecha.Format("2006-01-02"), cmd.Tipo)
		cd, cdErr := ausencia.NewCompensatoryDay(ausencia.NewCompensatoryDayParams{
			EmployeeID:  cmd.EmpleadoID,
			FechaOrigen: fecha,
			Motivo:      ausencia.DobleTurno,
			TurnoID:     &t.ID,
			Descripcion: desc,
		})
		if cdErr != nil {
			return nil, fmt.Errorf("turno creado pero error al generar dia compensatorio: %w", cdErr)
		}
		if err = s.compRepo.Create(ctx, cd); err != nil {
			return nil, fmt.Errorf("turno creado pero error al guardar dia compensatorio: %w", err)
		}
	} else if isRestDay {
		desc := fmt.Sprintf("Trabajo en dia de descanso el %s - %s", fecha.Format("2006-01-02"), cmd.Tipo)
		cd, cdErr := ausencia.NewCompensatoryDay(ausencia.NewCompensatoryDayParams{
			EmployeeID:  cmd.EmpleadoID,
			FechaOrigen: fecha,
			Motivo:      ausencia.DescansoLaborado,
			TurnoID:     &t.ID,
			Descripcion: desc,
		})
		if cdErr != nil {
			return nil, fmt.Errorf("turno creado pero error al generar dia compensatorio: %w", cdErr)
		}
		if err = s.compRepo.Create(ctx, cd); err != nil {
			return nil, fmt.Errorf("turno creado pero error al guardar dia compensatorio: %w", err)
		}
	}

	return t, nil
}

func (s *PlanificacionService) DeleteTurno(ctx context.Context, id string) error {
	return s.turnoRepo.Delete(ctx, id)
}

func (s *PlanificacionService) GetStaffingRequirements(ctx context.Context, planificacionID string) ([]planificacion.DotacionItem, error) {
	return s.getDotacionHandler.Handle(ctx, qryplanif.GetDotacionQuery{PlanificacionID: planificacionID})
}

func (s *PlanificacionService) GetPlanLeaves(ctx context.Context, planificacionID string) ([]*ausencia.LeaveRequest, error) {
	planif, err := s.planifRepo.FindByID(ctx, planificacionID)
	if err != nil {
		return nil, err
	}

	fechaInicio := isoWeekToDate(planif.Anio, planif.Semana, 1)
	fechaFin := isoWeekToDate(planif.Anio, planif.Semana, 7)

	return s.leaveRepo.FindByDateRange(ctx, fechaInicio, fechaFin)
}

func (s *PlanificacionService) GetSectores(ctx context.Context, planificacionID string) ([]*planificacion.SectorPlanificacion, error) {
	return s.sectorRepo.GetSectores(ctx, planificacionID)
}

func (s *PlanificacionService) UpdateSectores(ctx context.Context, cmd cmdplanif.UpdateSectoresCommand) error {
	return s.updateSectoresHandler.Handle(ctx, cmd)
}

func (s *PlanificacionService) UpdateDotacion(ctx context.Context, cmd cmdplanif.UpdateDotacionCommand) error {
	return s.updateDotacionHandler.Handle(ctx, cmd)
}

func isoWeekToDate(anio, semana, diaSemana int) time.Time {
	jan4 := time.Date(anio, 1, 4, 0, 0, 0, 0, time.UTC)
	daysSinceMonday := int(jan4.Weekday()) - 1
	if jan4.Weekday() == time.Sunday {
		daysSinceMonday = 6
	}
	mondayWeek1 := jan4.AddDate(0, 0, -daysSinceMonday)
	return mondayWeek1.AddDate(0, 0, (semana-1)*7+(diaSemana-1))
}

func diasDesde(anio, semana, diaSemana int) int {
	fecha := isoWeekToDate(anio, semana, diaSemana)
	inicioAnio := time.Date(anio, 1, 1, 0, 0, 0, 0, time.UTC)
	return int(fecha.Sub(inicioAnio).Hours()/24) + 1
}

