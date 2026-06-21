package services

import (
	"context"
	"fmt"

	domain "github.com/tuusuario/nursery-portal/internal/domain/employee"
	"github.com/tuusuario/nursery-portal/internal/ports"

	cmd "github.com/tuusuario/nursery-portal/internal/application/commands/employee"
)

type EmployeeService struct {
	repo              ports.EmployeeRepository
	authSvc           *AuthService
	createHandler     *cmd.CreateEmployeeHandler
	updateHandler     *cmd.UpdateEmployeeHandler
	deactivateHandler *cmd.DeactivateEmployeeHandler
}

func NewEmployeeService(repo ports.EmployeeRepository, authSvc *AuthService) *EmployeeService {
	return &EmployeeService{
		repo:              repo,
		authSvc:           authSvc,
		createHandler:     cmd.NewCreateEmployeeHandler(repo),
		updateHandler:     cmd.NewUpdateEmployeeHandler(repo),
		deactivateHandler: cmd.NewDeactivateEmployeeHandler(repo),
	}
}

func (s *EmployeeService) Create(ctx context.Context, params cmd.CreateEmployeeCommand) (*domain.Employee, string, error) {
	emp, err := s.createHandler.Handle(ctx, params)
	if err != nil {
		return nil, "", err
	}
	if s.authSvc == nil {
		return emp, "", nil
	}
	password, err := s.authSvc.CreateEmployeeAccount(ctx, emp)
	if err != nil {
		_ = s.repo.Delete(ctx, emp.ID)
		return nil, "", err
	}
	return emp, password, nil
}

func (s *EmployeeService) Update(ctx context.Context, params cmd.UpdateEmployeeCommand) (*domain.Employee, error) {
	emp, err := s.updateHandler.Handle(ctx, params)
	if err != nil {
		return nil, err
	}
	if s.authSvc != nil {
		if err := s.authSvc.UpdateEmployeeAccount(ctx, emp); err != nil {
			return nil, err
		}
	}
	return emp, nil
}

func (s *EmployeeService) Deactivate(ctx context.Context, id string) error {
	return s.deactivateHandler.Handle(ctx, cmd.DeactivateEmployeeCommand{ID: id})
}

func (s *EmployeeService) List(ctx context.Context) ([]*domain.Employee, error) {
	return s.repo.FindAll(ctx)
}

func (s *EmployeeService) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	emp, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}
	return emp, nil
}
