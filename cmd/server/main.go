package main

import (
	"log"
	"net/http"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/adapter/handler"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/adapter/repository"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/config"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/database"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/router"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Repositories
	projetoRepo := repository.NewPostgresProjetoRepository(db)
	fornecedorRepo := repository.NewPostgresFornecedorRepository(db)
	materialRepo := repository.NewPostgresMaterialRepository(db)
	responsavelRepo := repository.NewPostgresResponsavelRepository(db)
	solicitacaoRepo := repository.NewPostgresSolicitacaoRepository(db)
	tarefaRepo := repository.NewPostgresTarefaRepository(db)
	tempoRepo := repository.NewPostgresTempoRepository(db)
	fatoComprasRepo := repository.NewPostgresFatoComprasRepository(db)
	fatoEstoqueRepo := repository.NewPostgresFatoEstoqueRepository(db)
	fatoExecucaoRepo := repository.NewPostgresFatoExecucaoRepository(db)
	purchaseRepo := repository.NewPostgresPurchaseRepository(db)

	// Use Cases
	projetoUC := usecase.NewProjetoUseCase(projetoRepo)
	fornecedorUC := usecase.NewFornecedorUseCase(fornecedorRepo)
	materialUC := usecase.NewMaterialUseCase(materialRepo)
	responsavelUC := usecase.NewResponsavelUseCase(responsavelRepo)
	solicitacaoUC := usecase.NewSolicitacaoUseCase(solicitacaoRepo)
	tarefaUC := usecase.NewTarefaUseCase(tarefaRepo)
	tempoUC := usecase.NewTempoUseCase(tempoRepo)
	fatoComprasUC := usecase.NewFatoComprasUseCase(fatoComprasRepo)
	fatoEstoqueUC := usecase.NewFatoEstoqueUseCase(fatoEstoqueRepo)
	fatoExecucaoUC := usecase.NewFatoExecucaoUseCase(fatoExecucaoRepo)
	purchaseUC := usecase.NewPurchaseUseCase(purchaseRepo)

	// Handlers
	handlers := router.Handlers{
		Projeto:      handler.NewProjetoHandler(projetoUC),
		Fornecedor:   handler.NewFornecedorHandler(fornecedorUC),
		Material:     handler.NewMaterialHandler(materialUC),
		Responsavel:  handler.NewResponsavelHandler(responsavelUC),
		Solicitacao:  handler.NewSolicitacaoHandler(solicitacaoUC),
		Tarefa:       handler.NewTarefaHandler(tarefaUC),
		Tempo:        handler.NewTempoHandler(tempoUC),
		FatoCompras:  handler.NewFatoComprasHandler(fatoComprasUC),
		FatoEstoque:  handler.NewFatoEstoqueHandler(fatoEstoqueUC),
		FatoExecucao: handler.NewFatoExecucaoHandler(fatoExecucaoUC),
		Purchase:     handler.NewPurchaseHandler(purchaseUC),
	}

	r := router.NewRouter(handlers)

	log.Printf("Server running on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
