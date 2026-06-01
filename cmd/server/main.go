package main

import (
	"log"
	"net/http"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/adapter/handler"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/adapter/repository"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/config"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/auth"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/database"
	appmw "github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/middleware"
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
	logImportacaoRepo := repository.NewPostgresLogImportacaoRepository(db)
	purchaseRepo := repository.NewPostgresPurchaseRepository(db)
	userRepo := repository.NewPostgresUserRepository(db)
	tokenRepo := repository.NewPostgresTokenRepository(db)

	// Auth
	jwtSvc := auth.NewService(cfg.JWTSecret, cfg.JWTExpiry)
	authUC := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtSvc)
	authMW := appmw.Auth(jwtSvc, tokenRepo)

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
	logImportacaoUC := usecase.NewLogImportacaoUseCase(logImportacaoRepo)
	purchaseUC := usecase.NewPurchaseUseCase(purchaseRepo)

	// Handlers
	handlers := router.Handlers{
		Projeto:       handler.NewProjetoHandler(projetoUC),
		Fornecedor:    handler.NewFornecedorHandler(fornecedorUC),
		Material:      handler.NewMaterialHandler(materialUC),
		Responsavel:   handler.NewResponsavelHandler(responsavelUC),
		Solicitacao:   handler.NewSolicitacaoHandler(solicitacaoUC),
		Tarefa:        handler.NewTarefaHandler(tarefaUC),
		Tempo:         handler.NewTempoHandler(tempoUC),
		FatoCompras:   handler.NewFatoComprasHandler(fatoComprasUC),
		FatoEstoque:   handler.NewFatoEstoqueHandler(fatoEstoqueUC),
		FatoExecucao:  handler.NewFatoExecucaoHandler(fatoExecucaoUC),
		LogImportacao: handler.NewLogImportacaoHandler(logImportacaoUC),
		Purchase:      handler.NewPurchaseHandler(purchaseUC),
		Auth:          handler.NewAuthHandler(authUC),
	}

	r := router.NewRouter(handlers, authMW)

	log.Printf("Server running on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
