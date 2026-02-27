package provider

import (
	"fmt"

	"try_agent_go/internal/domain"
	"try_agent_go/internal/infra/llm"
)

func New(name string, configFn func(string) string) (domain.Provider, error) {
	switch name {

	case "gemini-2.5-flash":
		apiKey := configFn("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY not set")
		}
		return llm.NewGeminiProvider(apiKey, "gemini-2.5-flash"), nil

	case "deepseek v3.2":
		apiKey := configFn("DEEPSEEK_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("DEEPSEEK_API_KEY not set")
		}
		return llm.NewDeepSeekProvider(apiKey), nil

	case "zhipu":
		apiKey := configFn("ZHIPU_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ZHIPU_API_KEY not set")
		}
		return llm.NewZhiPuProvider(apiKey), nil

	default:
		return nil, fmt.Errorf(
			"unknown provider: %s (supported: gemini-2.5-flash, deepseek v3.2, zhipu)",
			name,
		)
	}
}
