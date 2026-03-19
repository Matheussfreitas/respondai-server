package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"goserver/internal/config"
	"goserver/internal/domain"
	"goserver/internal/repository"
	"goserver/internal/service/quizzes"
	"net/http"
	"time"
)

type QuizHandler struct {
	findManyQuizzesService *quizzes.FindManyQuizzesService
	createQuizService      *quizzes.CreateQuizService
	findQuizByIdService    *quizzes.FindQuizByIdService
	submitQuizService      *quizzes.SubmitQuizService
}

func NewQuizHandler(db *sql.DB) *QuizHandler {
	repo := repository.NewQuizRepository(db)
	return &QuizHandler{
		findManyQuizzesService: quizzes.NewFindManyQuizzesService(repo, db),
		createQuizService:      quizzes.NewCreateQuizService(repo, db),
		findQuizByIdService:    quizzes.NewFindQuizByIdService(repo, db),
		submitQuizService:      quizzes.NewSubmitQuizService(repo, db),
	}
}

type CreateQuizRequest struct {
	Tema        string `json:"tema"`
	NumQuestoes int    `json:"numQuestoes"`
	Dificuldade string `json:"dificuldade"`
}

type CreateQuizResponse struct {
	Message string `json:"message" example:"Quiz criado com sucesso"`
	Quiz    string `json:"quiz" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type SubmitQuizRequest struct {
	QuizID  string              `json:"quiz_id"`
	UserID  string              `json:"user_id"`
	Answers []domain.UserAnswer `json:"answers"`
}

type SubmitQuizResponse struct {
	Message string             `json:"message" example:"Quiz enviado com sucesso"`
	Quiz    *domain.ResultQuiz `json:"quiz"`
}

// CreateQuiz godoc
// @Summary Cria um novo quiz
// @Description Gera um quiz com IA e salva no banco.
// @Tags Quizzes
// @Accept json
// @Produce json
// @Param body body CreateQuizRequest true "Parâmetros do quiz"
// @Success 200 {object} CreateQuizResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /quizzes/create [post]
// @Security BearerAuth
func (h *QuizHandler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var req CreateQuizRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao ler JSON"})
		return
	}

	userId := r.Context().Value("user_id").(string)

	quiz, err := h.createQuizService.CreateQuiz(req.NumQuestoes, req.Dificuldade, req.Tema, userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, config.ErrGeminiRateLimited) {
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Limite da IA atingido (429). Tente novamente em instantes."})
			return
		}

		if errors.Is(err, config.ErrGeminiNotConfigured) {
			w.WriteHeader(http.StatusFailedDependency)
			_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Integração com IA não configurada no servidor."})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao criar quiz"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(CreateQuizResponse{Message: "Quiz criado com sucesso", Quiz: quiz})
}

// SubmitQuiz godoc
// @Summary Submete um quiz respondido
// @Description Submete as respostas de um quiz e calcula a pontuação do usuário.
// @Tags Quizzes
// @Accept json
// @Produce json
// @Param body body SubmitQuizRequest true "Respostas do quiz"
// @Success 200 {object} SubmitQuizResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /quizzes/submit [post]
// @Security BearerAuth
func (h *QuizHandler) SubmitQuiz(w http.ResponseWriter, r *http.Request) {
	var req SubmitQuizRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao ler JSON"})
		return
	}

	resultQuiz := domain.ResultQuiz{
		QuizID:      req.QuizID,
		UserID:      req.UserID,
		Answers:     req.Answers,
		CompletedAt: time.Now(),
	}

	userId := r.Context().Value("user_id").(string)
	quiz, err := h.submitQuizService.SubmitQuiz(resultQuiz, userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao enviar quiz"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SubmitQuizResponse{Message: "Quiz enviado com sucesso", Quiz: quiz})
}

// FindManyQuizzes godoc
// @Summary Lista quizzes do usuário
// @Description Retorna todos os quizzes do usuário autenticado.
// @Tags Quizzes
// @Produce json
// @Success 200 {array} domain.Quiz
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /quizzes [get]
// @Security BearerAuth
func (h *QuizHandler) FindManyQuizzes(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(string)

	quizzes, err := h.findManyQuizzesService.FindManyQuizzes(userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao buscar quizzes"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(quizzes)
}

// FindQuizById godoc
// @Summary Busca quiz por ID
// @Description Retorna um quiz específico do usuário autenticado.
// @Tags Quizzes
// @Produce json
// @Param id path string true "ID do Quiz"
// @Success 200 {object} domain.Quiz
// @Failure 401 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /quizzes/{id} [get]
// @Security BearerAuth
func (h *QuizHandler) FindQuizById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId := r.Context().Value("user_id").(string)

	quiz, err := h.findQuizByIdService.FindQuizById(id, userId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(domain.ErrorResponse{Message: "Erro ao buscar quiz"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(quiz)
}
