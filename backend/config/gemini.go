package config

import (
	"context"
	"errors"

	"os"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	GeminiClient *genai.Client
	geminiOnce   sync.Once
	geminiErr    error
)

func InitGemini() error {
	geminiOnce.Do(func() {
		GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))
		if err != nil {
			geminiErr = err
			return
		}
		GeminiClient = client
	})

	return geminiErr
}

// GetGeminiClient returns the initialized Gemini client or an error if not initialized
func GetGeminiClient() (*genai.Client, error) {
	if GeminiClient == nil {
		return nil, errors.New("gemini client is not initialized")
	}
	return GeminiClient, nil
}
