package cmd

import "fmt"

// SampleCommand executes the sample action.
// It validates input, processes the data, and returns a result.
func SampleCommand(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input must not be empty")
	}
	return "processed: " + input, nil
}

// SampleConfig holds configuration for the sample command.
type SampleConfig struct {
	Verbose bool   `yaml:"verbose"`
	Output  string `yaml:"output"`
}
