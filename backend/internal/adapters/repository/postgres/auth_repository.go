package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuusuario/nursery-portal/internal/domain/auth"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (*auth.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, username, password_hash, role, employee_id, must_change_password, created_at, updated_at
		FROM auth_users
		WHERE username = $1
	`, username)
	return scanAuthUser(row)
}

func (r *AuthRepository) FindUserBySessionHash(ctx context.Context, tokenHash string, now time.Time) (*auth.User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT u.id, u.username, u.password_hash, u.role, u.employee_id, u.must_change_password, u.created_at, u.updated_at
		FROM auth_sessions s
		JOIN auth_users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > $2
	`, tokenHash, now)
	return scanAuthUser(row)
}

func (r *AuthRepository) SetPasswordHash(ctx context.Context, userID, passwordHash string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE auth_users
		SET password_hash = $1, must_change_password = false, updated_at = NOW()
		WHERE id = $2 AND (password_hash IS NULL OR must_change_password = true)
	`, passwordHash, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("password already set")
	}
	return nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth_sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *AuthRepository) DeleteSession(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM auth_sessions WHERE token_hash = $1`, tokenHash)
	return err
}

func (r *AuthRepository) EnsureAdmin(ctx context.Context, username, passwordHash string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth_users (username, password_hash, role)
		VALUES ($1, $2, 'ADMIN')
		ON CONFLICT (username) DO UPDATE
		SET password_hash = EXCLUDED.password_hash,
		    role = 'ADMIN',
		    employee_id = NULL,
		    updated_at = NOW()
	`, username, passwordHash)
	return err
}

func (r *AuthRepository) CreateEmployeeUser(ctx context.Context, username string, role auth.Role, employeeID string, passwordHash string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO auth_users (username, role, employee_id, password_hash, must_change_password)
		VALUES ($1, $2, $3, $4, true)
	`, username, string(role), employeeID, passwordHash)
	if isUniqueViolation(err) {
		return errors.New("username already exists")
	}
	return err
}

func (r *AuthRepository) UpdateEmployeeUser(ctx context.Context, employeeID, username string, role auth.Role) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE auth_users
		SET username = $1, role = $2, updated_at = NOW()
		WHERE employee_id = $3
	`, username, string(role), employeeID)
	if isUniqueViolation(err) {
		return errors.New("username already exists")
	}
	if err == nil && tag.RowsAffected() == 0 {
		return r.CreateEmployeeUser(ctx, username, role, employeeID, "")
	}
	return err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func scanAuthUser(s scanner) (*auth.User, error) {
	var (
		id, username, role string
		passwordHash       *string
		employeeID         *string
		mustChangePassword bool
		createdAt          time.Time
		updatedAt          time.Time
	)
	if err := s.Scan(&id, &username, &passwordHash, &role, &employeeID, &mustChangePassword, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("auth user not found: %w", err)
		}
		return nil, err
	}
	return &auth.User{
		ID:                 id,
		Username:           username,
		PasswordHash:       passwordHash,
		Role:               auth.Role(role),
		EmployeeID:         employeeID,
		MustChangePassword: mustChangePassword,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}, nil
}
