package quizzes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"goserver/internal/config"
	"goserver/internal/domain"
	"goserver/internal/repository"
	"strings"
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

	quizJSON, err := extractJSONFromModelOutput(quizBuild)
	if err != nil {
		return "", fmt.Errorf("resposta da IA sem JSON válido: %w", err)
	}

	if err := json.Unmarshal([]byte(quizJSON), &quizExpected); err != nil {
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

func extractJSONFromModelOutput(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("resposta vazia")
	}

	if json.Valid([]byte(trimmed)) {
		return trimmed, nil
	}

	// Common model behavior: wrapping JSON in ```json ... ```
	trimmed = stripMarkdownCodeFence(trimmed)
	if json.Valid([]byte(trimmed)) {
		return trimmed, nil
	}

	for i := 0; i < len(trimmed); i++ {
		if trimmed[i] != '{' && trimmed[i] != '[' {
			continue
		}

		var candidate json.RawMessage
		decoder := json.NewDecoder(strings.NewReader(trimmed[i:]))
		if err := decoder.Decode(&candidate); err == nil && json.Valid(candidate) {
			return string(candidate), nil
		}
	}

	return "", errors.New("não foi possível extrair JSON")
}

func stripMarkdownCodeFence(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) < 2 {
		return input
	}

	first := strings.TrimSpace(lines[0])
	last := strings.TrimSpace(lines[len(lines)-1])
	if !strings.HasPrefix(first, "```") || last != "```" {
		return input
	}

	return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
}
