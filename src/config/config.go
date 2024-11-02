package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

var llmFormats = [2]string{"openai", "ollama"}

type Configuration struct {
	DatabaseURL string `json:"database_url"`
	StoragePath string `json:"storage_path"`
	LLMConfig   struct {
		Url          string `json:"url"`
		Format       string `json:"format"`
		Model        string `json:"model"`
		Token        string `json:"token"`
		SystemPrompt string `json:"system_prompt"`
	} `json:"llm_config"`
}

func LoadConfig(configFilePath string) (*Configuration, error) {
	// Open the config file
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, errors.New("unable to open config file: " + err.Error())
	}
	defer file.Close()

	// Read the config file
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("unable to read config file: " + err.Error())
	}

	// Unmarshal the config file
	var config Configuration
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, errors.New("unable to unmarshal config file: " + err.Error())
	}

	return &config, config.verify()
}

func (c *Configuration) verify() error {
	switch {
	case c.StoragePath == "":
		return errors.New("storage_path cannot be empty")
	case c.DatabaseURL == "":
		return errors.New("database_url cannot be empty")
	case c.LLMConfig.Url == "":
		return errors.New("llm_config.url cannot be empty")
	case c.LLMConfig.Format == "":
		return errors.New("llm_config.format cannot be empty")
	case func() bool {
		for _, format := range llmFormats {
			if c.LLMConfig.Format == format {
				return false
			}
		}
		return true
	}():
		return errors.New("llm_config.format must be one of: " + strings.Join(llmFormats[:], ", "))
	case c.LLMConfig.Model == "":
		return errors.New("llm_config.model cannot be empty")
	// case c.LLMConfig.Token == "":
	// 	return errors.New("llm_config.token cannot be empty")
	default:
		return nil
	}
}
