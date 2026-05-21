package handler

import (
	"encoding/json"
	"net/http"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/usecase"
)

type PurchaseHandler struct {
	useCase *usecase.PurchaseUseCase
}

func NewPurchaseHandler(uc *usecase.PurchaseUseCase) *PurchaseHandler {
	return &PurchaseHandler{useCase: uc}
}

func (h *PurchaseHandler) GetPurchases(w http.ResponseWriter, r *http.Request) {
	filter := entity.PurchaseFilter{
		Type:      r.URL.Query().Get("type"),
		Status:    r.URL.Query().Get("status"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	data, err := h.useCase.GetPurchases(r.Context(), filter)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch purchases"}`, http.StatusInternalServerError)
		return
	}

	// In case there are no results, ensure we return an empty array instead of null
	if data == nil {
		data = []entity.Purchase{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *PurchaseHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	filter := entity.PurchaseFilter{
		Type:      r.URL.Query().Get("type"),
		Status:    r.URL.Query().Get("status"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	metrics, err := h.useCase.GetMetrics(r.Context(), filter)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch purchase metrics"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
