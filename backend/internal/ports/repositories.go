package ports

import (
	"context"
	"time"

	"github.com/tuusuario/nursery-portal/internal/domain/auth"
	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	"github.com/tuusuario/nursery-portal/internal/domain/employee"
	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/tuusuario/nursery-portal/internal/domain/turno"
)

type EmployeeRepository interface {
	Create(ctx context.Context, e *employee.Employee) error
	Update(ctx context.Context, e *employee.Employee) error
	FindByID(ctx context.Context, id string) (*employee.Employee, error)
	FindAll(ctx context.Context) ([]*employee.Employee, error)
	Delete(ctx context.Context, id string) error
}

type PlanificacionRepository interface {
	Create(ctx context.Context, p *planificacion.Planificacion) error
	Update(ctx context.Context, p *planificacion.Planificacion) error
	FindByID(ctx context.Context, id string) (*planificacion.Planificacion, error)
	FindAll(ctx context.Context) ([]*planificacion.Planificacion, error)
	Delete(ctx context.Context, id string) error
	DeleteCascade(ctx context.Context, id string) error
}

type TurnoRepository interface {
	Create(ctx context.Context, t *turno.Turno) error
	CreateBatch(ctx context.Context, turnos []*turno.Turno) error
	FindByTurnoID(ctx context.Context, id string) (*turno.Turno, error)
	FindByPlanificacion(ctx context.Context, planificacionID string) ([]*turno.Turno, error)
	FindByPlanificacionAndEmpleado(ctx context.Context, planificacionID, empleadoID string) ([]*turno.Turno, error)
	Update(ctx context.Context, t *turno.Turno) error
	Delete(ctx context.Context, id string) error
	DeleteByPlanificacion(ctx context.Context, planificacionID string) error
}

type SectorRepository interface {
	GetSectores(ctx context.Context, planificacionID string) ([]*planificacion.SectorPlanificacion, error)
	SaveSectores(ctx context.Context, planificacionID string, nombres []string) error
}

type DotacionRepository interface {
	GetDotacion(ctx context.Context, planificacionID string) ([]*planificacion.DotacionPlanificacion, error)
	SaveDotacion(ctx context.Context, items []*planificacion.DotacionPlanificacion) error
	DeleteByPlanificacion(ctx context.Context, planificacionID string) error
}

type AuthRepository interface {
	FindUserByUsername(ctx context.Context, username string) (*auth.User, error)
	FindUserBySessionHash(ctx context.Context, tokenHash string, now time.Time) (*auth.User, error)
	SetPasswordHash(ctx context.Context, userID, passwordHash string) error
	CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	DeleteSession(ctx context.Context, tokenHash string) error
	EnsureAdmin(ctx context.Context, username, passwordHash string) error
	CreateEmployeeUser(ctx context.Context, username string, role auth.Role, employeeID string, passwordHash string) error
	UpdateEmployeeUser(ctx context.Context, employeeID, username string, role auth.Role) error
}

type LeaveRequestRepository interface {
	Create(ctx context.Context, lr *ausencia.LeaveRequest) error
	FindByID(ctx context.Context, id string) (*ausencia.LeaveRequest, error)
	FindByEmployee(ctx context.Context, employeeID string) ([]*ausencia.LeaveRequest, error)
	FindAll(ctx context.Context) ([]*ausencia.LeaveRequest, error)
	FindByDateRange(ctx context.Context, fechaInicio, fechaFin time.Time) ([]*ausencia.LeaveRequest, error)
	FindApprovedByEmployeeAndDate(ctx context.Context, employeeID string, fecha time.Time) ([]*ausencia.LeaveRequest, error)
	Update(ctx context.Context, lr *ausencia.LeaveRequest) error
}

type CompensatoryDayRepository interface {
	Create(ctx context.Context, cd *ausencia.CompensatoryDay) error
	FindByEmployee(ctx context.Context, employeeID string) ([]*ausencia.CompensatoryDay, error)
	FindByID(ctx context.Context, id string) (*ausencia.CompensatoryDay, error)
	FindAvailableByEmployee(ctx context.Context, employeeID string) ([]*ausencia.CompensatoryDay, error)
	Update(ctx context.Context, cd *ausencia.CompensatoryDay) error
	CountAvailable(ctx context.Context, employeeID string) (int, error)
}

type ShiftSwapRequestRepository interface {
	Create(ctx context.Context, req *intercambio.ShiftSwapRequest) error
	FindByID(ctx context.Context, id string) (*intercambio.ShiftSwapRequest, error)
	FindAll(ctx context.Context) ([]*intercambio.ShiftSwapRequest, error)
	Update(ctx context.Context, req *intercambio.ShiftSwapRequest) error
	AddHistoryEntry(ctx context.Context, entry *intercambio.ShiftSwapHistoryEntry) error
	GetHistory(ctx context.Context, swapRequestID string) ([]*intercambio.ShiftSwapHistoryEntry, error)
}
