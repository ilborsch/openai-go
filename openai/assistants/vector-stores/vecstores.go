package vecstores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VectorStoreClient interface {
	Create(storeName string) (string, error)
	Delete(storeID string) error
	AddFile(storeID, fileID string) error
	GetFiles(storeID string) (GetVectorStoreFilesResponse, error)
	DeleteFile(storeID, fileID string) error
}

// VectorStores represents OpenAI API vector store domain
type VectorStores struct {
	APIKey string
}

// CreateVectorStoreResponse is used to unmarshal OpenAI API response in Create function
type CreateVectorStoreResponse struct {
	ID string `json:"id"`
}

// File is used to unmarshal OpenAI API response in GetFiles function
type File struct {
	FileID string `json:"id"`
}

// GetVectorStoreFilesResponse is used to unmarshal OpenAI API response in GetFiles function
type GetVectorStoreFilesResponse struct {
	Files []File `json:"data"`
}

// Create creates a vector store for files that later can be attached to an assistant.
// It is a new feature of assistants v2 API so I sincerely recommend to jump through this docs:
// https://platform.openai.com/docs/api-reference/vector-stores/object
func (v VectorStores) Create(storeName string) (string, error) {
	URL := "https://api.openai.com/v1/vector_stores"
	payload := fmt.Sprintf(`{"name": "%s"}`, storeName)
	requestBody := bytes.NewBuffer([]byte(payload))

	req, err := http.NewRequest("POST", URL, requestBody)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+v.APIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error creating vector store: %v %s ", resp.StatusCode, string(responseBody))
	}

	var response CreateVectorStoreResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	return response.ID, nil
}

// Delete deletes the Vector Store object specified by `storeID` from OpenAI platform.
func (v VectorStores) Delete(storeID string) error {
	URL := fmt.Sprintf("https://api.openai.com/v1/vector_stores/%s", storeID)
	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+v.APIKey)
	req.Header.Add("Content-Type", "application/json")
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
		return fmt.Errorf("failed to delete vector store: %v %s", resp.Status, string(responseBody))
	}
	return nil
}

// AddFile adds a file with `fileID` to the Vector Store object specified by `storeID` from OpenAI platform.
func (v VectorStores) AddFile(storeID, fileID string) error {
	URL := fmt.Sprintf("https://api.openai.com/v1/vector_stores/%s/files", storeID)
	jsonData := fmt.Sprintf(`{"file_id": "%s"}`, fileID)
	reqBody := bytes.NewBufferString(jsonData)

	req, err := http.NewRequest("POST", URL, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+v.APIKey)
	req.Header.Add("Content-Type", "application/json")
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
		return fmt.Errorf("error adding file: %v %s ", resp.StatusCode, string(responseBody))
	}
	return nil
}

// GetFiles returns a list of Files stored in the Vector Store object specified by `storeID` as a struct GetVectorStoreFilesResponse
func (v VectorStores) GetFiles(storeID string) (GetVectorStoreFilesResponse, error) {
	URL := fmt.Sprintf("https://api.openai.com/v1/vector_stores/%s/files", storeID)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return GetVectorStoreFilesResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+v.APIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OpenAI-Beta", "assistants=v2")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetVectorStoreFilesResponse{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetVectorStoreFilesResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GetVectorStoreFilesResponse{}, fmt.Errorf("error retrieving files: %v %s", resp.StatusCode, string(responseBody))
	}

	var filesResponse GetVectorStoreFilesResponse
	if err = json.Unmarshal(responseBody, &filesResponse); err != nil {
		return GetVectorStoreFilesResponse{}, err
	}

	return filesResponse, nil
}

// DeleteFile is used to delete a file from a Vector Store object specified by `storeID`
func (v VectorStores) DeleteFile(storeID, fileID string) error {
	URL := fmt.Sprintf("https://api.openai.com/v1/vector_stores/%s/files/%s", storeID, fileID)

	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+v.APIKey)
	req.Header.Add("Content-Type", "application/json")
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
		return fmt.Errorf("error removing file: %v %s ", resp.StatusCode, string(responseBody))
	}
	return nil
}
