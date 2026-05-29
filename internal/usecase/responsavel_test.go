package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockResponsavelRepo struct {
	responsaveis []entity.DimResponsavel
	err          error
}

func (m *mockResponsavelRepo) FindAll() ([]entity.DimResponsavel, error) {
	return m.responsaveis, m.err
}

// Teste 1: sucesso
func TestResponsavelUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimResponsavel{
		{
			IdResponsavel:   "1",
			NomeResponsavel: "João Silva",
			Email:           "joao@empresa.com",
			Cargo:           "Gerente",
			Status:          "Ativo",
		},
		{
			IdResponsavel:   "2",
			NomeResponsavel: "Maria Souza",
			Email:           "maria@empresa.com",
			Cargo:           "Coordenadora",
			Status:          "Ativo",
		},
	}

	mock := &mockResponsavelRepo{
		responsaveis: fake,
		err:          nil,
	}

	uc := NewResponsavelUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 responsáveis, recebeu: %d", len(result))
	}

	if result[0].NomeResponsavel != "João Silva" {
		t.Errorf("esperava responsável 'João Silva', recebeu %s", result[0].NomeResponsavel)
	}

	if result[1].Cargo != "Coordenadora" {
		t.Errorf("esperava cargo 'Coordenadora', recebeu %s", result[1].Cargo)
	}
}

// Teste 2: erro do repository
func TestResponsavelUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockResponsavelRepo{
		responsaveis: nil,
		err:          errors.New("database connection failed"),
	}

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

// Teste 3: lista vazia
func TestResponsavelUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockResponsavelRepo{
		responsaveis: []entity.DimResponsavel{},
		err:          nil,
	}

	uc := NewResponsavelUseCase(mock)

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