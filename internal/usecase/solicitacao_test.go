package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockSolicitacaoRepo struct {
	solicitacoes []entity.DimSolicitacao
	err           error
}

func (m *mockSolicitacaoRepo) FindAll() ([]entity.DimSolicitacao, error) {
	return m.solicitacoes, m.err
}

// Teste 1: sucesso
func TestSolicitacaoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.DimSolicitacao{
		{
			IdSolicitacao:     "1",
			CodigoSolicitacao: "SOL001",
			Descricao:         "Solicitação de compra",
			Status:            "Aberta",
		},
		{
			IdSolicitacao:     "2",
			CodigoSolicitacao: "SOL002",
			Descricao:         "Solicitação de manutenção",
			Status:            "Concluída",
		},
	}

	mock := &mockSolicitacaoRepo{
		solicitacoes: fake,
		err:          nil,
	}

	uc := NewSolicitacaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 solicitações, recebeu: %d", len(result))
	}

	if result[0].CodigoSolicitacao != "SOL001" {
		t.Errorf("esperava código SOL001, recebeu %s", result[0].CodigoSolicitacao)
	}

	if result[1].Descricao != "Solicitação de manutenção" {
		t.Errorf("esperava descrição 'Solicitação de manutenção', recebeu %s", result[1].Descricao)
	}
}

// Teste 2: erro do repository
func TestSolicitacaoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockSolicitacaoRepo{
		solicitacoes: nil,
		err:          errors.New("database connection failed"),
	}

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

// Teste 3: lista vazia
func TestSolicitacaoUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockSolicitacaoRepo{
		solicitacoes: []entity.DimSolicitacao{},
		err:          nil,
	}

	uc := NewSolicitacaoUseCase(mock)

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