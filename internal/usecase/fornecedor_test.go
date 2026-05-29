package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockFornecedorRepo struct {
	fornecedores []entity.DimFornecedor
	err          error
}

func (m *mockFornecedorRepo) FindAll() ([]entity.DimFornecedor, error) {
	return m.fornecedores, m.err
}

// Teste 1: sucesso
func TestFornecedorUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimFornecedor{
		{
			IdFornecedor:     "1",
			NomeFornecedor:   "Fornecedor Alpha",
			CnpjFornecedor:   "12345678000199",
			CidadeFornecedor: "São Paulo",
			StatusFornecedor: "Ativo",
		},
		{
			IdFornecedor:     "2",
			NomeFornecedor:   "Fornecedor Beta",
			CnpjFornecedor:   "98765432000188",
			CidadeFornecedor: "Campinas",
			StatusFornecedor: "Inativo",
		},
	}

	mock := &mockFornecedorRepo{
		fornecedores: fake,
		err:          nil,
	}

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

	if result[0].NomeFornecedor != "Fornecedor Alpha" {
		t.Errorf("esperava fornecedor 'Fornecedor Alpha', recebeu %s", result[0].NomeFornecedor)
	}

	if result[1].CidadeFornecedor != "Campinas" {
		t.Errorf("esperava cidade 'Campinas', recebeu %s", result[1].CidadeFornecedor)
	}
}

// Teste 2: erro do repository
func TestFornecedorUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFornecedorRepo{
		fornecedores: nil,
		err:          errors.New("database connection failed"),
	}

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

// Teste 3: lista vazia
func TestFornecedorUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFornecedorRepo{
		fornecedores: []entity.DimFornecedor{},
		err:          nil,
	}

	uc := NewFornecedorUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("esperava lista vazia, recebeu %d itens", len(result))
	}
}