package domain

import "time"

type Quiz struct {
	ID              string         `db:"id" json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID          string         `db:"user_id" json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title           string         `db:"title" json:"title" example:"Capitais do Mundo"`
	Content         string         `db:"content" json:"content" example:"Quiz rápido sobre capitais."`
	Difficulty      QuizDifficulty `db:"difficulty" json:"difficulty" enums:"easy,medium,hard" example:"easy"`
	NumberQuestions int            `db:"number_questions" json:"number_questions" example:"10"`
	Questions       []Question     `db:"-" json:"questions,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at" example:"2026-01-01T12:00:00Z"`
}

type QuizDifficulty string

const (
	Easy   QuizDifficulty = "easy"
	Medium QuizDifficulty = "medium"
	Hard   QuizDifficulty = "hard"
)
