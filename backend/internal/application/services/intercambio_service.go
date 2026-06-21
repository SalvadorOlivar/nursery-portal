package services

import (
	"context"

	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
	cmdinter "github.com/tuusuario/nursery-portal/internal/application/commands/intercambio"
	qryinter "github.com/tuusuario/nursery-portal/internal/application/queries/intercambio"
	"github.com/tuusuario/nursery-portal/internal/ports"
)

type IntercambioService struct {
	swapRepo  ports.ShiftSwapRequestRepository
	turnoRepo ports.TurnoRepository

	createHandler  *cmdinter.CreateSwapRequestHandler
	acceptHandler  *cmdinter.AcceptSwapRequestHandler
	rejectHandler  *cmdinter.RejectSwapRequestHandler
	approveHandler *cmdinter.ApproveSwapRequestHandler
	cancelHandler  *cmdinter.CancelSwapRequestHandler
	listHandler    *qryinter.ListSwapRequestsHandler
	historyHandler *qryinter.GetSwapRequestHistoryHandler
}

func NewIntercambioService(
	swapRepo ports.ShiftSwapRequestRepository,
	turnoRepo ports.TurnoRepository,
) *IntercambioService {
	return &IntercambioService{
		swapRepo:       swapRepo,
		turnoRepo:      turnoRepo,
		createHandler:  cmdinter.NewCreateSwapRequestHandler(swapRepo, turnoRepo),
		acceptHandler:  cmdinter.NewAcceptSwapRequestHandler(swapRepo),
		rejectHandler:  cmdinter.NewRejectSwapRequestHandler(swapRepo),
		approveHandler: cmdinter.NewApproveSwapRequestHandler(swapRepo, turnoRepo),
		cancelHandler:  cmdinter.NewCancelSwapRequestHandler(swapRepo),
		listHandler:    qryinter.NewListSwapRequestsHandler(swapRepo),
		historyHandler: qryinter.NewGetSwapRequestHistoryHandler(swapRepo),
	}
}

func (s *IntercambioService) CreateSwapRequest(ctx context.Context, cmd cmdinter.CreateSwapRequestCommand) (*intercambio.ShiftSwapRequest, error) {
	return s.createHandler.Handle(ctx, cmd)
}

func (s *IntercambioService) AcceptSwapRequest(ctx context.Context, id, actorID string) error {
	return s.acceptHandler.Handle(ctx, cmdinter.AcceptSwapRequestCommand{ID: id, ActorID: actorID})
}

func (s *IntercambioService) RejectSwapRequest(ctx context.Context, id, actorID string) error {
	return s.rejectHandler.Handle(ctx, cmdinter.RejectSwapRequestCommand{ID: id, ActorID: actorID})
}

func (s *IntercambioService) ApproveSwapRequest(ctx context.Context, id, actorID string) error {
	return s.approveHandler.Handle(ctx, cmdinter.ApproveSwapRequestCommand{ID: id, ActorID: actorID})
}

func (s *IntercambioService) CancelSwapRequest(ctx context.Context, id, actorID string) error {
	return s.cancelHandler.Handle(ctx, cmdinter.CancelSwapRequestCommand{ID: id, ActorID: actorID})
}

func (s *IntercambioService) ListSwapRequests(ctx context.Context, employeeID, role string) ([]*intercambio.ShiftSwapRequest, error) {
	return s.listHandler.Handle(ctx, qryinter.ListSwapRequestsQuery{EmployeeID: employeeID, Role: role})
}

func (s *IntercambioService) GetSwapRequest(ctx context.Context, id string) (*intercambio.ShiftSwapRequest, error) {
	return s.swapRepo.FindByID(ctx, id)
}

func (s *IntercambioService) GetSwapHistory(ctx context.Context, swapRequestID string) ([]*intercambio.ShiftSwapHistoryEntry, error) {
	return s.historyHandler.Handle(ctx, swapRequestID)
}
