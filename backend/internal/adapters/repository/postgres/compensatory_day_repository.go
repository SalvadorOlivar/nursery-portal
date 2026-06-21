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

type CompensatoryDayRepository struct {
	pool *pgxpool.Pool
}

func NewCompensatoryDayRepository(pool *pgxpool.Pool) *CompensatoryDayRepository {
	return &CompensatoryDayRepository{pool: pool}
}

func (r *CompensatoryDayRepository) Create(ctx context.Context, cd *ausencia.CompensatoryDay) error {
	query := `
		INSERT INTO compensatory_days (id, employee_id, fecha_origen, motivo, turno_id, descripcion, utilizado, fecha_uso, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		cd.ID, cd.EmployeeID, cd.FechaOrigen, string(cd.Motivo),
		cd.TurnoID, cd.Descripcion, cd.Utilizado, cd.FechaUso, cd.CreatedAt,
	)
	return err
}

func (r *CompensatoryDayRepository) FindByEmployee(ctx context.Context, employeeID string) ([]*ausencia.CompensatoryDay, error) {
	query := `
		SELECT id, employee_id, fecha_origen, motivo, turno_id, descripcion, utilizado, fecha_uso, created_at
		FROM compensatory_days
		WHERE employee_id = $1
		ORDER BY created_at DESC
	`
	return r.queryCompensatoryDays(ctx, query, employeeID)
}

func (r *CompensatoryDayRepository) FindByID(ctx context.Context, id string) (*ausencia.CompensatoryDay, error) {
	query := `
		SELECT id, employee_id, fecha_origen, motivo, turno_id, descripcion, utilizado, fecha_uso, created_at
		FROM compensatory_days
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	return scanCompensatoryDay(row)
}

func (r *CompensatoryDayRepository) FindAvailableByEmployee(ctx context.Context, employeeID string) ([]*ausencia.CompensatoryDay, error) {
	query := `
		SELECT id, employee_id, fecha_origen, motivo, turno_id, descripcion, utilizado, fecha_uso, created_at
		FROM compensatory_days
		WHERE employee_id = $1 AND utilizado = false
		ORDER BY created_at ASC
	`
	return r.queryCompensatoryDays(ctx, query, employeeID)
}

func (r *CompensatoryDayRepository) Update(ctx context.Context, cd *ausencia.CompensatoryDay) error {
	query := `
		UPDATE compensatory_days
		SET utilizado = $1, fecha_uso = $2
		WHERE id = $3
	`
	_, err := r.pool.Exec(ctx, query, cd.Utilizado, cd.FechaUso, cd.ID)
	return err
}

func (r *CompensatoryDayRepository) CountAvailable(ctx context.Context, employeeID string) (int, error) {
	query := `SELECT COUNT(*) FROM compensatory_days WHERE employee_id = $1 AND utilizado = false`
	var count int
	err := r.pool.QueryRow(ctx, query, employeeID).Scan(&count)
	return count, err
}

func (r *CompensatoryDayRepository) queryCompensatoryDays(ctx context.Context, query string, args ...any) ([]*ausencia.CompensatoryDay, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*ausencia.CompensatoryDay
	for rows.Next() {
		cd, err := scanCompensatoryDay(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, cd)
	}
	return result, rows.Err()
}

type compensatoryDayScanner interface {
	Scan(dest ...any) error
}

func scanCompensatoryDay(s compensatoryDayScanner) (*ausencia.CompensatoryDay, error) {
	var (
		id, employeeID, motivo, descripcion string
		fechaOrigen, createdAt              time.Time
		turnoID                             *string
		utilizado                           bool
		fechaUso                            *time.Time
	)
	err := s.Scan(&id, &employeeID, &fechaOrigen, &motivo, &turnoID, &descripcion, &utilizado, &fechaUso, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("compensatory day not found: %w", err)
		}
		return nil, err
	}

	return &ausencia.CompensatoryDay{
		ID:          id,
		EmployeeID:  employeeID,
		FechaOrigen: fechaOrigen,
		Motivo:      ausencia.MotivoCompensatorio(motivo),
		TurnoID:     turnoID,
		Descripcion: descripcion,
		Utilizado:   utilizado,
		FechaUso:    fechaUso,
		CreatedAt:   createdAt,
	}, nil
}
