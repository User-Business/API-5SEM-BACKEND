package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
	"github.com/DenariusData/API-5SEM-BACKEND/internal/infrastructure/auth"
)

type mockUserRepo struct {
	user *entity.Usuario
	err  error
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	return m.user, m.err
}

type mockTokenRepo struct {
	revokedJTI string
	err        error
}

func (m *mockTokenRepo) Revoke(ctx context.Context, jti string, exp time.Time) error {
	m.revokedJTI = jti
	return m.err
}

func newTestAuthUC(user *entity.Usuario, userErr error, tokenRepo *mockTokenRepo) *AuthUseCase {
	jwtSvc := auth.NewService("test-secret", time.Hour)
	return NewAuthUseCase(&mockUserRepo{user: user, err: userErr}, tokenRepo, jwtSvc)
}

func TestAuthUseCase_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("segredo"), bcrypt.DefaultCost)
	user := &entity.Usuario{ID: 1, Email: "a@b.com", SenhaHash: string(hash), Nome: "Ana", Role: entity.RoleAdmin}
	uc := newTestAuthUC(user, nil, &mockTokenRepo{})

	res, err := uc.Login(context.Background(), "a@b.com", "segredo")
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if res.Token == "" {
		t.Error("esperava token não vazio")
	}
	if res.Usuario.Role != entity.RoleAdmin {
		t.Errorf("esperava role admin, recebeu: %s", res.Usuario.Role)
	}
	if res.Usuario.Nome != "Ana" {
		t.Errorf("esperava nome Ana, recebeu: %s", res.Usuario.Nome)
	}
}

func TestAuthUseCase_Login_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("segredo"), bcrypt.DefaultCost)
	user := &entity.Usuario{ID: 1, Email: "a@b.com", SenhaHash: string(hash), Role: entity.RoleAdmin}
	uc := newTestAuthUC(user, nil, &mockTokenRepo{})

	_, err := uc.Login(context.Background(), "a@b.com", "errada")
	if !errors.Is(err, ErrCredenciaisInvalidas) {
		t.Fatalf("esperava ErrCredenciaisInvalidas, recebeu: %v", err)
	}
}

func TestAuthUseCase_Login_UserNotFound(t *testing.T) {
	uc := newTestAuthUC(nil, nil, &mockTokenRepo{})

	_, err := uc.Login(context.Background(), "naoexiste@b.com", "x")
	if !errors.Is(err, ErrCredenciaisInvalidas) {
		t.Fatalf("esperava ErrCredenciaisInvalidas, recebeu: %v", err)
	}
}

func TestAuthUseCase_Logout_RevokesToken(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	uc := newTestAuthUC(nil, nil, tokenRepo)

	err := uc.Logout(context.Background(), "jti-123", time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("esperava sem erro, recebeu: %v", err)
	}
	if tokenRepo.revokedJTI != "jti-123" {
		t.Errorf("esperava jti revogado jti-123, recebeu: %s", tokenRepo.revokedJTI)
	}
}
