package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockPurchaseRepo struct {
	purchases []entity.Purchase
	metrics   entity.PurchaseMetrics
	err       error
}

func (m *mockPurchaseRepo) FindPurchases(ctx context.Context, filter entity.PurchaseFilter) ([]entity.Purchase, error) {
	return m.purchases, m.err
}

func (m *mockPurchaseRepo) GetMetrics(ctx context.Context, filter entity.PurchaseFilter) (entity.PurchaseMetrics, error) {
	return m.metrics, m.err
}

// ==================== GET PURCHASES ====================

// Teste 1: sucesso
func TestPurchaseUseCase_GetPurchases_Success(t *testing.T) {
	// Arrange
	fake := []entity.Purchase{
		{
			IdCompra:     "1",
			Fornecedor:   "Fornecedor A",
			ValorCompra:  "1500.00",
			StatusCompra: "Aprovado",
		},
		{
			IdCompra:     "2",
			Fornecedor:   "Fornecedor B",
			ValorCompra:  "3000.00",
			StatusCompra: "Pendente",
		},
	}

	mock := &mockPurchaseRepo{
		purchases: fake,
		err:       nil,
	}

	uc := NewPurchaseUseCase(mock)

	// Act
	result, err := uc.GetPurchases(context.Background(), entity.PurchaseFilter{})

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 compras, recebeu %d", len(result))
	}

	if result[0].Fornecedor != "Fornecedor A" {
		t.Errorf("esperava fornecedor 'Fornecedor A', recebeu %s", result[0].Fornecedor)
	}

	if result[1].StatusCompra != "Pendente" {
		t.Errorf("esperava status 'Pendente', recebeu %s", result[1].StatusCompra)
	}
}

// Teste 2: erro
func TestPurchaseUseCase_GetPurchases_Error(t *testing.T) {
	// Arrange
	mock := &mockPurchaseRepo{
		purchases: nil,
		err:       errors.New("database connection failed"),
	}

	uc := NewPurchaseUseCase(mock)

	// Act
	result, err := uc.GetPurchases(context.Background(), entity.PurchaseFilter{})

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != nil {
		t.Errorf("esperava nil, recebeu %v", result)
	}
}

// Teste 3: lista vazia
func TestPurchaseUseCase_GetPurchases_Empty(t *testing.T) {
	// Arrange
	mock := &mockPurchaseRepo{
		purchases: []entity.Purchase{},
		err:       nil,
	}

	uc := NewPurchaseUseCase(mock)

	// Act
	result, err := uc.GetPurchases(context.Background(), entity.PurchaseFilter{})

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("esperava lista vazia, recebeu %d itens", len(result))
	}
}

// ==================== GET METRICS ====================

// Teste 1: sucesso
func TestPurchaseUseCase_GetMetrics_Success(t *testing.T) {
	// Arrange
	fakeMetrics := entity.PurchaseMetrics{
		TotalCompras:  10,
		ValorTotal:    "50000.00",
		ComprasAtivas: 7,
	}

	mock := &mockPurchaseRepo{
		metrics: fakeMetrics,
		err:     nil,
	}

	uc := NewPurchaseUseCase(mock)

	// Act
	result, err := uc.GetMetrics(context.Background(), entity.PurchaseFilter{})

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if result.TotalCompras != 10 {
		t.Errorf("esperava 10 compras, recebeu %d", result.TotalCompras)
	}

	if result.ComprasAtivas != 7 {
		t.Errorf("esperava 7 compras ativas, recebeu %d", result.ComprasAtivas)
	}
}

// Teste 2: erro
func TestPurchaseUseCase_GetMetrics_Error(t *testing.T) {
	// Arrange
	mock := &mockPurchaseRepo{
		metrics: entity.PurchaseMetrics{},
		err:     errors.New("metrics query failed"),
	}

	uc := NewPurchaseUseCase(mock)

	// Act
	result, err := uc.GetMetrics(context.Background(), entity.PurchaseFilter{})

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != (entity.PurchaseMetrics{}) {
		t.Errorf("esperava métricas vazias, recebeu %v", result)
	}
}