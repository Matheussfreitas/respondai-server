package quizzes

import (
	"database/sql"
	"goserver/internal/domain"
	"goserver/internal/repository"
)

type FindManyQuizzesService struct {
	db   *sql.DB
	repo *repository.QuizRepository
}

func NewFindManyQuizzesService(repo *repository.QuizRepository, db *sql.DB) *FindManyQuizzesService {
	return &FindManyQuizzesService{
		repo: repo,
		db:   db,
	}
}

func (s *FindManyQuizzesService) FindManyQuizzesHandler(userId string) ([]domain.Quiz, error) {
	return s.FindManyQuizzes(userId)
}

func (s *FindManyQuizzesService) FindManyQuizzes(userId string) ([]domain.Quiz, error) {
	return s.repo.FindManyQuizzes(userId)
}
