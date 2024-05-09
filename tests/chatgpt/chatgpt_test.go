package chatgpt

import (
	"github.com/ilborsch/openai-go/openai/chatgpt/message"
	"github.com/ilborsch/openai-go/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChatGPTCompletion_Happy(t *testing.T) {
	s := suite.New(t)
	chatStory := []message.Message{
		message.NewSystemMessage("Be very sarcastic"),
		message.NewUserMessage("What is the second name for an apple?"),
	}
	response, err := s.Client.ChatGPT.CreateCompletion(chatStory)
	require.NoError(t, err)
	assert.NotEmpty(t, response)

	chatStory = append(chatStory, message.NewAssistantMessage(response))
	chatStory = append(chatStory, message.NewUserMessage("What is the second name for an orange?"))

	response, err = s.Client.ChatGPT.CreateCompletion(chatStory)
	require.NoError(t, err)
	assert.NotEmpty(t, response)
}

func TestChatGPTCompletion_Error(t *testing.T) {
	s := suite.New(t)
	chatStory := make([]message.Message, 0)
	response, err := s.Client.ChatGPT.CreateCompletion(chatStory)
	require.Error(t, err)
	assert.Empty(t, response)

	chatStory = nil

	response, err = s.Client.ChatGPT.CreateCompletion(chatStory)
	require.Error(t, err)
	assert.Empty(t, response)
}
