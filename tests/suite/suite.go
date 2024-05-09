package suite

import (
	"github.com/ilborsch/openai-go/openai"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

type Suite struct {
	*testing.T
	Client *openai.OpenAI
}

func New(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	APIKey := mustLoadAPIKey()
	client := openai.New(APIKey)
	return &Suite{
		T:      t,
		Client: client,
	}
}

func mustLoadAPIKey() string {
	if err := godotenv.Load("../../tests/test-data/.env"); err != nil {
		panic("error loading .env file: " + err.Error())
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY not found in .env file")
	}
	return apiKey
}
