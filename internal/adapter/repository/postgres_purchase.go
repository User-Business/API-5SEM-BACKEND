package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type PostgresPurchaseRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPurchaseRepository(db *pgxpool.Pool) *PostgresPurchaseRepository {
	return &PostgresPurchaseRepository{db: db}
}

func (r *PostgresPurchaseRepository) FindPurchases(ctx context.Context, filter entity.PurchaseFilter) ([]entity.Purchase, error) {
	query := `
		WITH combined AS (
			SELECT 
				id_solicitacao AS id, 
				'SC' AS type, 
				numero_solicitacao AS numero, 
				status, 
				data_solicitacao::text AS data_criacao,
				NULL AS data_previsao_entrega,
				NULL::int AS duracao_dias,
				NULL::boolean AS atrasado
			FROM solicitacao_compra

			UNION ALL

			SELECT 
				id_pedido AS id, 
				'PC' AS type, 
				numero_pedido AS numero, 
				status, 
				data_pedido::text AS data_criacao,
				data_previsao_entrega::text AS data_previsao_entrega,
				CASE 
					WHEN status IN ('Concluído', 'Entregue', 'Cancelado', 'Concluido') THEN data_previsao_entrega - data_pedido
					ELSE CURRENT_DATE - data_pedido
				END AS duracao_dias,
				CASE 
					WHEN status NOT IN ('Concluído', 'Entregue', 'Cancelado', 'Concluido') AND CURRENT_DATE > data_previsao_entrega THEN TRUE
					ELSE FALSE
				END AS atrasado
			FROM pedido_compra
		)
		SELECT id, type, numero, COALESCE(status, ''), COALESCE(data_criacao, ''), data_previsao_entrega, duracao_dias, atrasado 
		FROM combined
		WHERE 1=1
	`

	args := []interface{}{}
	argId := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argId)
		args = append(args, filter.Type)
		argId++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argId)
		args = append(args, filter.Status)
		argId++
	}
	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND data_criacao >= $%d", argId)
		args = append(args, filter.StartDate)
		argId++
	}
	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND data_criacao <= $%d", argId)
		args = append(args, filter.EndDate)
		argId++
	}

	query += " ORDER BY data_criacao DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entity.Purchase
	for rows.Next() {
		var p entity.Purchase
		err := rows.Scan(
			&p.ID, &p.Type, &p.Number, &p.Status, &p.Date,
			&p.ExpectedDeliveryDate, &p.DurationDays, &p.IsDelayed,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, p)
	}

	return results, rows.Err()
}

func (r *PostgresPurchaseRepository) GetMetrics(ctx context.Context, filter entity.PurchaseFilter) (entity.PurchaseMetrics, error) {
	purchases, err := r.FindPurchases(ctx, filter)
	if err != nil {
		return entity.PurchaseMetrics{}, err
	}

	metrics := entity.PurchaseMetrics{
		TotalPurchases: len(purchases),
		StatusCounts:   []entity.StatusCount{},
	}

	statusMap := make(map[string]int)
	var totalDurationPC int
	var countPC int

	for _, p := range purchases {
		statusMap[p.Status]++
		if p.Type == "SC" {
			metrics.TotalSC++
		} else if p.Type == "PC" {
			metrics.TotalPC++
			if p.DurationDays != nil {
				totalDurationPC += *p.DurationDays
				countPC++
			}
		}
	}

	if countPC > 0 {
		metrics.AverageDurationPC = float64(totalDurationPC) / float64(countPC)
	}

	for status, count := range statusMap {
		metrics.StatusCounts = append(metrics.StatusCounts, entity.StatusCount{
			Status: status,
			Count:  count,
		})
	}

	return metrics, nil
}
