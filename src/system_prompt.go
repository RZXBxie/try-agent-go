package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getPromptsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		// fallback: relative to working directory
		return filepath.Join("src", "prompts")
	}
	return filepath.Join(filepath.Dir(filename), "prompts")
}

// GetSystemPrompt reads a Markdown file from src/prompts/ directory.
// Switch styles by changing the name parameter.
func GetSystemPrompt(name string) (string, error) {
	if name == "" {
		name = "personal-assistant"
	}
	filePath := filepath.Join(getPromptsDir(), name+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %q: %w", filePath, err)
	}
	return strings.TrimSpace(string(data)), nil
}
