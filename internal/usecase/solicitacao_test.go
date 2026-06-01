package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockSolicitacaoRepo struct {
	data []entity.DimSolicitacao
	err  error
}

func (m *mockSolicitacaoRepo) FindAll() ([]entity.DimSolicitacao, error) {
	return m.data, m.err
}

func TestSolicitacaoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimSolicitacao{
		{SkSolicitacao: "1", IdSolicitacao: "S1", NumeroSolicitacao: "SC-100"},
		{SkSolicitacao: "2", IdSolicitacao: "S2", NumeroSolicitacao: "SC-200"},
	}
	mock := &mockSolicitacaoRepo{data: fake, err: nil}
	uc := NewSolicitacaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 solicitacoes, recebeu: %d", len(result))
	}
	if result[0].NumeroSolicitacao != "SC-100" {
		t.Errorf("esperava SC-100, recebeu: %s", result[0].NumeroSolicitacao)
	}
}

func TestSolicitacaoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockSolicitacaoRepo{data: nil, err: errors.New("db error")}
	uc := NewSolicitacaoUseCase(mock)

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

func TestSolicitacaoUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockSolicitacaoRepo{data: []entity.DimSolicitacao{}, err: nil}
	uc := NewSolicitacaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("esperava 0 solicitacoes, recebeu: %d", len(result))
	}
}
