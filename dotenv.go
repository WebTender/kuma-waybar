package main

import (
	"os"
	"strings"
)

var DEFAULT_ENV_FILES = []string{".env", "~/.kuma-waybar.env"}

// readEnv reads a .env file and returns a map of the key-value pairs
//
// filePath: The path to the .env file. Use "" to use the default .env file
func readEnv(filePath string) (map[string]string, error) {
	if filePath == "" {
		for _, file := range DEFAULT_ENV_FILES {
			if (strings.HasPrefix(file, "~")) {
				file = strings.Replace(file, "~", os.Getenv("HOME"), 1)
			}
			dotenv, err := readEnv(file)
			if err == nil {
				return dotenv, nil
			}
		}

		return map[string]string{}, nil
	}

	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return map[string]string{}, err
	}

	return parseDotEnv(string(fileContents))
}

func parseDotEnv(fileContents string) (map[string]string, error) {
	dotenv := map[string]string{}
	lines := strings.Split(fileContents, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		if line[0] == '#' {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}

		dotenv[key] = value
	}

	return dotenv, nil
}
