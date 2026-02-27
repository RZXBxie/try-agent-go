package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ===========================================================================
// ZhiPuClient — 智谱 GLM 系列
// ===========================================================================
// 复用 deepseek_client.go 中的 ChatMessage / ChatCompletionRequest / ChatCompletionResponse

type ZhiPuClient struct {
	apiKey string
	model  string
}

func NewZhiPuClient(apiKey string) *ZhiPuClient {
	return &ZhiPuClient{
		apiKey: apiKey,
		model:  "glm-4-flash",
	}
}

// ---------------------------------------------------------------------------
// Provider 接口实现
// ---------------------------------------------------------------------------

func (c *ZhiPuClient) Name() string {
	return "zhipu"
}

func (c *ZhiPuClient) SendMessage(history []Message, systemInstruction string) (string, error) {
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
	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", bytes.NewReader(bodyBytes))
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
		return "", fmt.Errorf("ZhiPu API error %d: %s", resp.StatusCode, string(respBytes))
	}

	var data ChatCompletionResponse
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(data.Choices) == 0 {
		return "", fmt.Errorf("ZhiPu API returned no choices")
	}

	return data.Choices[0].Message.Content, nil
}
