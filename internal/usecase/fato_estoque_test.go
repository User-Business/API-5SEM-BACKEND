package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockFatoEstoqueRepo struct {
	data []entity.FatoEstoqueMateriais
	err  error
}

func (m *mockFatoEstoqueRepo) FindAll() ([]entity.FatoEstoqueMateriais, error) {
	return m.data, m.err
}

func TestFatoEstoqueUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.FatoEstoqueMateriais{
		{SkFato: "1", SkProjeto: "P1", QuantidadeEstoque: "10"},
		{SkFato: "2", SkProjeto: "P2", QuantidadeEstoque: "20"},
	}
	mock := &mockFatoEstoqueRepo{data: fake, err: nil}
	uc := NewFatoEstoqueUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 itens de estoque, recebeu: %d", len(result))
	}
	if result[0].SkProjeto != "P1" {
		t.Errorf("esperava P1, recebeu: %s", result[0].SkProjeto)
	}
}

func TestFatoEstoqueUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFatoEstoqueRepo{data: nil, err: errors.New("db error")}
	uc := NewFatoEstoqueUseCase(mock)

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

func TestFatoEstoqueUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFatoEstoqueRepo{data: []entity.FatoEstoqueMateriais{}, err: nil}
	uc := NewFatoEstoqueUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 itens de estoque, recebeu: %d", len(result))
	}
}
