package config

import (
	"context"
	"errors"
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
		// apiKey := os.Getenv("GEMINI_API_KEY")
		// if apiKey == "" {
		// 	geminiErr = errors.New("GEMINI_API_KEY environment variable is not set")
		// 	return
		// }

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyAIFYLl-0rfX89tignMqQLNR-SccQnD31Q"))
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
