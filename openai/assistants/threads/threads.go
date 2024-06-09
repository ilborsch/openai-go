package threads

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ThreadClient interface {
	CreateThread() (string, error)
}

// Threads represents OpenAI API thread domain
type Threads struct {
	APIKey string
}

// CreateThreadResponse is used to unmarshal OpenAI API response
type CreateThreadResponse struct {
	ThreadID string `json:"id"`
}

// CreateThread creates an empty thread object.
// Returns its ID
func (t Threads) CreateThread() (string, error) {
	const URL = "https://api.openai.com/v1/threads"

	request, err := http.NewRequest(http.MethodPost, URL, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+t.APIKey)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error creating thread: %v %s ", resp.StatusCode, string(responseBody))
	}

	var response CreateThreadResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	return response.ThreadID, nil
}
