package utils

import (
	"fmt"
	"image/png"
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

func VerifySkin(skinPath string) error {
	inputFile, err := os.Open(skinPath)
	if err != nil {
		return fmt.Errorf("Could not open skin path: %v\n", err)
	}

	img, err := png.Decode(inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Could not decode skin: %v\n", err)
	}

	// it's a valid skin
	if img.Bounds().Dx() <= 64 && img.Bounds().Dy() <= 64 {
		return fmt.Errorf("Skin is to big")
	}

	return nil
}
