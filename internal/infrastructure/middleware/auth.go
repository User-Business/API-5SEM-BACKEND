package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/auth"
)

type contextKey string

const claimsKey contextKey = "claims"

// TokenParser valida um token e devolve suas claims.
type TokenParser interface {
	Parse(token string) (*auth.Claims, error)
}

// Revoker informa se um jti foi revogado (denylist).
type Revoker interface {
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

// ClaimsFromContext recupera as claims injetadas pelo middleware Auth.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*auth.Claims)
	return c, ok
}

func unauthorized(w http.ResponseWriter) {
	http.Error(w, `{"error":"não autorizado"}`, http.StatusUnauthorized)
}

// Auth valida o Bearer token, checa a denylist e injeta as claims no contexto.
func Auth(parser TokenParser, revoker Revoker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				unauthorized(w)
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")

			claims, err := parser.Parse(tokenStr)
			if err != nil {
				unauthorized(w)
				return
			}

			revoked, err := revoker.IsRevoked(r.Context(), claims.ID)
			if err != nil {
				http.Error(w, `{"error":"erro interno"}`, http.StatusInternalServerError)
				return
			}
			if revoked {
				unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole bloqueia (403) requisições cujo usuário não tenha uma das roles permitidas.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				unauthorized(w)
				return
			}
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, `{"error":"acesso negado"}`, http.StatusForbidden)
		})
	}
}
