package runs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	StatusQueued         = "queued"
	StatusInProgress     = "in_progress"
	StatusRequiresAction = "requires_action"
	StatusCancelling     = "cancelling"
	StatusCancelled      = "cancelled"
	StatusFailed         = "failed"
	StatusCompleted      = "completed"
	StatusExpired        = "expired"
)

type Runs struct {
	APIKey string
}

type CreateRunRequest struct {
	AssistantID string `json:"assistant_id"`
}

type CreateRunResponse struct {
	RunID string `json:"id"`
}

type GetRunResponse struct {
	Status      string `json:"status"`
	ThreadID    string `json:"thread_id"`
	AssistantID string `json:"assistant_id"`
}

func (r Runs) Create(threadID string, assistantID string) (string, error) {
	URL := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", threadID)

	payload := CreateRunRequest{
		AssistantID: assistantID,
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create run: %d %s", resp.StatusCode, string(responseBody))
	}

	var response CreateRunResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	return response.RunID, nil
}

func (r Runs) GetRun(threadID, runID string) (GetRunResponse, error) {
	URL := fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs/%s", threadID, runID)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return GetRunResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetRunResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetRunResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GetRunResponse{}, fmt.Errorf("error retrieving run: %v %s ", resp.StatusCode, string(responseBody))
	}

	var response GetRunResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return GetRunResponse{}, fmt.Errorf("failed to marshal response body: %v", err)
	}
	return response, nil
}
