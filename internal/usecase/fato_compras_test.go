package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockFatoComprasRepo struct {
	compras []entity.FatoCompras
	err     error
}

func (m *mockFatoComprasRepo) FindAll() ([]entity.FatoCompras, error) {
	return m.compras, m.err
}

// Teste 1: sucesso
func TestFatoComprasUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.FatoCompras{
    	{
    		SkFato:              "1",
    		SkProjeto:           "101",
    		SkFornecedor:        "201",
    		SkSolicitacao:       "301",
    		SkTempo:             "401",
    		ValorTotalPedido:    "1500.00",
    		ValorAlocadoProjeto: "1200.00",
    	},
    	{
    		SkFato:              "2",
    		SkProjeto:           "102",
    		SkFornecedor:        "202",
    		SkSolicitacao:       "302",
    		SkTempo:             "402",
    		ValorTotalPedido:    "3200.00",
    		ValorAlocadoProjeto: "2800.00",
    	},
    }

	mock := &mockFatoComprasRepo{
		compras: fake,
		err:     nil,
	}

	uc := NewFatoComprasUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 compras, recebeu: %d", len(result))
	}

	if result[0].ValorTotalPedido != "1500.00" {
    	t.Errorf("esperava valor 1500.00, recebeu %s", result[0].ValorTotalPedido)
    }

	if result[1].ValorAlocadoProjeto != "2800.00" {
    	t.Errorf("esperava valor alocado 2800.00, recebeu %s", result[1].ValorAlocadoProjeto)
    }
}

// Teste 2: erro do repository
func TestFatoComprasUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockFatoComprasRepo{
		compras: nil,
		err:     errors.New("database connection failed"),
	}

	uc := NewFatoComprasUseCase(mock)

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
func TestFatoComprasUseCase_GetAll_Empty(t *testing.T) {
	// Arrange
	mock := &mockFatoComprasRepo{
		compras: []entity.FatoCompras{},
		err:     nil,
	}

	uc := NewFatoComprasUseCase(mock)

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