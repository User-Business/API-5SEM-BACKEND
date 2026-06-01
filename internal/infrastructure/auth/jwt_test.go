package auth

import (
	"testing"
	"time"
)

func TestService_GenerateAndParse_Success(t *testing.T) {
	svc := NewService("test-secret", time.Hour)

	token, jti, exp, err := svc.Generate(42, "admin")
	if err != nil {
		t.Fatalf("esperava sem erro ao gerar, recebeu: %v", err)
	}
	if token == "" {
		t.Fatal("esperava token não vazio")
	}
	if jti == "" {
		t.Fatal("esperava jti não vazio")
	}
	if !exp.After(time.Now()) {
		t.Fatal("esperava exp no futuro")
	}

	claims, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("esperava sem erro ao parsear, recebeu: %v", err)
	}
	if claims.Role != "admin" {
		t.Errorf("esperava role admin, recebeu: %s", claims.Role)
	}
	if claims.Subject != "42" {
		t.Errorf("esperava subject 42, recebeu: %s", claims.Subject)
	}
	if claims.ID != jti {
		t.Errorf("esperava jti %s, recebeu: %s", jti, claims.ID)
	}
}

func TestService_Parse_InvalidSignature(t *testing.T) {
	svc := NewService("test-secret", time.Hour)
	token, _, _, _ := svc.Generate(1, "admin")

	other := NewService("outro-secret", time.Hour)
	if _, err := other.Parse(token); err == nil {
		t.Fatal("esperava erro de assinatura inválida, recebeu nil")
	}
}

func TestService_Parse_Expired(t *testing.T) {
	svc := NewService("test-secret", -time.Hour) // já expirado
	token, _, _, _ := svc.Generate(1, "admin")

	if _, err := svc.Parse(token); err == nil {
		t.Fatal("esperava erro de token expirado, recebeu nil")
	}
}
