package domain

type Conversation struct {
	history []Message
}

func NewConversation() *Conversation {
	return &Conversation{}
}

func (c *Conversation) AddUserMessage(content string) {
	c.history = append(c.history, NewUserMessage(content))
}

func (c *Conversation) AddModelMessage(content string) {
	c.history = append(c.history, NewModelMessage(content))
}

func (c *Conversation) UndoLast() {
	if len(c.history) > 0 {
		c.history = c.history[:len(c.history)-1]
	}
}

func (c *Conversation) History() []Message {
	out := make([]Message, len(c.history))
	copy(out, c.history)
	return out
}

func (c *Conversation) Reset() {
	c.history = nil
}
