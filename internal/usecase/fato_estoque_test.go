package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockFatoEstoqueRepo struct {
	estoques []entity.FatoEstoqueMateriais
	err      error
}

func (m *mockFatoEstoqueRepo) FindAll() ([]entity.FatoEstoqueMateriais, error) {
	return m.estoques, m.err
}

// Teste 1: sucesso
func TestFatoEstoqueUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.FatoEstoqueMateriais{
		{
			IdMaterial:        "1",
			CodigoMaterial:    "MAT001",
			DescricaoMaterial: "Material A",
			QuantidadeEstoque: 100,
			StatusEstoque:     "Disponível",
		},
		{
			IdMaterial:        "2",
			CodigoMaterial:    "MAT002",
			DescricaoMaterial: "Material B",
			QuantidadeEstoque: 50,
			StatusEstoque:     "Baixo",
		},
	}

	mock := &mockFatoEstoqueRepo{
		estoques: fake,
		err:      nil,
	}

	uc := NewFatoEstoqueUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 registros, recebeu: %d", len(result))
	}

	if result[0].CodigoMaterial != "MAT001" {
		t.Errorf("esperava codigo MAT001, recebeu %s", result[0].CodigoMaterial)
	}

	if result[1].DescricaoMaterial != "Material B" {
		t.Errorf("esperava descricao 'Material B', recebeu %s", result[1].DescricaoMaterial)
	}
}

// Teste 2: erro do repository
func TestFatoEstoqueUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFatoEstoqueRepo{
		estoques: nil,
		err:      errors.New("database connection failed"),
	}

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

// Teste 3: lista vazia
func TestFatoEstoqueUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFatoEstoqueRepo{
		estoques: []entity.FatoEstoqueMateriais{},
		err:      nil,
	}

	uc := NewFatoEstoqueUseCase(mock)

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