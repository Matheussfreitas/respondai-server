package handler

import (
	"database/sql"
	"goserver/internal/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Routes struct {
	mux  *http.ServeMux
	auth *AuthController
	quiz *QuizHandler
}

func NewRouter(db *sql.DB) *Routes {
	return &Routes{
		mux:  http.NewServeMux(),
		auth: NewAuthController(db),
		quiz: NewQuizHandler(db),
	}
}

func (r *Routes) GetHandler() http.Handler {
	return r.mux
}

func (r *Routes) RegisterRoutes() {
	// Rotas Públicas
	r.mux.HandleFunc("GET /health", HealthCheck)
	r.mux.HandleFunc("POST /login", r.auth.Login)
	r.mux.HandleFunc("POST /register", r.auth.Register)
	r.mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	r.mux.HandleFunc("GET /reference", ScalarReferenceHandler("/swagger/doc.json"))
	r.mux.HandleFunc("GET /reference/", ScalarReferenceHandler("/swagger/doc.json"))

	// Rotas Protegidas
	r.mux.Handle("GET /me", middleware.AuthMiddleware(http.HandlerFunc(r.auth.Me)))
	r.mux.Handle("GET /quizzes", middleware.AuthMiddleware(http.HandlerFunc(r.quiz.FindManyQuizzes)))
	r.mux.Handle("GET /quizzes/{id}", middleware.AuthMiddleware(http.HandlerFunc(r.quiz.FindQuizById)))
	r.mux.Handle("POST /quizzes/create", middleware.AuthMiddleware(http.HandlerFunc(r.quiz.CreateQuiz)))
	r.mux.Handle("POST /quizzes/submit", middleware.AuthMiddleware(http.HandlerFunc(r.quiz.SubmitQuiz)))
}
