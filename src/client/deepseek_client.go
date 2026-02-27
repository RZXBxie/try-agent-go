package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ===========================================================================
// 通用消息类型（所有 Provider 共用）
// ===========================================================================

// Message represents a single message in the conversation history.
type Message struct {
	Role    string // "user" or "model"
	Content string
}

// ===========================================================================
// OpenAI 兼容的共享类型（DeepSeek / ZhiPu 共用）
// ===========================================================================

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// ===========================================================================
// DeepSeekClient
// ===========================================================================

type DeepSeekClient struct {
	apiKey string
	model  string
}

func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey: apiKey,
		model:  "deepseek-chat",
	}
}

// ---------------------------------------------------------------------------
// Provider 接口实现
// ---------------------------------------------------------------------------

func (c *DeepSeekClient) Name() string {
	return "deepseek"
}

func (c *DeepSeekClient) SendMessage(history []Message, systemInstruction string) (string, error) {
	// 1. 构造请求体：system prompt + 完整对话历史
	messages := []ChatMessage{
		{Role: "system", Content: systemInstruction},
	}
	for _, m := range history {
		role := m.Role
		if role == "model" {
			role = "assistant"
		}
		messages = append(messages, ChatMessage{Role: role, Content: m.Content})
	}

	reqBody := ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 2. 发送 HTTP 请求
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// 3. 读取 & 解析响应
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API error %d: %s", resp.StatusCode, string(respBytes))
	}

	var data ChatCompletionResponse
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(data.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek API returned no choices")
	}

	return data.Choices[0].Message.Content, nil
}
