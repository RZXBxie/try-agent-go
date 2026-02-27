package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"try_agent_go/internal/domain"
)

type ZhiPuProvider struct {
	apiKey string
	model  string
}

func NewZhiPuProvider(apiKey string) *ZhiPuProvider {
	return &ZhiPuProvider{
		apiKey: apiKey,
		model:  "glm-4-flash",
	}
}

func (p *ZhiPuProvider) Name() string {
	return "zhipu"
}

func (p *ZhiPuProvider) SendMessage(history []domain.Message, systemInstruction string) (string, error) {
	messages := []chatMessage{
		{Role: "system", Content: systemInstruction},
	}
	for _, m := range history {
		role := m.Role
		if role == domain.RoleModel {
			role = "assistant"
		}
		messages = append(messages, chatMessage{Role: role, Content: m.Content})
	}

	reqBody := chatCompletionRequest{
		Model:    p.model,
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ZhiPu API error %d: %s", resp.StatusCode, string(respBytes))
	}

	var data chatCompletionResponse
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(data.Choices) == 0 {
		return "", fmt.Errorf("ZhiPu API returned no choices")
	}

	return data.Choices[0].Message.Content, nil
}
