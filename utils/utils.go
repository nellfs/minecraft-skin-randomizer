package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FormatPath(inputPath string) (fullPath string, err error) {
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		return inputPath, fmt.Errorf("Error gettint the absolute path: %w\n", err)
	}

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return inputPath, fmt.Errorf("This file path does not exist: %w\n", err)
	}

	return absolutePath, nil
}
