package typespeed

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Prompts        []Prompt `yaml:"prompts"`
	ActivePromptID int
	SeenIDs        map[int]int

	// how many prompts are active for the specific
	// type requested. Used for determining when game should
	// end if all prompts are used
	ActivePromptsLen int

	// The type of game mode
	PromptType string

	// Color for printing
	PromptTypeColor string
}

type Prompt struct {
	ID         int    `yaml:"id"`
	Text       string `yaml:"text"`
	Difficulty string `yaml:"difficulty"`
	Type       string `yaml:"type"`
}

func parseYAML(filePath string) (*Config, error) {
	path := "library.yaml"
	if filePath != "" {
		path = filePath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
