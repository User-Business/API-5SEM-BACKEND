package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/adapter/handler"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type Handlers struct {
	Projeto      *handler.ProjetoHandler
	Fornecedor   *handler.FornecedorHandler
	Material     *handler.MaterialHandler
	Responsavel  *handler.ResponsavelHandler
	Solicitacao  *handler.SolicitacaoHandler
	Tarefa       *handler.TarefaHandler
	Tempo        *handler.TempoHandler
	FatoCompras  *handler.FatoComprasHandler
	FatoEstoque  *handler.FatoEstoqueHandler
	FatoExecucao *handler.FatoExecucaoHandler
	Purchase     *handler.PurchaseHandler
}

func NewRouter(h Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors)

	r.Route("/api", func(r chi.Router) {
		r.Route("/dim", func(r chi.Router) {
			r.Get("/projetos", h.Projeto.GetAll)
			r.Get("/fornecedores", h.Fornecedor.GetAll)
			r.Get("/materiais", h.Material.GetAll)
			r.Get("/responsaveis", h.Responsavel.GetAll)
			r.Get("/solicitacoes", h.Solicitacao.GetAll)
			r.Get("/tarefas", h.Tarefa.GetAll)
			r.Get("/tempo", h.Tempo.GetAll)
			r.Get("/tempo-gasto", h.Tempo.GetTempoGasto)
		})

		r.Route("/fato", func(r chi.Router) {
			r.Get("/compras", h.FatoCompras.GetAll)
			r.Get("/estoque-materiais", h.FatoEstoque.GetAll)
			r.Get("/execucao-tarefas", h.FatoExecucao.GetAll)
		})

		r.Route("/programa", func(r chi.Router) {
			r.Get("/investimento", h.Projeto.GetInvestimentoByPrograma)
		})

		r.Route("/projetos", func(r chi.Router) {
			r.Get("/materiais", h.Projeto.GetMateriaisByProjeto)
		})

		r.Route("/purchases", func(r chi.Router) {
			r.Get("/", h.Purchase.GetPurchases)
			r.Get("/metrics", h.Purchase.GetMetrics)
		})
	})

	return r
}
