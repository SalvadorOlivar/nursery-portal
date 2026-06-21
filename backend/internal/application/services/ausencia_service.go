package services

import (
	"context"
	"time"

	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
	cmdcomp "github.com/tuusuario/nursery-portal/internal/application/commands/compensatory"
	cmdleave "github.com/tuusuario/nursery-portal/internal/application/commands/leave"
	qrycomp "github.com/tuusuario/nursery-portal/internal/application/queries/compensatory"
	qryleave "github.com/tuusuario/nursery-portal/internal/application/queries/leave"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type AusenciaService struct {
	leaveRepo      ports.LeaveRequestRepository
	compRepo       ports.CompensatoryDayRepository

	createLeaveHandler     *cmdleave.CreateLeaveRequestHandler
	approveLeaveHandler    *cmdleave.ApproveLeaveRequestHandler
	rejectLeaveHandler     *cmdleave.RejectLeaveRequestHandler
	listLeaveHandler       *qryleave.ListLeaveRequestsHandler
	createCompHandler      *cmdcomp.CreateCompensatoryDayHandler
	useCompHandler         *cmdcomp.UseCompensatoryDayHandler
	listCompHandler        *qrycomp.ListCompensatoryDaysHandler
}

func NewAusenciaService(
	leaveRepo ports.LeaveRequestRepository,
	compRepo ports.CompensatoryDayRepository,
) *AusenciaService {
	return &AusenciaService{
		leaveRepo:            leaveRepo,
		compRepo:             compRepo,
		createLeaveHandler:   cmdleave.NewCreateLeaveRequestHandler(leaveRepo),
		approveLeaveHandler:  cmdleave.NewApproveLeaveRequestHandler(leaveRepo),
		rejectLeaveHandler:   cmdleave.NewRejectLeaveRequestHandler(leaveRepo),
		listLeaveHandler:     qryleave.NewListLeaveRequestsHandler(leaveRepo),
		createCompHandler:    cmdcomp.NewCreateCompensatoryDayHandler(compRepo),
		useCompHandler:       cmdcomp.NewUseCompensatoryDayHandler(compRepo),
		listCompHandler:      qrycomp.NewListCompensatoryDaysHandler(compRepo),
	}
}

func (s *AusenciaService) CreateLeaveRequest(ctx context.Context, cmd cmdleave.CreateLeaveRequestCommand) (*ausencia.LeaveRequest, error) {
	return s.createLeaveHandler.Handle(ctx, cmd)
}

func (s *AusenciaService) ApproveLeaveRequest(ctx context.Context, id, approvedBy string) error {
	return s.approveLeaveHandler.Handle(ctx, cmdleave.ApproveLeaveRequestCommand{ID: id, ApprovedBy: approvedBy})
}

func (s *AusenciaService) RejectLeaveRequest(ctx context.Context, id, approvedBy string) error {
	return s.rejectLeaveHandler.Handle(ctx, cmdleave.RejectLeaveRequestCommand{ID: id, ApprovedBy: approvedBy})
}

func (s *AusenciaService) ListLeaveRequests(ctx context.Context, employeeID string) ([]*ausencia.LeaveRequest, error) {
	return s.listLeaveHandler.Handle(ctx, qryleave.ListLeaveRequestsQuery{EmployeeID: employeeID})
}

func (s *AusenciaService) GetLeaveRequest(ctx context.Context, id string) (*ausencia.LeaveRequest, error) {
	return s.leaveRepo.FindByID(ctx, id)
}

func (s *AusenciaService) FindApprovedLeaves(ctx context.Context, employeeID string, fecha time.Time) ([]*ausencia.LeaveRequest, error) {
	return s.leaveRepo.FindApprovedByEmployeeAndDate(ctx, employeeID, fecha)
}

func (s *AusenciaService) CreateCompensatoryDay(ctx context.Context, cmd cmdcomp.CreateCompensatoryDayCommand) (*ausencia.CompensatoryDay, error) {
	return s.createCompHandler.Handle(ctx, cmd)
}

func (s *AusenciaService) UseCompensatoryDay(ctx context.Context, id string, fechaUso time.Time) error {
	return s.useCompHandler.Handle(ctx, cmdcomp.UseCompensatoryDayCommand{ID: id, FechaUso: fechaUso})
}

func (s *AusenciaService) ListCompensatoryDays(ctx context.Context, employeeID string) (*qrycomp.CompensatoryDaysResult, error) {
	return s.listCompHandler.Handle(ctx, qrycomp.ListCompensatoryDaysQuery{EmployeeID: employeeID})
}

func (s *AusenciaService) CountAvailableCompensatoryDays(ctx context.Context, employeeID string) (int, error) {
	return s.compRepo.CountAvailable(ctx, employeeID)
}
