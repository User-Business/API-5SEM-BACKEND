package usecase

import (
	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

type LogImportacaoRepository interface {
	Create(log *entity.LogImportacao) (int, error)
	Update(log *entity.LogImportacao) error
	CreateError(errLog *entity.LogImportacaoErro) error
	CreateErrorsBatch(errLogs []entity.LogImportacaoErro) error
	FindAll() ([]entity.LogImportacao, error)
	FindByID(id int) (*entity.LogImportacao, error)
	FindErrorsByLogID(logID int) ([]entity.LogImportacaoErro, error)
}

type LogImportacaoUseCase struct {
	repo LogImportacaoRepository
}

func NewLogImportacaoUseCase(repo LogImportacaoRepository) *LogImportacaoUseCase {
	return &LogImportacaoUseCase{repo: repo}
}

func (uc *LogImportacaoUseCase) Create(log *entity.LogImportacao) (int, error) {
	return uc.repo.Create(log)
}

func (uc *LogImportacaoUseCase) Update(log *entity.LogImportacao) error {
	return uc.repo.Update(log)
}

func (uc *LogImportacaoUseCase) CreateError(errLog *entity.LogImportacaoErro) error {
	return uc.repo.CreateError(errLog)
}

func (uc *LogImportacaoUseCase) CreateErrorsBatch(errLogs []entity.LogImportacaoErro) error {
	return uc.repo.CreateErrorsBatch(errLogs)
}

func (uc *LogImportacaoUseCase) GetAll() ([]entity.LogImportacao, error) {
	return uc.repo.FindAll()
}

func (uc *LogImportacaoUseCase) GetByID(id int) (*entity.LogImportacao, error) {
	return uc.repo.FindByID(id)
}

func (uc *LogImportacaoUseCase) GetErrorsByLogID(logID int) ([]entity.LogImportacaoErro, error) {
	return uc.repo.FindErrorsByLogID(logID)
}
