package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// TestRBACEndpointsAdmin verifica, para cada endpoint que alimenta os
// dashboards, as três situações de autorização:
//   - sem token            -> 401 (não autenticado)
//   - token role "compras" -> 403 (autenticado, mas sem permissão)
//   - token role "admin"   -> 200 (autorizado)
//
// Cada endpoint está ligado a uma user story do backlog (ver coluna "US").
func TestRBACEndpointsAdmin(t *testing.T) {
	endpoints := []struct {
		nome string
		path string
		us   string
	}{
		{"projetos", "/api/dim/projetos", "US base (visão de projetos)"},
		{"fornecedores", "/api/dim/fornecedores", "US10 (busca de fornecedores)"},
		{"materiais", "/api/dim/materiais", "US06"},
		{"tempo-gasto", "/api/dim/tempo-gasto", "US07 (tempo por tarefa)"},
		{"fato-compras", "/api/fato/compras", "US base (custo)"},
		{"fato-estoque", "/api/fato/estoque-materiais", "US06 (consumo de material)"},
		{"fato-execucao", "/api/fato/execucao-tarefas", "US07"},
		{"investimento-programa", "/api/programa/investimento", "US05 (investimento por programa)"},
		{"materiais-por-projeto", "/api/projetos/materiais", "US06"},
	}

	for _, ep := range endpoints {
		t.Run(ep.nome, func(t *testing.T) {
			env := newTestEnv(t)

			cases := []struct {
				situacao   string
				token      string
				wantStatus int
			}{
				{"sem token", "", http.StatusUnauthorized},
				{"role compras", env.tokenFor(entity.RoleCompras), http.StatusForbidden},
				{"role admin", env.tokenFor(entity.RoleAdmin), http.StatusOK},
			}

			for _, c := range cases {
				t.Run(c.situacao, func(t *testing.T) {
					resp := env.request(t, http.MethodGet, ep.path, c.token, nil)
					defer resp.Body.Close()

					if resp.StatusCode != c.wantStatus {
						t.Fatalf("%s [%s]: esperado %d, obtido %d",
							ep.path, c.situacao, c.wantStatus, resp.StatusCode)
					}
				})
			}
		})
	}
}

// TestPurchasesPermitidoParaComprasEAdmin verifica que /api/purchases aceita
// tanto admin quanto o perfil de compras (rota compartilhada), mas continua
// negando quem não está autenticado.
func TestPurchasesPermitidoParaComprasEAdmin(t *testing.T) {
	env := newTestEnv(t)

	cases := []struct {
		situacao   string
		token      string
		wantStatus int
	}{
		{"sem token", "", http.StatusUnauthorized},
		{"role admin", env.tokenFor(entity.RoleAdmin), http.StatusOK},
		{"role compras", env.tokenFor(entity.RoleCompras), http.StatusOK},
	}

	for _, c := range cases {
		t.Run(c.situacao, func(t *testing.T) {
			resp := env.request(t, http.MethodGet, "/api/purchases", c.token, nil)
			defer resp.Body.Close()
			if resp.StatusCode != c.wantStatus {
				t.Fatalf("/api/purchases [%s]: esperado %d, obtido %d",
					c.situacao, c.wantStatus, resp.StatusCode)
			}
		})
	}
}

// TestInvestimentoPorProgramaRetornaDados é o cenário "feliz" da US05:
// Dado dados de investimento por programa no banco, Quando um admin consulta
// /api/programa/investimento, Então recebe 200 e o JSON com os investimentos.
func TestInvestimentoPorProgramaRetornaDados(t *testing.T) {
	env := newTestEnv(t)

	// Arrange: semeia o repositório com um programa de investimento.
	env.projeto.investimentos = []entity.ProgramaInvestimento{
		{CodigoPrograma: "PRG-01", NomePrograma: "Programa Alfa", InvestimentoTotal: 1500.50},
	}

	// Act
	resp := env.request(t, http.MethodGet, "/api/programa/investimento",
		env.tokenFor(entity.RoleAdmin), nil)
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: esperado 200, obtido %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: esperado application/json, obtido %q", ct)
	}

	var got []entity.ProgramaInvestimento
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decodificar resposta: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("esperava 1 programa, obtido %d", len(got))
	}
	if got[0].CodigoPrograma != "PRG-01" || got[0].InvestimentoTotal != 1500.50 {
		t.Errorf("dados inesperados: %+v", got[0])
	}
}

// TestErroNoRepositorioVira500 garante o contrato de erro: se a camada de
// dados falha, o endpoint responde 500 (e não vaza pânico). Aqui forçamos o
// repositório a devolver erro.
func TestErroNoRepositorioVira500(t *testing.T) {
	env := newTestEnv(t)
	env.projeto.err = errInjetado

	resp := env.request(t, http.MethodGet, "/api/dim/projetos",
		env.tokenFor(entity.RoleAdmin), nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status: esperado 500, obtido %d", resp.StatusCode)
	}
}
