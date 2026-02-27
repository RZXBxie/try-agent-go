package main

import (
	"fmt"
	"os"
	"try_agent_go/src/client"
)

// ---------------------------------------------------------------------------
// Provider — 统一的 LLM 后端接口
// ---------------------------------------------------------------------------

// Provider defines a common interface for all LLM backends.
type Provider interface {
	SendMessage(history []client.Message, systemInstruction string) (string, error)
	Name() string
}

// ---------------------------------------------------------------------------
// Factory — 根据名称创建对应的 Provider
// ---------------------------------------------------------------------------

// CreateProvider creates a Provider by name.
func CreateProvider(name string) (Provider, error) {
	switch name {

	// --- Gemini ---
	case "gemini-2.5-flash":
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY not set")
		}
		return client.NewGeminiClient(apiKey, "gemini-2.5-flash"), nil

	// --- DeepSeek ---
	case "deepseek v3.2":
		apiKey := os.Getenv("DEEPSEEK_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("DEEPSEEK_API_KEY not set")
		}
		return client.NewDeepSeekClient(apiKey), nil

	// --- ZhiPu (GLM) ---
	case "zhipu":
		apiKey := os.Getenv("ZHIPU_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ZHIPU_API_KEY not set")
		}
		return client.NewZhiPuClient(apiKey), nil

	default:
		return nil, fmt.Errorf(
			"unknown provider: %s (supported: gemini-2.5-flash, deepseek v3.2, zhipu)",
			name,
		)
	}
}
