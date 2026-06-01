package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/middleware"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/usecase"
)

type AuthHandler struct {
	useCase *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: uc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"payload inválido"}`, http.StatusBadRequest)
		return
	}

	res, err := h.useCase.Login(r.Context(), req.Email, req.Senha)
	if errors.Is(err, usecase.ErrCredenciaisInvalidas) {
		http.Error(w, `{"error":"credenciais inválidas"}`, http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"erro interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"não autorizado"}`, http.StatusUnauthorized)
		return
	}

	var exp time.Time
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Time
	}

	if err := h.useCase.Logout(r.Context(), claims.ID, exp); err != nil {
		http.Error(w, `{"error":"erro interno"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
