package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"goserver/internal/domain"
	"goserver/internal/middleware"
	"goserver/internal/repository"
	"goserver/internal/service"
	"net/http"
)

type AuthController struct {
	authService *service.AuthService
}

// NewAuthController atua como o 'construtor'
func NewAuthController(db *sql.DB) *AuthController {
	repo := repository.NewUserRepository(db)
	return &AuthController{authService: service.NewAuthService(db, repo)}
}

// Exemplo de método da 'classe'
// RegisterRoutesAuth foi movido para router.go para centralização

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string       `json:"message" example:"Login realizado com sucesso"`
	User    *domain.User `json:"user"`
	Token   string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type RegisterResponse struct {
	Message string `json:"message" example:"Usuário criado com sucesso"`
}

type MeResponse struct {
	Email   string `json:"email" example:"john.doe@example.com"`
	Message string `json:"message" example:"Dados do usuário autenticado"`
}

// Login godoc
// @Summary Login de usuário
// @Description Autentica um usuário e retorna um token JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Credenciais"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao ler JSON"})
		return
	}

	user, token, err := c.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao fazer login"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LoginResponse{Message: "Login realizado com sucesso", User: user, Token: token})
}

// Register godoc
// @Summary Registro de usuário
// @Description Cria um novo usuário.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Dados de cadastro"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 409 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao ler JSON"})
		return
	}

	if _, err := c.authService.Register(r.Context(), req.Name, req.Email, req.Password); err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, service.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict) // 409 Conflict
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao fazer cadastro"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(RegisterResponse{Message: "Usuário criado com sucesso"})
}

// Me godoc
// @Summary Retorna os dados do usuário autenticado
// @Description Retorna o email do usuário autenticado com base no token JWT fornecido.
// @Tags Auth
// @Produce json
// @Success 200 {object} MeResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /me [get]
// @Security BearerAuth
func (c *AuthController) Me(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserEmailKey).(string)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(MeResponse{Email: email, Message: "Dados do usuário autenticado"})
}
