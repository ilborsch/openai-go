package assistants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ilborsch/openai-go/openai/assistants/messages"
	"github.com/ilborsch/openai-go/openai/assistants/runs"
	"github.com/ilborsch/openai-go/openai/assistants/threads"
	vecstores "github.com/ilborsch/openai-go/openai/assistants/vector-stores"
	"io"
	"net/http"
)

const (
	DefaultModel        = "gpt-3.5-turbo"
	ToolFileSearch      = "file_search"
	ToolCodeInterpreter = "code_interpreter"
)

type Assistants struct {
	APIKey string
	vecstores.VectorStores
	messages.Messages
	runs.Runs
	threads.Threads
}

type CreateAssistantRequest struct {
	Instructions  string `json:"instructions"`
	Name          string `json:"name"`
	Tools         []Tool `json:"tools"`
	Model         string `json:"model"`
	ToolResources `json:"tool_resources"`
}

type ToolResources struct {
	FileSearch VectorStoreIDs `json:"file_search"`
}

type VectorStoreIDs struct {
	IDs []string `json:"vector_store_ids"`
}

type Tool struct {
	Type string `json:"type"`
}

type CreateAssistantResponse struct {
	AssistantID string `json:"id"`
}

type GetAssistantResponse struct {
	Name         string  `json:"name"`
	Model        string  `json:"model"`
	Instructions string  `json:"instructions"`
	Tools        []Tool  `json:"tools"`
	Temperature  float32 `json:"temperature"`
}

func (a Assistants) Create(name, instructions, vectorStoreID string, tools []Tool) (string, error) {
	const URL = "https://api.openai.com/v1/assistants"

	assistantConfig := CreateAssistantRequest{
		Name:         name,
		Instructions: instructions,
		Model:        DefaultModel,
		Tools:        tools,
		ToolResources: ToolResources{
			FileSearch: VectorStoreIDs{
				IDs: []string{vectorStoreID},
			},
		},
	}
	if len(tools) == 0 { // Default set of tools
		assistantConfig.Tools = []Tool{
			{Type: ToolFileSearch},
		}
	}

	requestBody, err := json.Marshal(&assistantConfig)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+a.APIKey)
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
		return string(responseBody), fmt.Errorf("reponse status code: %v %s ", resp.StatusCode, string(responseBody))
	}

	var response CreateAssistantResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	return response.AssistantID, nil
}

func (a Assistants) GetAssistant(ID string) (GetAssistantResponse, error) {
	URL := "https://api.openai.com/v1/assistants/" + ID
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return GetAssistantResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+a.APIKey)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return GetAssistantResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetAssistantResponse{}, fmt.Errorf("response status code: %v", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetAssistantResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GetAssistantResponse{}, fmt.Errorf("error creating assistant: %v %s", resp.StatusCode, string(responseBody))
	}

	var response GetAssistantResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return GetAssistantResponse{}, err
	}
	return response, nil
}

func (a Assistants) Modify(assistantID, newInstructions, newModel string, newTemperature float32) error {
	URL := "https://api.openai.com/v1/assistants/" + assistantID

	requestMap := make(map[string]any)
	if newInstructions != "" {
		requestMap["instructions"] = newInstructions
	}
	if newModel != "" {
		requestMap["model"] = newModel
	}
	if newTemperature != 0.0 {
		requestMap["temperature"] = newTemperature
	}

	requestBody, err := json.Marshal(requestMap)
	request, err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+a.APIKey)
	request.Header.Set("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error modifying assistant: %v %s", resp.StatusCode, string(responseBody))
	}
	return nil
}

func (a Assistants) Delete(ID string) error {
	req, err := http.NewRequest("DELETE", "https://api.openai.com/v1/assistants/"+ID, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+a.APIKey)
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting assistant: %v %s", resp.StatusCode, string(responseBody))
	}
	return nil
}
