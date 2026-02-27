package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(promptsDir, name string) (string, error) {
	if name == "" {
		name = "personal-assistant"
	}
	filePath := filepath.Join(promptsDir, name+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read prompt file %q: %w", filePath, err)
	}
	return strings.TrimSpace(string(data)), nil
}
