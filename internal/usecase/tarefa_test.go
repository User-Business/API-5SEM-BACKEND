package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockTarefaRepo struct {
	tarefas []entity.DimTarefa
	err     error
}

func (m *mockTarefaRepo) FindAll() ([]entity.DimTarefa, error) {
	return m.tarefas, m.err
}

// Teste 1: sucesso
func TestTarefaUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimTarefa{
		{
			IdTarefa:      "1",
			CodigoTarefa:  "TAR001",
			NomeTarefa:    "Planejamento",
			StatusTarefa:  "Concluída",
		},
		{
			IdTarefa:      "2",
			CodigoTarefa:  "TAR002",
			NomeTarefa:    "Execução",
			StatusTarefa:  "Em andamento",
		},
	}

	mock := &mockTarefaRepo{
		tarefas: fake,
		err:     nil,
	}

	uc := NewTarefaUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 tarefas, recebeu: %d", len(result))
	}

	if result[0].CodigoTarefa != "TAR001" {
		t.Errorf("esperava código TAR001, recebeu %s", result[0].CodigoTarefa)
	}

	if result[1].NomeTarefa != "Execução" {
		t.Errorf("esperava tarefa 'Execução', recebeu %s", result[1].NomeTarefa)
	}
}

// Teste 2: erro do repository
func TestTarefaUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockTarefaRepo{
		tarefas: nil,
		err:     errors.New("database connection failed"),
	}

	uc := NewTarefaUseCase(mock)

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
func TestTarefaUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockTarefaRepo{
		tarefas: []entity.DimTarefa{},
		err:     nil,
	}

	uc := NewTarefaUseCase(mock)

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