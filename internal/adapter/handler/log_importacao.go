package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/usecase"
)

type LogImportacaoHandler struct {
	useCase *usecase.LogImportacaoUseCase
}

func NewLogImportacaoHandler(uc *usecase.LogImportacaoUseCase) *LogImportacaoHandler {
	return &LogImportacaoHandler{useCase: uc}
}

func (h *LogImportacaoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var logEntity entity.LogImportacao
	if err := json.NewDecoder(r.Body).Decode(&logEntity); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	id, err := h.useCase.Create(&logEntity)
	if err != nil {
		http.Error(w, `{"error":"failed to create import log"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "log": logEntity})
}

func (h *LogImportacaoHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid log ID"}`, http.StatusBadRequest)
		return
	}

	var logEntity entity.LogImportacao
	if err := json.NewDecoder(r.Body).Decode(&logEntity); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	logEntity.ID = id

	if err := h.useCase.Update(&logEntity); err != nil {
		http.Error(w, `{"error":"failed to update import log"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntity)
}

func (h *LogImportacaoHandler) CreateError(w http.ResponseWriter, r *http.Request) {
	logIDStr := chi.URLParam(r, "id")
	logID, err := strconv.Atoi(logIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid log ID"}`, http.StatusBadRequest)
		return
	}

	var errLog entity.LogImportacaoErro
	if err := json.NewDecoder(r.Body).Decode(&errLog); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	errLog.LogImportacaoID = logID

	if err := h.useCase.CreateError(&errLog); err != nil {
		http.Error(w, `{"error":"failed to create error log"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(errLog)
}

func (h *LogImportacaoHandler) CreateErrorsBatch(w http.ResponseWriter, r *http.Request) {
	logIDStr := chi.URLParam(r, "id")
	logID, err := strconv.Atoi(logIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid log ID"}`, http.StatusBadRequest)
		return
	}

	var errLogs []entity.LogImportacaoErro
	if err := json.NewDecoder(r.Body).Decode(&errLogs); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	for i := range errLogs {
		errLogs[i].LogImportacaoID = logID
	}

	if err := h.useCase.CreateErrorsBatch(errLogs); err != nil {
		http.Error(w, `{"error":"failed to batch create error logs"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "error logs inserted successfully"})
}

func (h *LogImportacaoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	logs, err := h.useCase.GetAll()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch import logs"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *LogImportacaoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid log ID"}`, http.StatusBadRequest)
		return
	}

	logEntity, err := h.useCase.GetByID(id)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch import log"}`, http.StatusInternalServerError)
		return
	}
	if logEntity == nil {
		http.Error(w, `{"error":"import log not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntity)
}

func (h *LogImportacaoHandler) GetErrors(w http.ResponseWriter, r *http.Request) {
	logIDStr := chi.URLParam(r, "id")
	logID, err := strconv.Atoi(logIDStr)
	if err != nil {
		http.Error(w, `{"error":"invalid log ID"}`, http.StatusBadRequest)
		return
	}

	errors, err := h.useCase.GetErrorsByLogID(logID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch import errors"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errors)
}
