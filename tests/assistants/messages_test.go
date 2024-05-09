package assistants

import (
	"github.com/ilborsch/openai-go/openai/assistants"
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAddMessageToThread_Happy(t *testing.T) {
	s := suite.New(t)
	threadID, err := s.Client.Threads.Create()
	require.NoError(t, err)
	require.NotEmpty(t, threadID)

	err = s.Client.Messages.AddMessageToThread(threadID, "random_message")
	require.NoError(t, err)
}

func TestGetThreadMessages_Happy(t *testing.T) {
	s := suite.New(t)

	threadID, err := s.Client.Threads.Create()
	require.NoError(t, err)
	require.NotEmpty(t, threadID)

	messageContent := "random_message"
	err = s.Client.Messages.AddMessageToThread(threadID, messageContent)
	require.NoError(t, err)

	messages, err := s.Client.Messages.GetThreadMessages(threadID)
	require.NoError(t, err)
	require.NotEmpty(t, messages)
	require.NotEmpty(t, messages.Data)
	require.Equal(t, messageContent, messages.Data[0].Content[0].Text.Value)
}

func TestGetThreadMessages_InvalidThreadID(t *testing.T) {
	s := suite.New(t)

	threadID := "invalid_thread_id_123"

	messages, err := s.Client.Messages.GetThreadMessages(threadID)
	require.Error(t, err)
	require.Empty(t, messages)
	require.Empty(t, messages.Data)
}

func TestLatestAssistantResponse_Happy(t *testing.T) {
	s := suite.New(t)

	name := "Testing Assistant"
	instructions := "Test instructions"
	tools := []assistants.Tool{
		{
			Type: assistants.ToolFileSearch,
		},
	}

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	assistantID, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, assistantID)

	threadID, err := s.Client.Threads.Create()
	require.NoError(t, err)
	require.NotEmpty(t, threadID)

	messageContent := "Hi, Can you help me?"
	err = s.Client.Messages.AddMessageToThread(threadID, messageContent)
	require.NoError(t, err)

	runID, err := s.Client.Runs.Create(threadID, assistantID)
	require.NoError(t, err)
	require.NotEmpty(t, runID)

	// simulate pooling
	time.Sleep(7 * time.Second)

	response, err := s.Client.Messages.LatestAssistantResponse(threadID)
	require.NoError(t, err)
	require.NotEmpty(t, response)
}

func TestLatestAssistantResponse_InvalidThreadID(t *testing.T) {
	s := suite.New(t)

	threadID := "invalid_thread_id_123"

	response, err := s.Client.Messages.LatestAssistantResponse(threadID)
	require.Error(t, err)
	require.Empty(t, response)
}
