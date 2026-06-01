package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockTarefaRepo struct {
	data []entity.DimTarefa
	err  error
}

func (m *mockTarefaRepo) FindAll() ([]entity.DimTarefa, error) {
	return m.data, m.err
}

func TestTarefaUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimTarefa{
		{SkTarefa: "1", IdTarefa: "T1", Titulo: "Tarefa Alpha"},
		{SkTarefa: "2", IdTarefa: "T2", Titulo: "Tarefa Beta"},
	}
	mock := &mockTarefaRepo{data: fake, err: nil}
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
	if result[0].Titulo != "Tarefa Alpha" {
		t.Errorf("esperava Tarefa Alpha, recebeu: %s", result[0].Titulo)
	}
}

func TestTarefaUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockTarefaRepo{data: nil, err: errors.New("db error")}
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

func TestTarefaUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockTarefaRepo{data: []entity.DimTarefa{}, err: nil}
	uc := NewTarefaUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 tarefas, recebeu: %d", len(result))
	}
}
