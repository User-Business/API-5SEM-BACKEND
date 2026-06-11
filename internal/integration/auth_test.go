package integration

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// TestLogin cobre o endpoint público POST /api/auth/login.
//
// Cenários (Dado/Quando/Então):
//   - TC-AUTH-01: Dado um usuário válido, Quando faz login com a senha correta,
//     Então recebe 200 e um token JWT.
//   - TC-AUTH-02: Dado um usuário válido, Quando erra a senha, Então recebe 401.
//   - TC-AUTH-03: Dado um e-mail inexistente, Quando tenta logar, Então recebe 401.
//   - TC-AUTH-04: Dado um corpo malformado, Quando envia ao login, Então recebe 400.
func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantToken  bool
	}{
		{
			name:       "TC-AUTH-01 sucesso",
			body:       `{"email":"` + defaultEmail + `","senha":"` + defaultPassword + `"}`,
			wantStatus: http.StatusOK,
			wantToken:  true,
		},
		{
			name:       "TC-AUTH-02 senha incorreta",
			body:       `{"email":"` + defaultEmail + `","senha":"errada"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "TC-AUTH-03 email inexistente",
			body:       `{"email":"ninguem@teste.com","senha":"qualquer"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "TC-AUTH-04 corpo malformado",
			body:       `{ isso nao e json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env := newTestEnv(t)

			resp := env.request(t, http.MethodPost, "/api/auth/login", "", strings.NewReader(tc.body))
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("status: esperado %d, obtido %d", tc.wantStatus, resp.StatusCode)
			}

			if tc.wantToken {
				var got entity.LoginResponse
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatalf("decodificar resposta: %v", err)
				}
				if got.Token == "" {
					t.Error("esperava um token JWT, veio vazio")
				}
				if got.Usuario.Role != entity.RoleAdmin {
					t.Errorf("role: esperado %q, obtido %q", entity.RoleAdmin, got.Usuario.Role)
				}
			}
		})
	}
}

// TestRotaProtegidaSemToken garante que uma rota autenticada rejeita (401)
// requisições sem token Bearer. (TC-AUTH-05)
func TestRotaProtegidaSemToken(t *testing.T) {
	env := newTestEnv(t)

	resp := env.request(t, http.MethodGet, "/api/dim/projetos", "", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status: esperado 401, obtido %d", resp.StatusCode)
	}
}

// TestRotaProtegidaComTokenInvalido garante que um token mal formado é
// rejeitado com 401. (TC-AUTH-06)
func TestRotaProtegidaComTokenInvalido(t *testing.T) {
	env := newTestEnv(t)

	resp := env.request(t, http.MethodGet, "/api/dim/projetos", "token-invalido-xyz", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status: esperado 401, obtido %d", resp.StatusCode)
	}
}

// TestLogoutRevogaToken verifica o fluxo de logout: após o logout, o mesmo
// token passa a ser rejeitado (denylist). (TC-AUTH-07)
func TestLogoutRevogaToken(t *testing.T) {
	env := newTestEnv(t)
	token := env.tokenFor(entity.RoleAdmin)

	// Antes do logout: token funciona.
	resp := env.request(t, http.MethodGet, "/api/dim/projetos", token, nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("antes do logout esperava 200, obtido %d", resp.StatusCode)
	}

	// Logout.
	resp = env.request(t, http.MethodPost, "/api/auth/logout", token, nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("logout esperava 200, obtido %d", resp.StatusCode)
	}

	// Depois do logout: mesmo token é rejeitado.
	resp = env.request(t, http.MethodGet, "/api/dim/projetos", token, nil)
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("depois do logout esperava 401, obtido %d", resp.StatusCode)
	}
}
