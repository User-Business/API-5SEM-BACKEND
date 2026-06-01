package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

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

// GetPurchases
func TestPurchaseUseCase_GetPurchases_Success(t *testing.T) {
	fake := []entity.Purchase{
		{ID: 1, Type: "SC", Number: "123", Status: "Aprovado"},
		{ID: 2, Type: "PC", Number: "456", Status: "Pendente"},
	}
	mock := &mockPurchaseRepo{purchases: fake, err: nil}
	uc := NewPurchaseUseCase(mock)

	result, err := uc.GetPurchases(context.Background(), entity.PurchaseFilter{})

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("esperava 2 compras, recebeu: %d", len(result))
	}
	if result[0].Number != "123" {
		t.Errorf("esperava numero 123, recebeu: %s", result[0].Number)
	}
}

func TestPurchaseUseCase_GetPurchases_Error(t *testing.T) {
	mock := &mockPurchaseRepo{purchases: nil, err: errors.New("db error")}
	uc := NewPurchaseUseCase(mock)

	result, err := uc.GetPurchases(context.Background(), entity.PurchaseFilter{})

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if result != nil {
		t.Errorf("esperava nil, recebeu: %v", result)
	}
}

// GetMetrics
func TestPurchaseUseCase_GetMetrics_Success(t *testing.T) {
	fakeMetrics := entity.PurchaseMetrics{
		TotalPurchases: 10,
		TotalSC:        6,
		TotalPC:        4,
	}
	mock := &mockPurchaseRepo{metrics: fakeMetrics, err: nil}
	uc := NewPurchaseUseCase(mock)

	result, err := uc.GetMetrics(context.Background(), entity.PurchaseFilter{})

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if result.TotalPurchases != 10 {
		t.Errorf("esperava 10 compras totais, recebeu: %d", result.TotalPurchases)
	}
}

func TestPurchaseUseCase_GetMetrics_Error(t *testing.T) {
	mock := &mockPurchaseRepo{err: errors.New("db error")}
	uc := NewPurchaseUseCase(mock)

	_, err := uc.GetMetrics(context.Background(), entity.PurchaseFilter{})

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}
