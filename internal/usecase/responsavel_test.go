package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockResponsavelRepo struct {
	data []entity.DimResponsavel
	err  error
}

func (m *mockResponsavelRepo) FindAll() ([]entity.DimResponsavel, error) {
	return m.data, m.err
}

func TestResponsavelUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimResponsavel{
		{SkResponsavel: "1", NomeResponsavel: "Alice", Tipo: "Gerente"},
		{SkResponsavel: "2", NomeResponsavel: "Bob", Tipo: "Desenvolvedor"},
	}
	mock := &mockResponsavelRepo{data: fake, err: nil}
	uc := NewResponsavelUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 responsaveis, recebeu: %d", len(result))
	}
	if result[0].NomeResponsavel != "Alice" {
		t.Errorf("esperava Alice, recebeu: %s", result[0].NomeResponsavel)
	}
}

func TestResponsavelUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockResponsavelRepo{data: nil, err: errors.New("db error")}
	uc := NewResponsavelUseCase(mock)

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

func TestResponsavelUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockResponsavelRepo{data: []entity.DimResponsavel{}, err: nil}
	uc := NewResponsavelUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 responsaveis, recebeu: %d", len(result))
	}
}
