package quizzes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"goserver/internal/config"
	"goserver/internal/domain"
	"goserver/internal/repository"
	"time"
)

type CreateQuizService struct {
	repo *repository.QuizRepository
	db   *sql.DB
}

func NewCreateQuizService(repo *repository.QuizRepository, db *sql.DB) *CreateQuizService {
	return &CreateQuizService{
		repo: repo,
		db:   db,
	}
}

type QuizExpected struct {
	QuizTitle string `json:"quiz_title"`
	Questions []struct {
		Statement    string `json:"statement"`
		Alternatives []struct {
			Text string `json:"text"`
		} `json:"alternatives"`
		CorrectIndex int    `json:"correct_index"`
		Explanation  string `json:"explanation"`
	} `json:"questions"`
}

func (s *CreateQuizService) CreateQuizHandler(userID string, numQuestoes int, dificuldade, tema string) (string, error) {
	return s.CreateQuiz(numQuestoes, dificuldade, tema, userID)
}

func (s *CreateQuizService) CreateQuiz(numQuestoes int, dificuldade, tema, userID string) (string, error) {
	prompt := BuildPrompt(tema, numQuestoes, dificuldade)

	quizBuild, err := config.Gemini(prompt)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var quizExpected QuizExpected

	if err := json.Unmarshal([]byte(quizBuild), &quizExpected); err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(quizExpected)

	var questions []domain.Question
	for _, q := range quizExpected.Questions {
		var answers []string
		for _, alt := range q.Alternatives {
			answers = append(answers, alt.Text)
		}
		questions = append(questions, domain.Question{
			Statement:     q.Statement,
			Answers:       answers,
			CorrectAnswer: q.CorrectIndex,
			Explanation:   q.Explanation,
		})
	}

	quiz := domain.Quiz{
		Title:           quizExpected.QuizTitle,
		Content:         tema,
		Questions:       questions,
		Difficulty:      domain.QuizDifficulty(dificuldade),
		NumberQuestions: numQuestoes,
		CreatedAt:       time.Now(),
	}

	quizCreated, err := s.repo.CreateQuiz(userID, quiz)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return quizCreated.ID, nil
}

func BuildPrompt(tema string, numQuestoes int, dificuldade string) string {
	return fmt.Sprintf(`Atue como um professor especialista no tema %s. 
    Gere um quiz com %d questões de nível %s.
    Retorne apenas JSON no formato:
    {
      "quiz_title": string,
      "questions": [
        {
          "statement": string,
          "alternatives": [{"text": string}],
          "correct_index": int,
          "explanation": string
        }
      ]
    }
    Não inclua explicações fora do JSON.`, tema, numQuestoes, dificuldade)
}
