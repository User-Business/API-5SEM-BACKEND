package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockFornecedorRepo struct {
	data []entity.DimFornecedor
	err  error
}

func (m *mockFornecedorRepo) FindAll() ([]entity.DimFornecedor, error) {
	return m.data, m.err
}

func TestFornecedorUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimFornecedor{
		{SkFornecedor: "1", IdFornecedor: "F1", RazaoSocial: "Fornecedor Alpha"},
		{SkFornecedor: "2", IdFornecedor: "F2", RazaoSocial: "Fornecedor Beta"},
	}
	mock := &mockFornecedorRepo{data: fake, err: nil}
	uc := NewFornecedorUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 fornecedores, recebeu: %d", len(result))
	}
	if result[0].RazaoSocial != "Fornecedor Alpha" {
		t.Errorf("esperava Fornecedor Alpha, recebeu: %s", result[0].RazaoSocial)
	}
}

func TestFornecedorUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFornecedorRepo{data: nil, err: errors.New("db error")}
	uc := NewFornecedorUseCase(mock)

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

func TestFornecedorUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFornecedorRepo{data: []entity.DimFornecedor{}, err: nil}
	uc := NewFornecedorUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 fornecedores, recebeu: %d", len(result))
	}
}
