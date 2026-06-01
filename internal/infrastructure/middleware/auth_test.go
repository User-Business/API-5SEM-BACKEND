package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/auth"
)

type fakeParser struct {
	claims *auth.Claims
	err    error
}

func (f *fakeParser) Parse(token string) (*auth.Claims, error) {
	return f.claims, f.err
}

type fakeRevoker struct {
	revoked bool
	err     error
}

func (f *fakeRevoker) IsRevoked(ctx context.Context, jti string) (bool, error) {
	return f.revoked, f.err
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func validClaims() *auth.Claims {
	c := &auth.Claims{Role: "admin"}
	c.ID = "jti-1"
	c.Subject = "1"
	return c
}

func TestAuth_NoHeader_401(t *testing.T) {
	mw := Auth(&fakeParser{claims: validClaims()}, &fakeRevoker{})
	req := httptest.NewRequest(http.MethodGet, "/api/dim/projetos", nil)
	rec := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, recebeu: %d", rec.Code)
	}
}

func TestAuth_InvalidToken_401(t *testing.T) {
	mw := Auth(&fakeParser{err: errors.New("inválido")}, &fakeRevoker{})
	req := httptest.NewRequest(http.MethodGet, "/api/dim/projetos", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	rec := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, recebeu: %d", rec.Code)
	}
}

func TestAuth_RevokedToken_401(t *testing.T) {
	mw := Auth(&fakeParser{claims: validClaims()}, &fakeRevoker{revoked: true})
	req := httptest.NewRequest(http.MethodGet, "/api/dim/projetos", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	rec := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("esperava 401, recebeu: %d", rec.Code)
	}
}

func TestAuth_ValidToken_PassesAndSetsClaims(t *testing.T) {
	mw := Auth(&fakeParser{claims: validClaims()}, &fakeRevoker{})
	req := httptest.NewRequest(http.MethodGet, "/api/dim/projetos", nil)
	req.Header.Set("Authorization", "Bearer xxx")
	rec := httptest.NewRecorder()

	var gotRole string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := ClaimsFromContext(r.Context()); ok {
			gotRole = c.Role
		}
		w.WriteHeader(http.StatusOK)
	})

	mw(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, recebeu: %d", rec.Code)
	}
	if gotRole != "admin" {
		t.Errorf("esperava role admin no contexto, recebeu: %s", gotRole)
	}
}

func TestRequireRole_Forbidden_403(t *testing.T) {
	mw := RequireRole("admin")
	req := httptest.NewRequest(http.MethodGet, "/api/dim/projetos", nil)
	c := &auth.Claims{Role: "compras"}
	req = req.WithContext(context.WithValue(req.Context(), claimsKey, c))
	rec := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("esperava 403, recebeu: %d", rec.Code)
	}
}

func TestRequireRole_Allowed_200(t *testing.T) {
	mw := RequireRole("admin", "compras")
	req := httptest.NewRequest(http.MethodGet, "/api/purchases/", nil)
	c := &auth.Claims{Role: "compras"}
	req = req.WithContext(context.WithValue(req.Context(), claimsKey, c))
	rec := httptest.NewRecorder()

	mw(okHandler()).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava 200, recebeu: %d", rec.Code)
	}
}
