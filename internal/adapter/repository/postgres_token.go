package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTokenRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTokenRepository(db *pgxpool.Pool) *PostgresTokenRepository {
	return &PostgresTokenRepository{db: db}
}

// Revoke adiciona o jti à denylist com seu instante de expiração.
func (r *PostgresTokenRepository) Revoke(ctx context.Context, jti string, exp time.Time) error {
	query := `INSERT INTO token_revogado (jti, expira_em) VALUES ($1, $2) ON CONFLICT (jti) DO NOTHING`
	_, err := r.db.Exec(ctx, query, jti, exp)
	return err
}

// IsRevoked indica se o jti está na denylist.
func (r *PostgresTokenRepository) IsRevoked(ctx context.Context, jti string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM token_revogado WHERE jti = $1)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, jti).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
