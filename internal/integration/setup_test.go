package integration

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/adapter/handler"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/auth"
	appmw "github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/middleware"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/router"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/usecase"
)

// Credenciais do usuário admin semeado por padrão em cada teste.
const (
	testJWTSecret   = "test-secret-nao-usar-em-producao"
	defaultEmail    = "admin@teste.com"
	defaultPassword = "senha-correta-123"
)

// errInjetado é usado pelos testes para forçar uma falha na camada de dados.
var errInjetado = errors.New("erro simulado de banco")

// testEnv é o ambiente de um teste: o servidor HTTP real + acesso aos
// repositórios falsos, para que cada caso possa semear dados ou forçar erros.
type testEnv struct {
	server  *httptest.Server
	jwt     *auth.Service
	users   *fakeUserRepo
	tokens  *fakeTokenRepo
	projeto *fakeProjetoRepo
	compras *fakeFatoComprasRepo
}

// newTestEnv monta a aplicação inteira com repositórios em memória e devolve
// um servidor httptest pronto para receber requisições. O servidor é fechado
// automaticamente ao fim do teste.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	jwtSvc := auth.NewService(testJWTSecret, time.Hour)

	users := &fakeUserRepo{users: map[string]*entity.Usuario{}}
	tokens := &fakeTokenRepo{revoked: map[string]bool{}}

	// Usuário admin padrão (para os testes de login).
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("gerar hash de senha: %v", err)
	}
	users.users[defaultEmail] = &entity.Usuario{
		ID:        1,
		Email:     defaultEmail,
		SenhaHash: string(hash),
		Nome:      "Admin Teste",
		Role:      entity.RoleAdmin,
	}

	projeto := &fakeProjetoRepo{}
	compras := &fakeFatoComprasRepo{}

	// Auth + middleware reais.
	authUC := usecase.NewAuthUseCase(users, tokens, jwtSvc)
	authMW := appmw.Auth(jwtSvc, tokens)

	// Handlers reais, cada um envolvendo um use case real sobre um repo falso.
	handlers := router.Handlers{
		Projeto:       handler.NewProjetoHandler(usecase.NewProjetoUseCase(projeto)),
		Fornecedor:    handler.NewFornecedorHandler(usecase.NewFornecedorUseCase(&fakeFornecedorRepo{})),
		Material:      handler.NewMaterialHandler(usecase.NewMaterialUseCase(&fakeMaterialRepo{})),
		Responsavel:   handler.NewResponsavelHandler(usecase.NewResponsavelUseCase(&fakeResponsavelRepo{})),
		Solicitacao:   handler.NewSolicitacaoHandler(usecase.NewSolicitacaoUseCase(&fakeSolicitacaoRepo{})),
		Tarefa:        handler.NewTarefaHandler(usecase.NewTarefaUseCase(&fakeTarefaRepo{})),
		Tempo:         handler.NewTempoHandler(usecase.NewTempoUseCase(&fakeTempoRepo{})),
		FatoCompras:   handler.NewFatoComprasHandler(usecase.NewFatoComprasUseCase(compras)),
		FatoEstoque:   handler.NewFatoEstoqueHandler(usecase.NewFatoEstoqueUseCase(&fakeFatoEstoqueRepo{})),
		FatoExecucao:  handler.NewFatoExecucaoHandler(usecase.NewFatoExecucaoUseCase(&fakeFatoExecucaoRepo{})),
		LogImportacao: handler.NewLogImportacaoHandler(usecase.NewLogImportacaoUseCase(&fakeLogImportacaoRepo{})),
		Purchase:      handler.NewPurchaseHandler(usecase.NewPurchaseUseCase(&fakePurchaseRepo{})),
		Auth:          handler.NewAuthHandler(authUC),
	}

	srv := httptest.NewServer(router.NewRouter(handlers, authMW))
	t.Cleanup(srv.Close)

	return &testEnv{
		server:  srv,
		jwt:     jwtSvc,
		users:   users,
		tokens:  tokens,
		projeto: projeto,
		compras: compras,
	}
}

// tokenFor gera um JWT válido para a role informada, sem precisar fazer login.
// Útil para testar autorização (RBAC) de forma direta.
func (e *testEnv) tokenFor(role string) string {
	token, _, _, err := e.jwt.Generate(1, role)
	if err != nil {
		panic("falha ao gerar token de teste: " + err.Error())
	}
	return token
}

// request dispara uma requisição contra o servidor de teste. Se token != "",
// envia o header Authorization: Bearer. O corpo da resposta deve ser fechado
// pelo chamador.
func (e *testEnv) request(t *testing.T, method, path, token string, body io.Reader) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, e.server.URL+path, body)
	if err != nil {
		t.Fatalf("montar request %s %s: %v", method, path, err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("executar request %s %s: %v", method, path, err)
	}
	return resp
}

// ---------------------------------------------------------------------------
// Repositórios falsos (em memória). Cada um implementa a interface que o
// respectivo use case espera, devolvendo dados semeados e/ou um erro forçado.
// ---------------------------------------------------------------------------

type fakeUserRepo struct {
	users map[string]*entity.Usuario // chaveado por e-mail
}

func (f *fakeUserRepo) FindByEmail(_ context.Context, email string) (*entity.Usuario, error) {
	return f.users[email], nil // nil quando não existe -> credenciais inválidas
}

type fakeTokenRepo struct {
	revoked map[string]bool
}

func (f *fakeTokenRepo) Revoke(_ context.Context, jti string, _ time.Time) error {
	f.revoked[jti] = true
	return nil
}

func (f *fakeTokenRepo) IsRevoked(_ context.Context, jti string) (bool, error) {
	return f.revoked[jti], nil
}

type fakeProjetoRepo struct {
	projetos      []entity.DimProjeto
	investimentos []entity.ProgramaInvestimento
	materiais     []entity.ProjetoMaterial
	err           error
}

func (f *fakeProjetoRepo) FindAll() ([]entity.DimProjeto, error) { return f.projetos, f.err }
func (f *fakeProjetoRepo) FindInvestimentoByPrograma() ([]entity.ProgramaInvestimento, error) {
	return f.investimentos, f.err
}
func (f *fakeProjetoRepo) FindMateriaisByProjeto() ([]entity.ProjetoMaterial, error) {
	return f.materiais, f.err
}

type fakeFornecedorRepo struct {
	data []entity.DimFornecedor
	err  error
}

func (f *fakeFornecedorRepo) FindAll() ([]entity.DimFornecedor, error) { return f.data, f.err }

type fakeMaterialRepo struct {
	data []entity.DimMaterial
	err  error
}

func (f *fakeMaterialRepo) FindAll() ([]entity.DimMaterial, error) { return f.data, f.err }

type fakeResponsavelRepo struct {
	data []entity.DimResponsavel
	err  error
}

func (f *fakeResponsavelRepo) FindAll() ([]entity.DimResponsavel, error) { return f.data, f.err }

type fakeSolicitacaoRepo struct {
	data []entity.DimSolicitacao
	err  error
}

func (f *fakeSolicitacaoRepo) FindAll() ([]entity.DimSolicitacao, error) { return f.data, f.err }

type fakeTarefaRepo struct {
	data []entity.DimTarefa
	err  error
}

func (f *fakeTarefaRepo) FindAll() ([]entity.DimTarefa, error) { return f.data, f.err }

type fakeTempoRepo struct {
	data  []entity.DimTempo
	gasto interface{}
	err   error
}

func (f *fakeTempoRepo) FindAll() ([]entity.DimTempo, error) { return f.data, f.err }
func (f *fakeTempoRepo) GetTempoGasto() (interface{}, error)  { return f.gasto, f.err }

type fakeFatoComprasRepo struct {
	data []entity.FatoCompras
	err  error
}

func (f *fakeFatoComprasRepo) FindAll() ([]entity.FatoCompras, error) { return f.data, f.err }

type fakeFatoEstoqueRepo struct {
	data []entity.FatoEstoqueMateriais
	err  error
}

func (f *fakeFatoEstoqueRepo) FindAll() ([]entity.FatoEstoqueMateriais, error) { return f.data, f.err }

type fakeFatoExecucaoRepo struct {
	data []entity.FatoExecucaoTarefas
	err  error
}

func (f *fakeFatoExecucaoRepo) FindAll() ([]entity.FatoExecucaoTarefas, error) { return f.data, f.err }

type fakeLogImportacaoRepo struct {
	logs   []entity.LogImportacao
	errs   []entity.LogImportacaoErro
	nextID int
	err    error
}

func (f *fakeLogImportacaoRepo) Create(*entity.LogImportacao) (int, error) {
	f.nextID++
	return f.nextID, f.err
}
func (f *fakeLogImportacaoRepo) Update(*entity.LogImportacao) error             { return f.err }
func (f *fakeLogImportacaoRepo) CreateError(*entity.LogImportacaoErro) error    { return f.err }
func (f *fakeLogImportacaoRepo) CreateErrorsBatch([]entity.LogImportacaoErro) error { return f.err }
func (f *fakeLogImportacaoRepo) FindAll() ([]entity.LogImportacao, error)       { return f.logs, f.err }
func (f *fakeLogImportacaoRepo) FindByID(int) (*entity.LogImportacao, error) {
	if len(f.logs) > 0 {
		return &f.logs[0], f.err
	}
	return nil, f.err
}
func (f *fakeLogImportacaoRepo) FindErrorsByLogID(int) ([]entity.LogImportacaoErro, error) {
	return f.errs, f.err
}

type fakePurchaseRepo struct {
	purchases []entity.Purchase
	metrics   entity.PurchaseMetrics
	err       error
}

func (f *fakePurchaseRepo) FindPurchases(context.Context, entity.PurchaseFilter) ([]entity.Purchase, error) {
	return f.purchases, f.err
}
func (f *fakePurchaseRepo) GetMetrics(context.Context, entity.PurchaseFilter) (entity.PurchaseMetrics, error) {
	return f.metrics, f.err
}
