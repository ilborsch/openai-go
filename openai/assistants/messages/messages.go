package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Messages struct {
	APIKey string
}

type AddMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ThreadMessages struct {
	Data []ThreadMessage `json:"data"`
}

type ThreadMessage struct {
	Content []ThreadMessageContent `json:"content"`
}

type ThreadMessageContent struct {
	Text MessageValue `json:"text"`
}

type MessageValue struct {
	Value string `json:"value"`
}

func (m Messages) AddMessageToThread(threadID string, message string) error {
	URL := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", threadID)

	payload := AddMessageRequest{
		Role:    "user",
		Content: message,
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.APIKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create message: %d %s ", resp.StatusCode, string(responseBody))
	}
	return nil
}

func (m Messages) GetThreadMessages(threadID string) (ThreadMessages, error) {
	URL := fmt.Sprintf("https://api.openai.com/v1/threads/%s/messages", threadID)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return ThreadMessages{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.APIKey)
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ThreadMessages{}, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ThreadMessages{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ThreadMessages{}, fmt.Errorf("failed to retrieve messages: %d %s", resp.StatusCode, string(responseBody))
	}

	var response ThreadMessages
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return ThreadMessages{}, err
	}
	return response, nil
}

func (m Messages) LatestAssistantResponse(threadID string) (string, error) {
	chatStory, err := m.GetThreadMessages(threadID)
	if err != nil {
		return "", err
	}
	if len(chatStory.Data) == 0 {
		return "", fmt.Errorf("thread %s has no messages", threadID)
	}
	return chatStory.Data[0].Content[0].Text.Value, nil
}
