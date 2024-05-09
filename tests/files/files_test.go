package files

import (
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestFileUpload_Happy(t *testing.T) {
	s := suite.New(t)
	filename := "../test-data/test_file.txt"
	file, err := os.Open(filename)
	defer file.Close()
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, fileData)

	id, err := s.Client.Files.UploadFile(filename, fileData)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestFileUpload_NonExistingFile(t *testing.T) {
	s := suite.New(t)
	filename := "../test-data/test_file123.txt"
	file, err := os.Open(filename)
	defer file.Close()
	require.Error(t, err)

	fileData, err := io.ReadAll(file)
	require.Error(t, err)
	assert.Empty(t, fileData)

	id, err := s.Client.Files.UploadFile(filename, fileData)
	require.Error(t, err)
	require.Empty(t, id)
}

func TestFileUpload_InvalidFormat(t *testing.T) {
	s := suite.New(t)
	filename := "../test-data/file.zip"
	file, err := os.Open(filename)
	defer file.Close()
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)

	id, err := s.Client.Files.UploadFile(filename, fileData)
	require.Error(t, err)
	require.Empty(t, id)
}

func TestDeleteFile_Happy(t *testing.T) {
	s := suite.New(t)
	filename := "./test-data/test_file.txt"
	file, err := os.Open(filename)
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, fileData)

	fileID, err := s.Client.Files.UploadFile(filename, fileData)
	require.NoError(t, err)
	require.NotEmpty(t, fileID)

	err = s.Client.Files.DeleteFile(fileID)
	require.NoError(t, err)
}

func TestDeleteFile_NonExistingFile(t *testing.T) {
	s := suite.New(t)
	err := s.Client.Files.DeleteFile("random_id_123")
	require.Error(t, err)
}
