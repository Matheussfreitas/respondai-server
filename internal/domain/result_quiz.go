package domain

import "time"

type ResultQuiz struct {
	ID             string       `db:"id" json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	QuizID         string       `db:"quiz_id" json:"quiz_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID         string       `db:"user_id" json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Score          int          `db:"score" json:"score" example:"7"`
	TotalQuestions int          `db:"total_questions" json:"total_questions" example:"10"`
	Answers        []UserAnswer `db:"answers" json:"answers"`
	CompletedAt    time.Time    `db:"completed_at" json:"completed_at" example:"2026-01-01T12:00:00Z"`
}

type UserAnswer struct {
	QuestionID string `json:"question_id" example:"123e4567-e89b-12d3-a456-426614174000"` // UUID
	UserChoice int    `json:"user_choice" example:"2"`
	IsCorrect  bool   `json:"is_correct" example:"false"`
}
