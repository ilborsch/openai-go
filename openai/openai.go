package openai

import (
	"github.com/ilborsch/openai-go/openai/assistants"
	"github.com/ilborsch/openai-go/openai/assistants/messages"
	"github.com/ilborsch/openai-go/openai/assistants/runs"
	"github.com/ilborsch/openai-go/openai/assistants/threads"
	vecstores "github.com/ilborsch/openai-go/openai/assistants/vector-stores"
	"github.com/ilborsch/openai-go/openai/chatgpt"
	"github.com/ilborsch/openai-go/openai/files"
)

// OpenAI is a main client and centre of user interaction with the openai-go library.
// Combines all subdomains which may be accessed by full composition relation
// ( f.e. `OpenAI.Assistants.VectorStores.Create(...)` ),
// or alternatively by shortened syntax ( f.e. `OpenAI.VectorStores.Create(...)` ).
// Both examples do absolutely the same, it's just a syntax sugar from Go language.
type OpenAI struct {
	apiKey string
	chatgpt.ChatGPTClient
	files.FileClient
	assistants.AssistantClient
}

// New initializes a new OpenAI instance and returns it
func New(apiKey string) *OpenAI {
	if apiKey == "" {
		panic("api key cannot be empty")
	}
	return &OpenAI{
		apiKey: apiKey,
		ChatGPTClient: chatgpt.ChatGPT{
			APIKey: apiKey,
			Model:  chatgpt.DefaultModel,
		},
		FileClient: files.Files{
			APIKey: apiKey,
		},
		AssistantClient: assistants.Assistants{
			APIKey: apiKey,
			VectorStoreClient: vecstores.VectorStores{
				APIKey: apiKey,
			},
			MessageClient: messages.Messages{
				APIKey: apiKey,
			},
			RunClient: runs.Runs{
				APIKey: apiKey,
			},
			ThreadClient: threads.Threads{
				APIKey: apiKey,
			},
		},
	}
}
