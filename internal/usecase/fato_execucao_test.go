package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockFatoExecucaoRepo struct {
	execucoes []entity.FatoExecucaoTarefas
	err       error
}

func (m *mockFatoExecucaoRepo) FindAll() ([]entity.FatoExecucaoTarefas, error) {
	return m.execucoes, m.err
}

// Teste 1: sucesso
func TestFatoExecucaoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.FatoExecucaoTarefas{
		{
			IdTarefa:      "1",
			NomeTarefa:    "Planejamento",
			StatusTarefa:  "Concluída",
			Responsavel:   "João",
		},
		{
			IdTarefa:      "2",
			NomeTarefa:    "Implementação",
			StatusTarefa:  "Em andamento",
			Responsavel:   "Maria",
		},
	}

	mock := &mockFatoExecucaoRepo{
		execucoes: fake,
		err:       nil,
	}

	uc := NewFatoExecucaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 registros, recebeu: %d", len(result))
	}

	if result[0].NomeTarefa != "Planejamento" {
		t.Errorf("esperava tarefa 'Planejamento', recebeu %s", result[0].NomeTarefa)
	}

	if result[1].Responsavel != "Maria" {
		t.Errorf("esperava responsável 'Maria', recebeu %s", result[1].Responsavel)
	}
}

// Teste 2: erro do repository
func TestFatoExecucaoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFatoExecucaoRepo{
		execucoes: nil,
		err:       errors.New("database connection failed"),
	}

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

// Teste 3: lista vazia
func TestFatoExecucaoUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFatoExecucaoRepo{
		execucoes: []entity.FatoExecucaoTarefas{},
		err:       nil,
	}

	uc := NewFatoExecucaoUseCase(mock)

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