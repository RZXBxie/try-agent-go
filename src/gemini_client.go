package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type GeminiRequest struct {
	SystemInstruction struct {
		Parts []Part `json:"parts"`
	} `json:"system_instruction"`
	Contents []Content `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []Part `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type GeminiClient struct {
	apiKey string
	model  string
}

func NewGeminiClient(apiKey, model string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		model:  model,
	}
}

func (c *GeminiClient) SendMessage(text string, systemInstruction string) (string, error) {
	body := GeminiRequest{
		Contents: []Content{
			{Role: "user", Parts: []Part{{Text: text}}},
		},
	}
	body.SystemInstruction.Parts = []Part{{Text: systemInstruction}}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.model, c.apiKey,
	)

	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error %d: %s", resp.StatusCode, string(respBytes))
	}

	var data GeminiResponse
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data.Candidates[0].Content.Parts[0].Text, nil
}
