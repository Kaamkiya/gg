package main

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Prompts []Prompt `yaml:"prompts"`
	ActivePromptID int
}

type Prompt struct {
	ID         int    `yaml:"id"`
	Text       string `yaml:"text"`
	Difficulty string `yaml:"difficulty"`
}

func parseYAML(filePath string) (*Config, error) {
	path := "library.yaml"
	if filePath != "" {
		path = filePath
	}

	data, err := os.ReadFile(path)
	if err != nil { return nil, err }

	var cfg Config

	if err = yaml.Unmarshal(data, &cfg); err != nil { return nil, err}

	return &cfg, nil

}
