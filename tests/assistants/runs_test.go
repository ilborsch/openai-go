package assistants

import (
	"github.com/ilborsch/openai-go/openai/assistants"
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateRun_Happy(t *testing.T) {
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

	runID, err := s.Client.Runs.Create(threadID, assistantID)
	require.NoError(t, err)
	require.NotEmpty(t, runID)
}

func TestCreateRun_InvalidThreadID(t *testing.T) {
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

	threadID := "invalid_thread_id_123"

	runID, err := s.Client.Runs.Create(threadID, assistantID)
	require.Error(t, err)
	require.Empty(t, runID)
}

func TestCreateRun_InvalidAssistantID(t *testing.T) {
	s := suite.New(t)

	assistantID := "invalid_assistant_id_123"

	threadID, err := s.Client.Threads.Create()
	require.NoError(t, err)

	runID, err := s.Client.Runs.Create(threadID, assistantID)
	require.Error(t, err)
	require.Empty(t, runID)
}

func TestGetRun_Happy(t *testing.T) {
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

	runID, err := s.Client.Runs.Create(threadID, assistantID)
	require.NoError(t, err)
	require.NotEmpty(t, runID)

	run, err := s.Client.Runs.GetRun(threadID, runID)
	require.NoError(t, err)
	require.NotEmpty(t, run)
	require.NotEmpty(t, run.Status)
	require.Equal(t, threadID, run.ThreadID)
	require.Equal(t, assistantID, run.AssistantID)
}

func TestGetRun_InvalidID(t *testing.T) {
	s := suite.New(t)

	threadID, err := s.Client.Threads.Create()
	require.NoError(t, err)

	runID := "invalid_run_id_123"
	run, err := s.Client.Runs.GetRun(threadID, runID)
	require.Error(t, err)
	require.Empty(t, run)
}

func TestGetRun_InvalidThreadID(t *testing.T) {
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

	runID, err := s.Client.Runs.Create(threadID, assistantID)
	require.NoError(t, err)
	require.NotEmpty(t, runID)

	invalidThreadID := "invalid_thread_id_123"

	run, err := s.Client.Runs.GetRun(invalidThreadID, runID)
	require.Error(t, err)
	require.Empty(t, run)
}
