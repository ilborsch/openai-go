package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)


type FileClient interface {
	UploadFile(filename string, fileData []byte) (string, error)
	DeleteFile(fileID string) error
}

// Files represents OpenAI API files domain
type Files struct {
	APIKey string
}

// UploadFileResponse is used to unmarshal OpenAI API response in the UploadFile function
type UploadFileResponse struct {
	ID string `json:"id"`
}

// UploadFile uploads file with `filename` and binary data `fileData` into OpenAI portal
// where can be later used by an assistant
func (f Files) UploadFile(filename string, fileData []byte) (string, error) {
	const URL = "https://api.openai.com/v1/files"

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	_ = writer.WriteField("purpose", "assistants")

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, bytes.NewReader(fileData))
	err = writer.Close()
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", URL, &requestBody)
	request.Header.Set("Authorization", "Bearer "+f.APIKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error uplaoding file: %d %s", resp.StatusCode, string(bodyBytes))
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error deleting file: %d %s ", resp.StatusCode, string(responseBody))
	}
	var response UploadFileResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}
	return response.ID, nil
}

// DeleteFile deletes file object specified by `fileID` from OpenAI portal
func (f Files) DeleteFile(fileID string) error {
	URL := "https://api.openai.com/v1/files/" + fileID
	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+f.APIKey)

	// Create a new HTTP client and send the request
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
		return fmt.Errorf("failed to delete file: %d %s ", resp.StatusCode, string(responseBody))
	}
	return nil
}
