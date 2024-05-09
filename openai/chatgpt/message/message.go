package message

type Message interface {
	Message() string
}

type (
	UserMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	SystemMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	AssistantMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)

func NewUserMessage(content string) *UserMessage {
	return &UserMessage{
		Role:    "user",
		Content: content,
	}
}

func (um *UserMessage) Message() string {
	return um.Content
}

func NewSystemMessage(content string) *SystemMessage {
	return &SystemMessage{
		Role:    "system",
		Content: content,
	}
}

func (sm *SystemMessage) Message() string {
	return sm.Content
}

func NewAssistantMessage(content string) *AssistantMessage {
	return &AssistantMessage{
		Role:    "assistant",
		Content: content,
	}
}

func (am *AssistantMessage) Message() string {
	return am.Content
}
