package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuusuario/nursery-portal/internal/domain/planificacion"
	"github.com/google/uuid"
)

type DotacionRepository struct {
	pool *pgxpool.Pool
}

func NewDotacionRepository(pool *pgxpool.Pool) *DotacionRepository {
	return &DotacionRepository{pool: pool}
}

func (r *DotacionRepository) GetSectores(ctx context.Context, planificacionID string) ([]*planificacion.SectorPlanificacion, error) {
	query := `
		SELECT id, planificacion_id, nombre, created_at
		FROM planificacion_sectores
		WHERE planificacion_id = $1
		ORDER BY nombre
	`
	rows, err := r.pool.Query(ctx, query, planificacionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectores []*planificacion.SectorPlanificacion
	for rows.Next() {
		s, err := scanSectorPlanificacion(rows)
		if err != nil {
			return nil, err
		}
		sectores = append(sectores, s)
	}
	return sectores, rows.Err()
}

func (r *DotacionRepository) SaveSectores(ctx context.Context, planificacionID string, nombres []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM planificacion_sectores WHERE planificacion_id = $1`, planificacionID); err != nil {
		return err
	}

	for _, nombre := range nombres {
		if _, err := tx.Exec(ctx,
			`INSERT INTO planificacion_sectores (id, planificacion_id, nombre, created_at) VALUES ($1, $2, $3, $4)`,
			uuid.New().String(), planificacionID, nombre, time.Now().UTC(),
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *DotacionRepository) GetDotacion(ctx context.Context, planificacionID string) ([]*planificacion.DotacionPlanificacion, error) {
	query := `
		SELECT id, planificacion_id, sector, tipo_empleado, turno, cantidad_minima, created_at, updated_at
		FROM planificacion_dotacion
		WHERE planificacion_id = $1
		ORDER BY sector, tipo_empleado, turno
	`
	rows, err := r.pool.Query(ctx, query, planificacionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*planificacion.DotacionPlanificacion
	for rows.Next() {
		d, err := scanDotacionPlanificacion(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}

func (r *DotacionRepository) SaveDotacion(ctx context.Context, items []*planificacion.DotacionPlanificacion) error {
	if len(items) == 0 {
		return nil
	}

	planificacionID := items[0].PlanificacionID

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM planificacion_dotacion WHERE planificacion_id = $1`, planificacionID); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, d := range items {
		if d.ID == "" {
			d.ID = uuid.New().String()
		}
		if d.CreatedAt.IsZero() {
			d.CreatedAt = now
		}
		d.UpdatedAt = now

		if _, err := tx.Exec(ctx,
			`INSERT INTO planificacion_dotacion (id, planificacion_id, sector, tipo_empleado, turno, cantidad_minima, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			 ON CONFLICT (planificacion_id, sector, tipo_empleado, turno)
			 DO UPDATE SET cantidad_minima = $6, updated_at = $8`,
			d.ID, d.PlanificacionID, d.Sector, d.TipoEmpleado, d.Turno, d.CantidadMinima, d.CreatedAt, d.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *DotacionRepository) DeleteByPlanificacion(ctx context.Context, planificacionID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM planificacion_dotacion WHERE planificacion_id = $1`, planificacionID)
	return err
}

func scanSectorPlanificacion(s scanner) (*planificacion.SectorPlanificacion, error) {
	var id, planifID, nombre string
	var createdAt time.Time
	err := s.Scan(&id, &planifID, &nombre, &createdAt)
	if err != nil {
		return nil, err
	}
	return &planificacion.SectorPlanificacion{
		ID:              id,
		PlanificacionID: planifID,
		Nombre:          nombre,
		CreatedAt:       createdAt,
	}, nil
}

func scanDotacionPlanificacion(s scanner) (*planificacion.DotacionPlanificacion, error) {
	var id, planifID, sector, tipoEmp, turno string
	var cantidadMinima int
	var createdAt, updatedAt time.Time
	err := s.Scan(&id, &planifID, &sector, &tipoEmp, &turno, &cantidadMinima, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &planificacion.DotacionPlanificacion{
		ID:              id,
		PlanificacionID: planifID,
		Sector:          sector,
		TipoEmpleado:    tipoEmp,
		Turno:           turno,
		CantidadMinima:  cantidadMinima,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}, nil
}
