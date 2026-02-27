package domain

const (
	RoleUser  = "user"
	RoleModel = "model"
)

type Message struct {
	Role    string
	Content string
}

func NewUserMessage(content string) Message {
	return Message{Role: RoleUser, Content: content}
}

func NewModelMessage(content string) Message {
	return Message{Role: RoleModel, Content: content}
}
