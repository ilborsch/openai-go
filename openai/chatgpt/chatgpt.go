package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ilborsch/openai-go/openai/chatgpt/message"
	"io"
	"net/http"
)

const DefaultModel = "gpt-3.5-turbo"

type ChatGPT struct {
	APIKey string
	Model  string
}

type CreateCompletionRequest struct {
	Model    string            `json:"model"`
	Messages []message.Message `json:"messages"`
}

type CompletionChoice struct {
	Index   int `json:"index"`
	Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CreateCompletionResponse struct {
	Choices []CompletionChoice `json:"choices"`
}

func (c ChatGPT) CreateCompletion(chatStory []message.Message) (string, error) {
	const URL = "https://api.openai.com/v1/chat/completions"
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.APIKey,
	}
	completionModel := c.Model
	if c.Model == "" {
		completionModel = DefaultModel
	}

	payload := CreateCompletionRequest{
		Model:    completionModel,
		Messages: chatStory,
	}
	requestBytes, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	client := http.Client{}
	response, err := client.Do(request)
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error creating completion %v %s", response.StatusCode, string(responseBody))
	}

	var completionResponse CreateCompletionResponse
	if err = json.Unmarshal(responseBody, &completionResponse); err != nil {
		return "", nil
	}
	if len(completionResponse.Choices) == 0 {
		return "", fmt.Errorf("no response returned from chatgpt")
	}
	return completionResponse.Choices[0].Message.Content, nil
}
