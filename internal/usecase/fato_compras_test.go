package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockFatoComprasRepo struct {
	data []entity.FatoCompras
	err  error
}

func (m *mockFatoComprasRepo) FindAll() ([]entity.FatoCompras, error) {
	return m.data, m.err
}

func TestFatoComprasUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.FatoCompras{
		{SkFato: "1", SkProjeto: "P1", ValorTotalPedido: "100.0"},
		{SkFato: "2", SkProjeto: "P2", ValorTotalPedido: "200.0"},
	}
	mock := &mockFatoComprasRepo{data: fake, err: nil}
	uc := NewFatoComprasUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 compras, recebeu: %d", len(result))
	}
	if result[0].SkProjeto != "P1" {
		t.Errorf("esperava P1, recebeu: %s", result[0].SkProjeto)
	}
}

func TestFatoComprasUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFatoComprasRepo{data: nil, err: errors.New("db error")}
	uc := NewFatoComprasUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if result != nil {
		t.Errorf("esperava nil, recebeu: %v", result)
	}
}

func TestFatoComprasUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFatoComprasRepo{data: []entity.FatoCompras{}, err: nil}
	uc := NewFatoComprasUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 compras, recebeu: %d", len(result))
	}
}
