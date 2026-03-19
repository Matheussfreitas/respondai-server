package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"google.golang.org/genai"
)

var (
	ErrGeminiNotConfigured = errors.New("GEMINI_API_KEY não configurada no ambiente")
	ErrGeminiRateLimited   = errors.New("limite de requisições da Gemini excedido")
)

func Gemini(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", ErrGeminiNotConfigured
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("gemini: falha ao criar client: %w", err)
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-2.5-flash"
	}

	const maxAttempts = 4
	var lastErr error

	for attempt := 0; attempt < maxAttempts; attempt++ {
		result, err := client.Models.GenerateContent(ctx, model, genai.Text(prompt), nil)
		if err == nil {
			return result.Text(), nil
		}

		lastErr = err
		if !isGeminiRateLimitError(err) || attempt == maxAttempts-1 {
			break
		}

		backoff := (250 * time.Millisecond) * time.Duration(1<<attempt)
		time.Sleep(backoff)
	}

	if isGeminiRateLimitError(lastErr) {
		return "", fmt.Errorf("%w: %v", ErrGeminiRateLimited, lastErr)
	}

	return "", fmt.Errorf("gemini: falha ao gerar conteúdo: %w", lastErr)
}

func isGeminiRateLimitError(err error) bool {
	if err == nil {
		return false
	}

	var apiErr genai.APIError
	if errors.As(err, &apiErr) {
		if apiErr.Code == 429 {
			return true
		}
		if apiErr.Status == "RESOURCE_EXHAUSTED" {
			return true
		}
	}

	return false
}
