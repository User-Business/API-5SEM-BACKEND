package usecase

import (
	"errors"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// Mock do Repository
type mockLogImportacaoRepo struct {
	logs    []entity.LogImportacao
	log     *entity.LogImportacao
	errLogs []entity.LogImportacaoErro

	id  int
	err error
}

func (m *mockLogImportacaoRepo) Create(log *entity.LogImportacao) (int, error) {
	return m.id, m.err
}

func (m *mockLogImportacaoRepo) Update(log *entity.LogImportacao) error {
	return m.err
}

func (m *mockLogImportacaoRepo) CreateError(errLog *entity.LogImportacaoErro) error {
	return m.err
}

func (m *mockLogImportacaoRepo) CreateErrorsBatch(errLogs []entity.LogImportacaoErro) error {
	return m.err
}

func (m *mockLogImportacaoRepo) FindAll() ([]entity.LogImportacao, error) {
	return m.logs, m.err
}

func (m *mockLogImportacaoRepo) FindByID(id int) (*entity.LogImportacao, error) {
	return m.log, m.err
}

func (m *mockLogImportacaoRepo) FindErrorsByLogID(logID int) ([]entity.LogImportacaoErro, error) {
	return m.errLogs, m.err
}

// ==================== CREATE ====================

func TestLogImportacaoUseCase_Create_Success(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		id:  1,
		err: nil,
	}

	uc := NewLogImportacaoUseCase(mock)

	log := &entity.LogImportacao{}

	// Act
	result, err := uc.Create(log)

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if result != 1 {
		t.Errorf("esperava ID 1, recebeu %d", result)
	}
}

func TestLogImportacaoUseCase_Create_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		id:  0,
		err: errors.New("database error"),
	}

	uc := NewLogImportacaoUseCase(mock)

	log := &entity.LogImportacao{}

	// Act
	result, err := uc.Create(log)

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != 0 {
		t.Errorf("esperava ID 0, recebeu %d", result)
	}
}

// ==================== UPDATE ====================

func TestLogImportacaoUseCase_Update_Success(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{err: nil}

	uc := NewLogImportacaoUseCase(mock)

	log := &entity.LogImportacao{}

	// Act
	err := uc.Update(log)

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
}

func TestLogImportacaoUseCase_Update_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		err: errors.New("update failed"),
	}

	uc := NewLogImportacaoUseCase(mock)

	log := &entity.LogImportacao{}

	// Act
	err := uc.Update(log)

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// ==================== CREATE ERROR ====================

func TestLogImportacaoUseCase_CreateError_Success(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{err: nil}

	uc := NewLogImportacaoUseCase(mock)

	errLog := &entity.LogImportacaoErro{}

	// Act
	err := uc.CreateError(errLog)

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
}

func TestLogImportacaoUseCase_CreateError_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		err: errors.New("create error failed"),
	}

	uc := NewLogImportacaoUseCase(mock)

	errLog := &entity.LogImportacaoErro{}

	// Act
	err := uc.CreateError(errLog)

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// ==================== CREATE ERRORS BATCH ====================

func TestLogImportacaoUseCase_CreateErrorsBatch_Success(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{err: nil}

	uc := NewLogImportacaoUseCase(mock)

	errLogs := []entity.LogImportacaoErro{}

	// Act
	err := uc.CreateErrorsBatch(errLogs)

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
}

func TestLogImportacaoUseCase_CreateErrorsBatch_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		err: errors.New("batch insert failed"),
	}

	uc := NewLogImportacaoUseCase(mock)

	errLogs := []entity.LogImportacaoErro{}

	// Act
	err := uc.CreateErrorsBatch(errLogs)

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// ==================== GET ALL ====================

func TestLogImportacaoUseCase_GetAll_Success(t *testing.T) {
	// Arrange
	fake := []entity.LogImportacao{
		{Id: 1},
		{Id: 2},
	}

	mock := &mockLogImportacaoRepo{
		logs: fake,
		err:  nil,
	}

	uc := NewLogImportacaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 logs, recebeu %d", len(result))
	}
}

func TestLogImportacaoUseCase_GetAll_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		logs: nil,
		err:  errors.New("find all failed"),
	}

	uc := NewLogImportacaoUseCase(mock)

	// Act
	result, err := uc.GetAll()

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != nil {
		t.Errorf("esperava nil, recebeu %v", result)
	}
}

// ==================== GET BY ID ====================

func TestLogImportacaoUseCase_GetByID_Success(t *testing.T) {
	// Arrange
	fake := &entity.LogImportacao{
		Id: 1,
	}

	mock := &mockLogImportacaoRepo{
		log: fake,
		err: nil,
	}

	uc := NewLogImportacaoUseCase(mock)

	// Act
	result, err := uc.GetByID(1)

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if result.Id != 1 {
		t.Errorf("esperava ID 1, recebeu %d", result.Id)
	}
}

func TestLogImportacaoUseCase_GetByID_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		log: nil,
		err: errors.New("not found"),
	}

	uc := NewLogImportacaoUseCase(mock)

	// Act
	result, err := uc.GetByID(1)

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != nil {
		t.Errorf("esperava nil, recebeu %v", result)
	}
}

// ==================== GET ERRORS BY LOG ID ====================

func TestLogImportacaoUseCase_GetErrorsByLogID_Success(t *testing.T) {
	// Arrange
	fake := []entity.LogImportacaoErro{
		{Id: 1},
		{Id: 2},
	}

	mock := &mockLogImportacaoRepo{
		errLogs: fake,
		err:     nil,
	}

	uc := NewLogImportacaoUseCase(mock)

	// Act
	result, err := uc.GetErrorsByLogID(1)

	// Assert
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("esperava 2 erros, recebeu %d", len(result))
	}
}

func TestLogImportacaoUseCase_GetErrorsByLogID_Error(t *testing.T) {
	// Arrange
	mock := &mockLogImportacaoRepo{
		errLogs: nil,
		err:     errors.New("find errors failed"),
	}

	uc := NewLogImportacaoUseCase(mock)

	// Act
	result, err := uc.GetErrorsByLogID(1)

	// Assert
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}

	if result != nil {
		t.Errorf("esperava nil, recebeu %v", result)
	}
}