package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockMaterialRepo struct {
	materials []entity.DimMaterial
	err       error
}

func (m *mockMaterialRepo) FindAll() ([]entity.DimMaterial, error) {
	return m.materials, m.err
}

// Teste 1: sucesso
func TestMaterialUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimMaterial{
		{
			SkMaterial:     "1",
			IdMaterial:     "1",
			CodigoMaterial: "MAT001",
			Descricao:      "Material Teste 1",
			Categoria:      "Categoria A",
			Fabricante:     "Fabricante X",
			CustoEstimado:  "100.50",
			Status:         "Ativo",
		},
		{
			SkMaterial:     "2",
			IdMaterial:     "2",
			CodigoMaterial: "MAT002",
			Descricao:      "Material Teste 2",
			Categoria:      "Categoria B",
			Fabricante:     "Fabricante Y",
			CustoEstimado:  "200.75",
			Status:         "Ativo",
		},
	}

	mock := &mockMaterialRepo{
		materials: fake,
		err:       nil,
	}

	uc := NewMaterialUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 materiais, recebeu: %d", len(result))
	}

	if result[0].CodigoMaterial != "MAT001" {
		t.Errorf("esperava codigo MAT001, recebeu %s", result[0].CodigoMaterial)
	}

	if result[1].Descricao != "Material Teste 2" {
		t.Errorf("esperava descricao 'Material Teste 2', recebeu %s", result[1].Descricao)
	}
}

// Teste 2: erro do repository
func TestMaterialUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockMaterialRepo{
		materials: nil,
		err:       errors.New("database connection failed"),
	}

	uc := NewMaterialUseCase(mock)

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
func TestMaterialUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockMaterialRepo{
		materials: []entity.DimMaterial{},
		err:       nil,
	}

	uc := NewMaterialUseCase(mock)

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