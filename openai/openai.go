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
	chatgpt.ChatGPT
	files.Files
	assistants.Assistants
}

// New initializes a new OpenAI instance and returns it
func New(apiKey string) *OpenAI {
	if apiKey == "" {
		panic("api key cannot be empty")
	}
	return &OpenAI{
		apiKey: apiKey,
		ChatGPT: chatgpt.ChatGPT{
			APIKey: apiKey,
			Model:  chatgpt.DefaultModel,
		},
		Files: files.Files{
			APIKey: apiKey,
		},
		Assistants: assistants.Assistants{
			APIKey: apiKey,
			VectorStores: vecstores.VectorStores{
				APIKey: apiKey,
			},
			Messages: messages.Messages{
				APIKey: apiKey,
			},
			Runs: runs.Runs{
				APIKey: apiKey,
			},
			Threads: threads.Threads{
				APIKey: apiKey,
			},
		},
	}
}
