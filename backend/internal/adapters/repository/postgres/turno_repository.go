package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuusuario/nursery-portal/internal/domain/turno"
)

type TurnoRepository struct {
	pool *pgxpool.Pool
}

func NewTurnoRepository(pool *pgxpool.Pool) *TurnoRepository {
	return &TurnoRepository{pool: pool}
}

func (r *TurnoRepository) Create(ctx context.Context, t *turno.Turno) error {
	query := `
		INSERT INTO turnos (id, planificacion_id, empleado_id, dia_semana, turno, sector, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (planificacion_id, empleado_id, dia_semana, turno) DO NOTHING
	`
	_, err := r.pool.Exec(ctx, query,
		t.ID, t.PlanificacionID, t.EmpleadoID, t.DiaSemana, string(t.Tipo), t.Sector,
		t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *TurnoRepository) CreateBatch(ctx context.Context, turnos []*turno.Turno) error {
	if len(turnos) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO turnos (id, planificacion_id, empleado_id, dia_semana, turno, sector, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (planificacion_id, empleado_id, dia_semana, turno) DO NOTHING
	`
	for _, t := range turnos {
		batch.Queue(query, t.ID, t.PlanificacionID, t.EmpleadoID, t.DiaSemana, string(t.Tipo), t.Sector, t.CreatedAt, t.UpdatedAt)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for range turnos {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return br.Close()
}

func (r *TurnoRepository) FindByPlanificacion(ctx context.Context, planificacionID string) ([]*turno.Turno, error) {
	query := `
		SELECT id, planificacion_id, empleado_id, dia_semana, turno, sector, created_at, updated_at
		FROM turnos
		WHERE planificacion_id = $1
		ORDER BY empleado_id, dia_semana
	`
	rows, err := r.pool.Query(ctx, query, planificacionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turnos []*turno.Turno
	for rows.Next() {
		t, err := scanTurno(rows)
		if err != nil {
			return nil, err
		}
		turnos = append(turnos, t)
	}
	return turnos, rows.Err()
}

func (r *TurnoRepository) FindByPlanificacionAndEmpleado(ctx context.Context, planificacionID, empleadoID string) ([]*turno.Turno, error) {
	query := `
		SELECT id, planificacion_id, empleado_id, dia_semana, turno, sector, created_at, updated_at
		FROM turnos
		WHERE planificacion_id = $1 AND empleado_id = $2
		ORDER BY dia_semana
	`
	rows, err := r.pool.Query(ctx, query, planificacionID, empleadoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turnos []*turno.Turno
	for rows.Next() {
		t, err := scanTurno(rows)
		if err != nil {
			return nil, err
		}
		turnos = append(turnos, t)
	}
	return turnos, rows.Err()
}

func (r *TurnoRepository) FindByTurnoID(ctx context.Context, id string) (*turno.Turno, error) {
	query := `
		SELECT id, planificacion_id, empleado_id, dia_semana, turno, sector, created_at, updated_at
		FROM turnos
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	return scanTurno(row)
}

func (r *TurnoRepository) Update(ctx context.Context, t *turno.Turno) error {
	query := `
		UPDATE turnos SET empleado_id = $1, updated_at = $2 WHERE id = $3
	`
	_, err := r.pool.Exec(ctx, query, t.EmpleadoID, t.UpdatedAt, t.ID)
	return err
}

func (r *TurnoRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM turnos WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *TurnoRepository) DeleteByPlanificacion(ctx context.Context, planificacionID string) error {
	query := `DELETE FROM turnos WHERE planificacion_id = $1`
	_, err := r.pool.Exec(ctx, query, planificacionID)
	return err
}

func scanTurno(s scanner) (*turno.Turno, error) {
	var (
		id, planifID, empID, tipo, sector string
		diaSemana                          int
		createdAt, updatedAt               time.Time
	)
	err := s.Scan(&id, &planifID, &empID, &diaSemana, &tipo, &sector, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return &turno.Turno{
		ID:              id,
		PlanificacionID: planifID,
		EmpleadoID:      empID,
		DiaSemana:       diaSemana,
		Tipo:            turno.TipoTurno(tipo),
		Sector:          sector,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}
