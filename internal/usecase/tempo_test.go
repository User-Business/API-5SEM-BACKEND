package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockTempoRepo struct {
	data       []entity.DimTempo
	tempoGasto interface{}
	err        error
}

func (m *mockTempoRepo) FindAll() ([]entity.DimTempo, error) {
	return m.data, m.err
}

func (m *mockTempoRepo) GetTempoGasto() (interface{}, error) {
	return m.tempoGasto, m.err
}

// GetAll
func TestTempoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimTempo{
		{SkTempo: "1", DataCompleta: "2026-05-31"},
		{SkTempo: "2", DataCompleta: "2026-06-01"},
	}
	mock := &mockTempoRepo{data: fake, err: nil}
	uc := NewTempoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 registros de tempo, recebeu: %d", len(result))
	}
	if result[0].DataCompleta != "2026-05-31" {
		t.Errorf("esperava 2026-05-31, recebeu: %s", result[0].DataCompleta)
	}
}

func TestTempoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{data: nil, err: errors.New("db error")}
	uc := NewTempoUseCase(mock)

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

func TestTempoUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{data: []entity.DimTempo{}, err: nil}
	uc := NewTempoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 registros de tempo, recebeu: %d", len(result))
	}
}

// GetTempoGasto
func TestTempoUseCase_GetTempoGasto_Success(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{tempoGasto: "120h", err: nil}
	uc := NewTempoUseCase(mock)

	// Act
	result, err := uc.GetTempoGasto()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if result != "120h" {
		t.Errorf("esperava 120h, recebeu: %v", result)
	}
}

func TestTempoUseCase_GetTempoGasto_Error(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{tempoGasto: nil, err: errors.New("db error")}
	uc := NewTempoUseCase(mock)

	// Act
	_, err := uc.GetTempoGasto()

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}
