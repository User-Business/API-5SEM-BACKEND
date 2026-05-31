package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type mockLogImportacaoRepo struct {
	createdID int
	err       error
	logs      []entity.LogImportacao
	singleLog *entity.LogImportacao
	errors    []entity.LogImportacaoErro
}

func (m *mockLogImportacaoRepo) Create(log *entity.LogImportacao) (int, error) {
	return m.createdID, m.err
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
	return m.singleLog, m.err
}

func (m *mockLogImportacaoRepo) FindErrorsByLogID(logID int) ([]entity.LogImportacaoErro, error) {
	return m.errors, m.err
}

// Create
func TestLogImportacaoUseCase_Create_Success(t *testing.T) {
	mock := &mockLogImportacaoRepo{createdID: 42, err: nil}
	uc := NewLogImportacaoUseCase(mock)

	id, err := uc.Create(&entity.LogImportacao{NomeArquivo: "test.xlsx"})

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if id != 42 {
		t.Fatalf("esperava ID 42, recebeu: %d", id)
	}
}

func TestLogImportacaoUseCase_Create_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{createdID: 0, err: errors.New("insert error")}
	uc := NewLogImportacaoUseCase(mock)

	_, err := uc.Create(&entity.LogImportacao{NomeArquivo: "test.xlsx"})

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// Update
func TestLogImportacaoUseCase_Update_Success(t *testing.T) {
	mock := &mockLogImportacaoRepo{err: nil}
	uc := NewLogImportacaoUseCase(mock)

	err := uc.Update(&entity.LogImportacao{ID: 1})

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
}

func TestLogImportacaoUseCase_Update_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{err: errors.New("update error")}
	uc := NewLogImportacaoUseCase(mock)

	err := uc.Update(&entity.LogImportacao{ID: 1})

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// CreateError
func TestLogImportacaoUseCase_CreateError_Success(t *testing.T) {
	mock := &mockLogImportacaoRepo{err: nil}
	uc := NewLogImportacaoUseCase(mock)

	err := uc.CreateError(&entity.LogImportacaoErro{LogImportacaoID: 1})

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
}

func TestLogImportacaoUseCase_CreateError_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{err: errors.New("create error error")}
	uc := NewLogImportacaoUseCase(mock)

	err := uc.CreateError(&entity.LogImportacaoErro{LogImportacaoID: 1})

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// CreateErrorsBatch
func TestLogImportacaoUseCase_CreateErrorsBatch_Success(t *testing.T) {
	mock := &mockLogImportacaoRepo{err: nil}
	uc := NewLogImportacaoUseCase(mock)

	err := uc.CreateErrorsBatch([]entity.LogImportacaoErro{{LogImportacaoID: 1}})

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
}

func TestLogImportacaoUseCase_CreateErrorsBatch_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{err: errors.New("batch error")}
	uc := NewLogImportacaoUseCase(mock)

	err := uc.CreateErrorsBatch([]entity.LogImportacaoErro{{LogImportacaoID: 1}})

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

// GetAll
func TestLogImportacaoUseCase_GetAll_Success(t *testing.T) {
	fake := []entity.LogImportacao{
		{ID: 1, NomeArquivo: "a.xlsx", DataInicio: time.Now()},
	}
	mock := &mockLogImportacaoRepo{logs: fake, err: nil}
	uc := NewLogImportacaoUseCase(mock)

	res, err := uc.GetAll()

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("esperava 1 log, recebeu: %d", len(res))
	}
}

func TestLogImportacaoUseCase_GetAll_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{logs: nil, err: errors.New("db error")}
	uc := NewLogImportacaoUseCase(mock)

	res, err := uc.GetAll()

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if res != nil {
		t.Errorf("esperava nil, recebeu: %v", res)
	}
}

// GetByID
func TestLogImportacaoUseCase_GetByID_Success(t *testing.T) {
	fake := &entity.LogImportacao{ID: 1, NomeArquivo: "a.xlsx"}
	mock := &mockLogImportacaoRepo{singleLog: fake, err: nil}
	uc := NewLogImportacaoUseCase(mock)

	res, err := uc.GetByID(1)

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if res == nil || res.ID != 1 {
		t.Fatalf("esperava log com ID 1, recebeu: %v", res)
	}
}

func TestLogImportacaoUseCase_GetByID_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{singleLog: nil, err: errors.New("not found")}
	uc := NewLogImportacaoUseCase(mock)

	res, err := uc.GetByID(1)

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if res != nil {
		t.Errorf("esperava nil, recebeu: %v", res)
	}
}

// GetErrorsByLogID
func TestLogImportacaoUseCase_GetErrorsByLogID_Success(t *testing.T) {
	fake := []entity.LogImportacaoErro{
		{ID: 10, LogImportacaoID: 1, MotivoErro: "Invalido"},
	}
	mock := &mockLogImportacaoRepo{errors: fake, err: nil}
	uc := NewLogImportacaoUseCase(mock)

	res, err := uc.GetErrorsByLogID(1)

	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("esperava 1 erro, recebeu: %d", len(res))
	}
}

func TestLogImportacaoUseCase_GetErrorsByLogID_Error(t *testing.T) {
	mock := &mockLogImportacaoRepo{errors: nil, err: errors.New("db error")}
	uc := NewLogImportacaoUseCase(mock)

	res, err := uc.GetErrorsByLogID(1)

	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
	if res != nil {
		t.Errorf("esperava nil, recebeu: %v", res)
	}
}
