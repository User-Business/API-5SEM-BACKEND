package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockFatoExecucaoRepo struct {
	data []entity.FatoExecucaoTarefas
	err  error
}

func (m *mockFatoExecucaoRepo) FindAll() ([]entity.FatoExecucaoTarefas, error) {
	return m.data, m.err
}

func TestFatoExecucaoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.FatoExecucaoTarefas{
		{SkFato: "1", SkProjeto: "P1", HorasTrabalhadas: "8"},
		{SkFato: "2", SkProjeto: "P2", HorasTrabalhadas: "16"},
	}
	mock := &mockFatoExecucaoRepo{data: fake, err: nil}
	uc := NewFatoExecucaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 itens de execucao, recebeu: %d", len(result))
	}
	if result[0].SkProjeto != "P1" {
		t.Errorf("esperava P1, recebeu: %s", result[0].SkProjeto)
	}
}

func TestFatoExecucaoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFatoExecucaoRepo{data: nil, err: errors.New("db error")}
	uc := NewFatoExecucaoUseCase(mock)

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

func TestFatoExecucaoUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFatoExecucaoRepo{data: []entity.FatoExecucaoTarefas{}, err: nil}
	uc := NewFatoExecucaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 itens de execucao, recebeu: %d", len(result))
	}
}
