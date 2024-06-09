package messages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type MessageClient interface {
	AddMessageToThread(threadID string, message string) error
	GetThreadMessages(threadID string) (ThreadMessages, error)
	LatestAssistantResponse(threadID string) (string, error)
}

// Messages represents OpenAI API thread message domain
type Messages struct {
	APIKey string
}

// AddMessageRequest is used to structure payload in the AddMessageToThread request
type AddMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ThreadMessages is used to unmarshal OpenAI API response in the GetThreadMessages function
// ThreadMessages represents a list of all thread message which you could also call "chat story"
type ThreadMessages struct {
	Data []ThreadMessage `json:"data"`
}

// ThreadMessage is used to unmarshal OpenAI API response in the GetThreadMessages function
// ThreadMessage represents a Thread message entity of ThreadMessages list
type ThreadMessage struct {
	Role    string                 `json:"role"`
	Content []ThreadMessageContent `json:"content"`
}

// ThreadMessageContent is used to unmarshal OpenAI API response in the GetThreadMessages function
// ThreadMessageContent represents a contents of ThreadMessage struct
type ThreadMessageContent struct {
	Text MessageValue `json:"text"`
}

// MessageValue is used to unmarshal OpenAI API response in the GetThreadMessages function
// MessageValue represents a value of ThreadMessageContent
type MessageValue struct {
	Value string `json:"value"`
}

// AddMessageToThread adds user message to a thread.
// The thread is specified by the threadID argument
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

// GetThreadMessages returns a list of all messages that are stored in a thread as a ThreadMessages struct
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

// LatestAssistantResponse does almost same job as GetThreadMessages except that it extracts last assistant response from the thread messages list.
// See https://github.com/ilborsch/openai-go/blob/main/examples/cli-chat-bot-assistant.go for an example of how to use it correctly
func (m Messages) LatestAssistantResponse(threadID string) (string, error) {
	chatStory, err := m.GetThreadMessages(threadID)
	if err != nil {
		return "", err
	}
	if len(chatStory.Data) == 0 {
		return "", fmt.Errorf("thread %s has no messages", threadID)
	}
	for _, message := range chatStory.Data {
		if message.Role == RoleAssistant {
			return message.Content[0].Text.Value, nil
		}
	}
	return "", fmt.Errorf("thread %s has no assistant responses", threadID)
}
