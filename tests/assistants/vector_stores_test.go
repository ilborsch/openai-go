package assistants

import (
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestCreateVectorStore_Happy(t *testing.T) {
	s := suite.New(t)

	vsID, err := s.Client.Assistants.VectorStores.Create("test store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)
}

func TestCreateVectorStoreNoName_Happy(t *testing.T) {
	s := suite.New(t)

	vsID, err := s.Client.Assistants.VectorStores.Create("")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)
}

func TestDeleteVectorStore_Happy(t *testing.T) {
	s := suite.New(t)

	vsID, err := s.Client.Assistants.VectorStores.Create("test store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	err = s.Client.Assistants.VectorStores.Delete(vsID)
	require.NoError(t, err)
}

func TestDeleteVectorStore_InvalidID(t *testing.T) {
	s := suite.New(t)

	vsID := "invalid_store_id_123"
	err := s.Client.Assistants.VectorStores.Delete(vsID)
	require.Error(t, err)
}

func TestAddFileToVectorStorage_Happy(t *testing.T) {
	s := suite.New(t)

	filename := "../test-data/test_file.txt"
	file, err := os.Open(filename)
	defer file.Close()
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, fileData)

	fileID, err := s.Client.Files.UploadFile(filename, fileData)
	require.NoError(t, err)
	require.NotEmpty(t, fileID)

	vsID, err := s.Client.Assistants.VectorStores.Create("testing storage")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	err = s.Client.Assistants.VectorStores.AddFile(vsID, fileID)
	require.NoError(t, err)
}

func TestAddFileToVectorStorage_InvalidFileID(t *testing.T) {
	s := suite.New(t)

	fileID := "invalid_file_id_123"

	vsID, err := s.Client.Assistants.VectorStores.Create("testing storage")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	err = s.Client.Assistants.VectorStores.AddFile(vsID, fileID)
	require.Error(t, err)
}

func TestGetFilesFromVectorStore_Happy(t *testing.T) {
	s := suite.New(t)

	filename := "../test-data/test_file.txt"
	file, err := os.Open(filename)
	defer file.Close()
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, fileData)

	fileID, err := s.Client.Files.UploadFile(filename, fileData)
	require.NoError(t, err)
	require.NotEmpty(t, fileID)

	vsID, err := s.Client.Assistants.VectorStores.Create("testing storage")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	err = s.Client.Assistants.VectorStores.AddFile(vsID, fileID)
	require.NoError(t, err)

	filesResponse, err := s.Client.Assistants.VectorStores.GetFiles(vsID)
	require.NoError(t, err)
	require.NotEmpty(t, filesResponse)
	require.NotEmpty(t, filesResponse.Files)

	for _, file := range filesResponse.Files {
		if file.FileID == fileID {
			return
		}
	}
	t.Fatal("function GetFiles did not return all required files")
}

func TestGetFilesFromVectorStore_InvalidStoreID(t *testing.T) {
	s := suite.New(t)

	vsID := "invalid_store_id_123"

	filesResponse, err := s.Client.Assistants.VectorStores.GetFiles(vsID)
	require.Error(t, err)
	require.Empty(t, filesResponse)
}

func TestDeleteFile_Happy(t *testing.T) {
	s := suite.New(t)

	filename := "../test-data/test_file.txt"
	file, err := os.Open(filename)
	defer file.Close()
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, fileData)

	fileID, err := s.Client.Files.UploadFile(filename, fileData)
	require.NoError(t, err)
	require.NotEmpty(t, fileID)

	vsID, err := s.Client.Assistants.VectorStores.Create("testing storage")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	err = s.Client.Assistants.VectorStores.AddFile(vsID, fileID)
	require.NoError(t, err)

	err = s.Client.Assistants.VectorStores.DeleteFile(vsID, fileID)
	require.NoError(t, err)
}

func TestDeleteFile_InvalidStoreID(t *testing.T) {
	s := suite.New(t)

	filename := "../test-data/test_file.txt"
	file, err := os.Open(filename)
	defer file.Close()
	require.NoError(t, err)

	fileData, err := io.ReadAll(file)
	require.NoError(t, err)
	assert.NotEmpty(t, fileData)

	fileID, err := s.Client.Files.UploadFile(filename, fileData)
	require.NoError(t, err)
	require.NotEmpty(t, fileID)

	vsID := "invalid_store_id_123"

	err = s.Client.Assistants.VectorStores.DeleteFile(vsID, fileID)
	require.Error(t, err)
}

func TestDeleteFile_InvalidFileID(t *testing.T) {
	s := suite.New(t)

	vsID, err := s.Client.Assistants.VectorStores.Create("testing storage")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	fileID := "invalid_file_id_123"

	err = s.Client.Assistants.VectorStores.DeleteFile(vsID, fileID)
	require.Error(t, err)
}
