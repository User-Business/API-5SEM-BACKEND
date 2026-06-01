package usecase

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/DenariusData/API-5SEM-BACKEND/internal/domain/entity"
)

// ErrCredenciaisInvalidas é retornado quando email ou senha não conferem.
var ErrCredenciaisInvalidas = errors.New("credenciais inválidas")

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.Usuario, error)
}

type TokenRepository interface {
	Revoke(ctx context.Context, jti string, exp time.Time) error
}

type TokenGenerator interface {
	Generate(userID int, role string) (token string, jti string, exp time.Time, err error)
}

type AuthUseCase struct {
	users  UserRepository
	tokens TokenRepository
	jwt    TokenGenerator
}

func NewAuthUseCase(users UserRepository, tokens TokenRepository, jwt TokenGenerator) *AuthUseCase {
	return &AuthUseCase{users: users, tokens: tokens, jwt: jwt}
}

// Login valida credenciais e devolve um token JWT.
func (uc *AuthUseCase) Login(ctx context.Context, email, senha string) (entity.LoginResponse, error) {
	user, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		return entity.LoginResponse{}, err
	}
	if user == nil {
		return entity.LoginResponse{}, ErrCredenciaisInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.SenhaHash), []byte(senha)); err != nil {
		return entity.LoginResponse{}, ErrCredenciaisInvalidas
	}

	token, _, _, err := uc.jwt.Generate(user.ID, user.Role)
	if err != nil {
		return entity.LoginResponse{}, err
	}

	return entity.LoginResponse{
		Token:   token,
		Usuario: entity.UsuarioPublico{Nome: user.Nome, Role: user.Role},
	}, nil
}

// Logout revoga o token (adiciona o jti à denylist).
func (uc *AuthUseCase) Logout(ctx context.Context, jti string, exp time.Time) error {
	return uc.tokens.Revoke(ctx, jti, exp)
}
