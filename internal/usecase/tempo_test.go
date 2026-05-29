package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockTempoRepo struct {
	tempos      []entity.DimTempo
	tempoGasto  interface{}
	err         error
}

func (m *mockTempoRepo) FindAll() ([]entity.DimTempo, error) {
	return m.tempos, m.err
}

func (m *mockTempoRepo) GetTempoGasto() (interface{}, error) {
	return m.tempoGasto, m.err
}

// ==================== GET ALL ====================

// Teste 1: sucesso
func TestTempoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimTempo{
		{
			IdTempo: "1",
			Ano:     2025,
			Mes:     5,
			Dia:     29,
		},
		{
			IdTempo: "2",
			Ano:     2025,
			Mes:     6,
			Dia:     1,
		},
	}

	mock := &mockTempoRepo{
		tempos: fake,
		err:    nil,
	}

	uc := NewTempoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 registros, recebeu: %d", len(result))
	}

	if result[0].Ano != 2025 {
		t.Errorf("esperava ano 2025, recebeu %d", result[0].Ano)
	}

	if result[1].Mes != 6 {
		t.Errorf("esperava mês 6, recebeu %d", result[1].Mes)
	}
}

// Teste 2: erro do repository
func TestTempoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{
		tempos: nil,
		err:    errors.New("database connection failed"),
	}

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

// Teste 3: lista vazia
func TestTempoUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{
		tempos: []entity.DimTempo{},
		err:    nil,
	}

	uc := NewTempoUseCase(mock)

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

// ==================== GET TEMPO GASTO ====================

// Teste 1: sucesso
func TestTempoUseCase_GetTempoGasto_Success(t *testing.T) {
	// Arrange
	fakeTempo := map[string]interface{}{
		"tempo_total_horas": 120,
		"projeto":           "Projeto Alpha",
	}

	mock := &mockTempoRepo{
		tempoGasto: fakeTempo,
		err:        nil,
	}

	uc := NewTempoUseCase(mock)

	// Act
	result, err := uc.GetTempoGasto()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if result == nil {
		t.Fatal("esperava resultado, recebeu nil")
	}
}

// Teste 2: erro
func TestTempoUseCase_GetTempoGasto_Error(t *testing.T) {
	// Arrange
	mock := &mockTempoRepo{
		tempoGasto: nil,
		err:        errors.New("query failed"),
	}

	uc := NewTempoUseCase(mock)

	// Act
	result, err := uc.GetTempoGasto()

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != nil {
		t.Errorf("esperava nil, recebeu %v", result)
	}
}