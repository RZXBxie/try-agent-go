package main

import "try_agent_go/src/client"

// ---------------------------------------------------------------------------
// Chat — 维护对话历史，实现多轮记忆
// ---------------------------------------------------------------------------

// Chat wraps a Provider with conversation history so the LLM can "remember"
// previous messages. Each call to Send appends the user message and the model
// reply to the history, and sends the full history on every API call.
type Chat struct {
	provider     Provider
	systemPrompt string
	history      []client.Message
}

// NewChat creates a Chat bound to the given provider and system prompt.
func NewChat(provider Provider, systemPrompt string) *Chat {
	return &Chat{
		provider:     provider,
		systemPrompt: systemPrompt,
	}
}

// Send sends a user message, appending it to history, calling the provider
// with the full conversation, and appending the model reply to history.
func (c *Chat) Send(text string) (string, error) {
	// 追加用户消息
	c.history = append(c.history, client.Message{Role: "user", Content: text})

	// 将完整历史发给 LLM
	reply, err := c.provider.SendMessage(c.history, c.systemPrompt)
	if err != nil {
		// 调用失败时移除刚追加的用户消息，保持历史一致
		c.history = c.history[:len(c.history)-1]
		return "", err
	}

	// 追加模型回复
	c.history = append(c.history, client.Message{Role: "model", Content: reply})

	return reply, nil
}
