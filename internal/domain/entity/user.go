package entity

import "time"

// Roles disponíveis no sistema.
const (
	RoleAdmin   = "admin"
	RoleCompras = "compras"
)

// Usuario representa um usuário do sistema (registro completo, uso interno).
type Usuario struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	SenhaHash string    `json:"-"`
	Nome      string    `json:"nome"`
	Role      string    `json:"role"`
	CriadoEm  time.Time `json:"criado_em"`
}

// UsuarioPublico é a representação retornada ao cliente (sem dados sensíveis).
type UsuarioPublico struct {
	Nome string `json:"nome"`
	Role string `json:"role"`
}

// LoginRequest é o corpo esperado em POST /api/auth/login.
type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// LoginResponse é o corpo retornado por um login bem-sucedido.
type LoginResponse struct {
	Token   string         `json:"token"`
	Usuario UsuarioPublico `json:"usuario"`
}
