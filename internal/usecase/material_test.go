package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockMaterialRepo struct {
	data []entity.DimMaterial
	err  error
}

func (m *mockMaterialRepo) FindAll() ([]entity.DimMaterial, error) {
	return m.data, m.err
}

func TestMaterialUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimMaterial{
		{SkMaterial: "1", IdMaterial: "M1", Descricao: "Parafuso"},
		{SkMaterial: "2", IdMaterial: "M2", Descricao: "Prego"},
	}
	mock := &mockMaterialRepo{data: fake, err: nil}
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
	if result[0].Descricao != "Parafuso" {
		t.Errorf("esperava Parafuso, recebeu: %s", result[0].Descricao)
	}
}

func TestMaterialUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockMaterialRepo{data: nil, err: errors.New("db error")}
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

func TestMaterialUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockMaterialRepo{data: []entity.DimMaterial{}, err: nil}
	uc := NewMaterialUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 materiais, recebeu: %d", len(result))
	}
}
