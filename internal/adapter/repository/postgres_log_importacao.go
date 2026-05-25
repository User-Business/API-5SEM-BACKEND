package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type PostgresLogImportacaoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresLogImportacaoRepository(db *pgxpool.Pool) *PostgresLogImportacaoRepository {
	return &PostgresLogImportacaoRepository{db: db}
}

func (r *PostgresLogImportacaoRepository) Create(log *entity.LogImportacao) (int, error) {
	query := `
		INSERT INTO log_importacao (nome_arquivo, data_inicio, data_fim, usuario_id, status, total_linhas_processadas, total_linhas_erro)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(context.Background(), query,
		log.NomeArquivo, log.DataInicio, log.DataFim, log.UsuarioID, log.Status, log.TotalLinhasProcessadas, log.TotalLinhasErro,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	log.ID = id
	return id, nil
}

func (r *PostgresLogImportacaoRepository) Update(log *entity.LogImportacao) error {
	query := `
		UPDATE log_importacao
		SET nome_arquivo = $1, data_inicio = $2, data_fim = $3, usuario_id = $4, status = $5, total_linhas_processadas = $6, total_linhas_erro = $7
		WHERE id = $8
	`
	_, err := r.db.Exec(context.Background(), query,
		log.NomeArquivo, log.DataInicio, log.DataFim, log.UsuarioID, log.Status, log.TotalLinhasProcessadas, log.TotalLinhasErro, log.ID,
	)
	return err
}

func (r *PostgresLogImportacaoRepository) CreateError(errLog *entity.LogImportacaoErro) error {
	query := `
		INSERT INTO log_importacao_erro (log_importacao_id, numero_linha, conteudo_original, motivo_erro)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(context.Background(), query,
		errLog.LogImportacaoID, errLog.NumeroLinha, errLog.ConteudoOriginal, errLog.MotivoErro,
	).Scan(&id)
	if err != nil {
		return err
	}
	errLog.ID = id
	return nil
}

func (r *PostgresLogImportacaoRepository) CreateErrorsBatch(errLogs []entity.LogImportacaoErro) error {
	if len(errLogs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, errLog := range errLogs {
		batch.Queue(
			"INSERT INTO log_importacao_erro (log_importacao_id, numero_linha, conteudo_original, motivo_erro) VALUES ($1, $2, $3, $4)",
			errLog.LogImportacaoID, errLog.NumeroLinha, errLog.ConteudoOriginal, errLog.MotivoErro,
		)
	}

	br := r.db.SendBatch(context.Background(), batch)
	defer br.Close()

	for i := 0; i < len(errLogs); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch insert at index %d: %w", i, err)
		}
	}
	return nil
}

func (r *PostgresLogImportacaoRepository) FindAll() ([]entity.LogImportacao, error) {
	query := `
		SELECT id, nome_arquivo, data_inicio, data_fim, usuario_id, status, total_linhas_processadas, total_linhas_erro
		FROM log_importacao
		ORDER BY data_inicio DESC
	`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entity.LogImportacao
	for rows.Next() {
		var log entity.LogImportacao
		err := rows.Scan(
			&log.ID, &log.NomeArquivo, &log.DataInicio, &log.DataFim, &log.UsuarioID, &log.Status, &log.TotalLinhasProcessadas, &log.TotalLinhasErro,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, log)
	}
	return results, rows.Err()
}

func (r *PostgresLogImportacaoRepository) FindByID(id int) (*entity.LogImportacao, error) {
	query := `
		SELECT id, nome_arquivo, data_inicio, data_fim, usuario_id, status, total_linhas_processadas, total_linhas_erro
		FROM log_importacao
		WHERE id = $1
	`
	var log entity.LogImportacao
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&log.ID, &log.NomeArquivo, &log.DataInicio, &log.DataFim, &log.UsuarioID, &log.Status, &log.TotalLinhasProcessadas, &log.TotalLinhasErro,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

func (r *PostgresLogImportacaoRepository) FindErrorsByLogID(logID int) ([]entity.LogImportacaoErro, error) {
	query := `
		SELECT id, log_importacao_id, numero_linha, conteudo_original, motivo_erro
		FROM log_importacao_erro
		WHERE log_importacao_id = $1
		ORDER BY numero_linha ASC
	`
	rows, err := r.db.Query(context.Background(), query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entity.LogImportacaoErro
	for rows.Next() {
		var logErr entity.LogImportacaoErro
		err := rows.Scan(
			&logErr.ID, &logErr.LogImportacaoID, &logErr.NumeroLinha, &logErr.ConteudoOriginal, &logErr.MotivoErro,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, logErr)
	}
	return results, rows.Err()
}
