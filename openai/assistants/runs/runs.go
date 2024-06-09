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

type RunClient interface {
	CreateRun(threadID string, assistantID string) (string, error)
	GetRun(threadID, runID string) (GetRunResponse, error)
}

// Runs represents OpenAI API run domain
type Runs struct {
	APIKey string
}

// CreateRunRequest is used to structure payload in the Create function
type CreateRunRequest struct {
	AssistantID string `json:"assistant_id"`
}

// CreateRunResponse is used to unmarshal OpenAI API response
type CreateRunResponse struct {
	RunID string `json:"id"`
}

// GetRunResponse is used to structure payload in the GetRun function
type GetRunResponse struct {
	Status      string `json:"status"`
	ThreadID    string `json:"thread_id"`
	AssistantID string `json:"assistant_id"`
}

// CreateRun creates run object for an assistant with the `assistantID` in the thread specified by `threadID`.
// Returns its ID.
// Recommended docs to better understand this approach: https://platform.openai.com/docs/assistants/overview
func (r Runs) CreateRun(threadID string, assistantID string) (string, error) {
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

// GetRun fetches run object given by `runID` parameter.
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
