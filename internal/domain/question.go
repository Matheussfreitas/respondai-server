package domain

type Question struct {
	ID            string   `db:"id" json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	QuizID        string   `db:"quiz_id" json:"quiz_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Statement     string   `db:"statement" json:"statement" example:"Qual é a capital da França?"`
	Answers       []string `db:"answers" json:"answers" example:"[\"Paris\",\"Londres\",\"Berlim\",\"Roma\"]"` // Lista de opções (ex: Paris, Londres)
	CorrectAnswer int      `db:"correct_answer" json:"correct_answer" example:"0"`                             // Índice da resposta correta
	Explanation   string   `db:"explanation" json:"explanation" example:"Paris é a capital da França."`
}
