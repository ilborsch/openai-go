package assistants

import (
	"github.com/ilborsch/openai-go/openai/assistants"
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAssistant_Happy(t *testing.T) {
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
}

func TestCreateAssistantNoTools_Happy(t *testing.T) {
	s := suite.New(t)
	name := "Testing Assistant"
	instructions := "Test instructions"
	tools := make([]assistants.Tool, 0)

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	id, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestCreateAssistantNoName_Happy(t *testing.T) {
	s := suite.New(t)
	name := ""
	instructions := "Test Instructions"
	tools := make([]assistants.Tool, 0)

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	id, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func TestGetAssistant_Happy(t *testing.T) {
	s := suite.New(t)
	name := "Testing Assistant"
	instructions := "Test instructions"
	tools := make([]assistants.Tool, 0)

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	id, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	assistant, err := s.Client.Assistants.GetAssistant(id)
	require.NoError(t, err)
	require.NotEmpty(t, assistant)
}

func TestGetAssistant_InvalidID(t *testing.T) {
	s := suite.New(t)
	id := "invalid_assistant_id_123"

	assistant, err := s.Client.Assistants.GetAssistant(id)
	require.Error(t, err)
	require.Empty(t, assistant)
}

func TestModifyAssistant_Happy(t *testing.T) {
	s := suite.New(t)
	name := "Testing Assistant"
	instructions := "Test instructions"
	tools := make([]assistants.Tool, 0)

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	id, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newInstructions := "new instructions"
	newModel := assistants.DefaultModel
	newTemperature := float32(0.8)

	err = s.Client.Assistants.Modify(id, newInstructions, newModel, newTemperature)
	require.NoError(t, err)
}

func TestModifyAssistant_NoUpdate(t *testing.T) {
	s := suite.New(t)
	name := "Testing Assistant"
	instructions := "Test instructions"
	tools := make([]assistants.Tool, 0)

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	id, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	newInstructions := ""
	newModel := ""
	newTemperature := float32(0)

	err = s.Client.Assistants.Modify(id, newInstructions, newModel, newTemperature)
	require.NoError(t, err)
}

func TestDeleteAssistant_Happy(t *testing.T) {
	s := suite.New(t)
	name := "Testing Assistant"
	instructions := "Test instructions"
	tools := make([]assistants.Tool, 0)

	vsID, err := s.Client.Assistants.VectorStores.Create("test_store")
	require.NoError(t, err)
	require.NotEmpty(t, vsID)

	id, err := s.Client.Assistants.Create(name, instructions, vsID, tools)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = s.Client.Assistants.Delete(id)
	require.NoError(t, err)
}

func TestDeleteAssistant_InvalidID(t *testing.T) {
	s := suite.New(t)
	id := "invalid_assistant_id_123"
	err := s.Client.Assistants.Delete(id)
	require.Error(t, err)
}
