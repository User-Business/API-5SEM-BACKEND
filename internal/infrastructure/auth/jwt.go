package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims são as claims do JWT da aplicação.
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Service gera e valida JWTs HS256.
type Service struct {
	secret []byte
	expiry time.Duration
}

func NewService(secret string, expiry time.Duration) *Service {
	return &Service{secret: []byte(secret), expiry: expiry}
}

// Generate cria um token assinado para o usuário informado.
// Retorna o token, o jti (id do token) e o instante de expiração.
func (s *Service) Generate(userID int, role string) (string, string, time.Time, error) {
	jti, err := newJTI()
	if err != nil {
		return "", "", time.Time{}, err
	}

	exp := time.Now().Add(s.expiry)
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return signed, jti, exp, nil
}

// Parse valida o token e devolve as claims.
func (s *Service) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func newJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
