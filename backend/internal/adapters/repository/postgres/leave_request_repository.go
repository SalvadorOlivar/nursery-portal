package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuusuario/nursery-portal/internal/domain/ausencia"
)

type LeaveRequestRepository struct {
	pool *pgxpool.Pool
}

func NewLeaveRequestRepository(pool *pgxpool.Pool) *LeaveRequestRepository {
	return &LeaveRequestRepository{pool: pool}
}

func (r *LeaveRequestRepository) Create(ctx context.Context, lr *ausencia.LeaveRequest) error {
	query := `
		INSERT INTO leave_requests (id, employee_id, fecha_inicio, fecha_fin, tipo, estado, motivo, aprobado_por, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		lr.ID, lr.EmployeeID, lr.FechaInicio, lr.FechaFin, string(lr.Tipo), string(lr.Estado),
		lr.Motivo, lr.AprobadoPor, lr.CreatedAt, lr.UpdatedAt,
	)
	return err
}

func (r *LeaveRequestRepository) FindByID(ctx context.Context, id string) (*ausencia.LeaveRequest, error) {
	query := `
		SELECT id, employee_id, fecha_inicio, fecha_fin, tipo, estado, motivo, aprobado_por, created_at, updated_at
		FROM leave_requests
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	return scanLeaveRequest(row)
}

func (r *LeaveRequestRepository) FindByEmployee(ctx context.Context, employeeID string) ([]*ausencia.LeaveRequest, error) {
	query := `
		SELECT id, employee_id, fecha_inicio, fecha_fin, tipo, estado, motivo, aprobado_por, created_at, updated_at
		FROM leave_requests
		WHERE employee_id = $1
		ORDER BY created_at DESC
	`
	return r.queryLeaveRequests(ctx, query, employeeID)
}

func (r *LeaveRequestRepository) FindAll(ctx context.Context) ([]*ausencia.LeaveRequest, error) {
	query := `
		SELECT id, employee_id, fecha_inicio, fecha_fin, tipo, estado, motivo, aprobado_por, created_at, updated_at
		FROM leave_requests
		ORDER BY created_at DESC
	`
	return r.queryLeaveRequests(ctx, query)
}

func (r *LeaveRequestRepository) FindByDateRange(ctx context.Context, fechaInicio, fechaFin time.Time) ([]*ausencia.LeaveRequest, error) {
	query := `
		SELECT id, employee_id, fecha_inicio, fecha_fin, tipo, estado, motivo, aprobado_por, created_at, updated_at
		FROM leave_requests
		WHERE fecha_inicio <= $2 AND fecha_fin >= $1 AND estado = 'APROBADO'
		ORDER BY employee_id, fecha_inicio
	`
	return r.queryLeaveRequests(ctx, query, fechaInicio, fechaFin)
}

func (r *LeaveRequestRepository) FindApprovedByEmployeeAndDate(ctx context.Context, employeeID string, fecha time.Time) ([]*ausencia.LeaveRequest, error) {
	query := `
		SELECT id, employee_id, fecha_inicio, fecha_fin, tipo, estado, motivo, aprobado_por, created_at, updated_at
		FROM leave_requests
		WHERE employee_id = $1 AND $2 BETWEEN fecha_inicio AND fecha_fin AND estado = 'APROBADO'
		ORDER BY created_at DESC
	`
	return r.queryLeaveRequests(ctx, query, employeeID, fecha)
}

func (r *LeaveRequestRepository) Update(ctx context.Context, lr *ausencia.LeaveRequest) error {
	query := `
		UPDATE leave_requests
		SET estado = $1, aprobado_por = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, string(lr.Estado), lr.AprobadoPor, lr.UpdatedAt, lr.ID)
	return err
}

func (r *LeaveRequestRepository) queryLeaveRequests(ctx context.Context, query string, args ...any) ([]*ausencia.LeaveRequest, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*ausencia.LeaveRequest
	for rows.Next() {
		lr, err := scanLeaveRequest(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, lr)
	}
	return result, rows.Err()
}

type leaveRequestScanner interface {
	Scan(dest ...any) error
}

func scanLeaveRequest(s leaveRequestScanner) (*ausencia.LeaveRequest, error) {
	var (
		id, employeeID, tipo, estado, motivo string
		fechaInicio, fechaFin, createdAt, updatedAt time.Time
		aprobadoPor                          *string
	)
	err := s.Scan(&id, &employeeID, &fechaInicio, &fechaFin, &tipo, &estado, &motivo, &aprobadoPor, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("leave request not found: %w", err)
		}
		return nil, err
	}

	t, _ := ausencia.ParseTipoAusencia(tipo)
	e := ausencia.EstadoAusencia(estado)

	return &ausencia.LeaveRequest{
		ID:          id,
		EmployeeID:  employeeID,
		FechaInicio: fechaInicio,
		FechaFin:    fechaFin,
		Tipo:        t,
		Estado:      e,
		Motivo:      motivo,
		AprobadoPor: aprobadoPor,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
