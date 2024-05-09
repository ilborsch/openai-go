package message

// Message represents a message object in a ChatGPT chatStory
// You need to pass full chat story in your request to `Create a Completion`
// in order to get a response from ChatGPT.
// Message has 3 implementations: UserMessage, SystemMessage and AssistantMessage.
type Message interface {
	Message() string
}

type (
	// UserMessage is used to represent a user message in a chat with ChatGPT
	UserMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	// SystemMessage is used to represent a system instructions to ChatGPT in a chat with it.
	SystemMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	// AssistantMessage represents ChatGPT response
	AssistantMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

// NewUserMessage initializes a new UserMessage object with content
func NewUserMessage(content string) *UserMessage {
	return &UserMessage{
		Role:    "user",
		Content: content,
	}
}

// Message is a getter function to get a message content.
// Is needed to implement the Message interface via Go duck typing.
func (um *UserMessage) Message() string {
	return um.Content
}

// NewSystemMessage initializes a new SystemMessage object with content
func NewSystemMessage(content string) *SystemMessage {
	return &SystemMessage{
		Role:    "system",
		Content: content,
	}
}

// Message is a getter function to get a message content.
// Is needed to implement the Message interface via Go duck typing.
func (sm *SystemMessage) Message() string {
	return sm.Content
}

// NewAssistantMessage initializes a new AssistantMessage object with content
func NewAssistantMessage(content string) *AssistantMessage {
	return &AssistantMessage{
		Role:    "assistant",
		Content: content,
	}
}

// Message is a getter function to get a message content.
// Is needed to implement the Message interface via Go duck typing.
func (am *AssistantMessage) Message() string {
	return am.Content
}
