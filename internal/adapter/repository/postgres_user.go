package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// FindByEmail retorna o usuário pelo email, ou (nil, nil) se não existir.
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	query := `SELECT id, email, senha_hash, nome, role, criado_em FROM usuario WHERE email = $1`

	var u entity.Usuario
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.SenhaHash, &u.Nome, &u.Role, &u.CriadoEm,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
