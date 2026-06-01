package entity

// Purchase represents either a Solicitação de Compra (SC) or Pedido de Compra (PC).
type Purchase struct {
	ID                   int     `json:"id"`
	Type                 string  `json:"type"` // "SC" or "PC"
	Number               string  `json:"numero"`
	Status               string  `json:"status"`
	Date                 string  `json:"data_criacao"`
	ExpectedDeliveryDate *string `json:"data_previsao_entrega,omitempty"`
	DurationDays         *int    `json:"duracao_dias,omitempty"`
	IsDelayed            *bool   `json:"atrasado,omitempty"`
}

// PurchaseFilter defines the filter parameters for querying purchases.
type PurchaseFilter struct {
	Type      string `json:"type"` // "SC", "PC", or empty for all
	Status    string `json:"status"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// StatusCount represents the count of purchases grouped by status.
type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// PurchaseMetrics represents the aggregated metrics response.
type PurchaseMetrics struct {
	TotalPurchases    int           `json:"total_purchases"`
	TotalSC           int           `json:"total_sc"`
	TotalPC           int           `json:"total_pc"`
	AverageDurationPC float64       `json:"average_duration_pc"`
	StatusCounts      []StatusCount `json:"status_counts"`
}
