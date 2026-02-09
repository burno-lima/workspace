package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"workspace/pkg/models"
	str "workspace/pkg/strings"
)

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configDir := filepath.Join(home, ".config", "workspace")
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "workspace.config")
}

func Load() (models.Config, error) {
	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return models.Config{ProjectEditors: make(map[string]string)}, err
	}

	config := models.Config{ProjectEditors: make(map[string]string)}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "default" {
				config.DefaultEditor = value
			} else {
				config.ProjectEditors[key] = value
			}
		}
	}
	return config, nil
}

func Save(config models.Config) error {
	configPath := GetConfigPath()
	var lines []string
	lines = append(lines, str.ConfigHeader.String())
	lines = append(lines, str.ConfigFormat.String())
	lines = append(lines, "")
	if config.DefaultEditor != "" {
		lines = append(lines, fmt.Sprintf(str.ConfigDefaultPrefix.String(), config.DefaultEditor))
		lines = append(lines, "")
	}
	for project, editor := range config.ProjectEditors {
		lines = append(lines, fmt.Sprintf("%s=%s", project, editor))
	}
	data := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(configPath, []byte(data), 0644)
}

func Reset() error {
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		return os.Remove(configPath)
	}
	return nil
}
