package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuusuario/nursery-portal/internal/domain/intercambio"
)

type IntercambioRepository struct {
	pool *pgxpool.Pool
}

func NewIntercambioRepository(pool *pgxpool.Pool) *IntercambioRepository {
	return &IntercambioRepository{pool: pool}
}

func (r *IntercambioRepository) Create(ctx context.Context, req *intercambio.ShiftSwapRequest) error {
	query := `
		INSERT INTO shift_swap_requests (id, planificacion_id, turno_solicitante_id, turno_destino_id, solicitante_id, destino_id, estado, aprobado_por, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		req.ID, req.PlanificacionID, req.TurnoSolicitanteID, req.TurnoDestinoID,
		req.SolicitanteID, req.DestinoID, string(req.Estado), req.AprobadoPor,
		req.CreatedAt, req.UpdatedAt,
	)
	return err
}

func (r *IntercambioRepository) FindByID(ctx context.Context, id string) (*intercambio.ShiftSwapRequest, error) {
	query := `
		SELECT id, planificacion_id, turno_solicitante_id, turno_destino_id, solicitante_id, destino_id, estado, aprobado_por, created_at, updated_at
		FROM shift_swap_requests
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	return scanSwapRequest(row)
}

func (r *IntercambioRepository) FindAll(ctx context.Context) ([]*intercambio.ShiftSwapRequest, error) {
	query := `
		SELECT id, planificacion_id, turno_solicitante_id, turno_destino_id, solicitante_id, destino_id, estado, aprobado_por, created_at, updated_at
		FROM shift_swap_requests
		ORDER BY created_at DESC
	`
	return r.querySwapRequests(ctx, query)
}

func (r *IntercambioRepository) Update(ctx context.Context, req *intercambio.ShiftSwapRequest) error {
	query := `
		UPDATE shift_swap_requests
		SET estado = $1, aprobado_por = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, string(req.Estado), req.AprobadoPor, req.UpdatedAt, req.ID)
	return err
}

func (r *IntercambioRepository) AddHistoryEntry(ctx context.Context, entry *intercambio.ShiftSwapHistoryEntry) error {
	query := `
		INSERT INTO shift_swap_history (id, swap_request_id, accion, actor_id, detalle, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		entry.ID, entry.SwapRequestID, string(entry.Accion), entry.ActorID,
		entry.Detalle, entry.CreatedAt,
	)
	return err
}

func (r *IntercambioRepository) GetHistory(ctx context.Context, swapRequestID string) ([]*intercambio.ShiftSwapHistoryEntry, error) {
	query := `
		SELECT id, swap_request_id, accion, actor_id, detalle, created_at
		FROM shift_swap_history
		WHERE swap_request_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, swapRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*intercambio.ShiftSwapHistoryEntry
	for rows.Next() {
		var (
			id, swapRequestID, accionStr, actorID string
			detalle                              *string
			createdAt                            time.Time
		)
		if err := rows.Scan(&id, &swapRequestID, &accionStr, &actorID, &detalle, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, &intercambio.ShiftSwapHistoryEntry{
			ID:            id,
			SwapRequestID: swapRequestID,
			Accion:        intercambio.AccionHistorial(accionStr),
			ActorID:       actorID,
			Detalle:       detalle,
			CreatedAt:     createdAt,
		})
	}
	return result, rows.Err()
}

func (r *IntercambioRepository) querySwapRequests(ctx context.Context, query string, args ...any) ([]*intercambio.ShiftSwapRequest, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*intercambio.ShiftSwapRequest
	for rows.Next() {
		req, err := scanSwapRequest(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, req)
	}
	return result, rows.Err()
}

type swapRequestScanner interface {
	Scan(dest ...any) error
}

func scanSwapRequest(s swapRequestScanner) (*intercambio.ShiftSwapRequest, error) {
	var (
		id, planificacionID, turnoSolID, turnoDestID, solicitanteID, destinoID, estado string
		aprobadoPor                                                                     *string
		createdAt, updatedAt                                                           time.Time
	)
	err := s.Scan(&id, &planificacionID, &turnoSolID, &turnoDestID, &solicitanteID, &destinoID, &estado, &aprobadoPor, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("swap request not found: %w", err)
		}
		return nil, err
	}

	return &intercambio.ShiftSwapRequest{
		ID:                 id,
		PlanificacionID:    planificacionID,
		TurnoSolicitanteID: turnoSolID,
		TurnoDestinoID:     turnoDestID,
		SolicitanteID:      solicitanteID,
		DestinoID:          destinoID,
		Estado:             intercambio.EstadoSwap(estado),
		AprobadoPor:        aprobadoPor,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}, nil
}
