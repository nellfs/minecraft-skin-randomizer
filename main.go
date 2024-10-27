package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// TODO: organize this code in more files PLEASEEEE

const (
	SkinsFolder      = "skins"
	BaseSkinFolder   = "base"
	HeadSkinFolder   = "head"
	Layer0SkinFolder = "layer_0"
	Layer1SkinFolder = "layer_1"

	Version           = "0.0.1"
	ConfigurationFile = "config.json"
)

type ConfigFile struct {
	ConfigVersion    string `json:"version"`
	EditableSkin     string `json:"edit_skin"`
	RandomizerFolder string `json:"randomizer_folder"`
	// todo: implement that in the future
	BaseSkinGeneration ConfigGeneration `json:"base_generation"`
}

// todo: implement that:
//
// each body part should have an configGeneration, so users can config behaviors
// for each skin body part
type ConfigGeneration struct {
	Recolor bool `json:"recolor"`
}

func SelectPath(inputPath string) (fullPath string, err error) {
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		return inputPath, fmt.Errorf("Error gettint the absolute path: %w\n", err)
	}

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return inputPath, fmt.Errorf("This file path does not exist: %w\n", err)
	}

	return absolutePath, nil
}

type SkinManager struct {
	RandomizerFolderPath string
	SkinPath             string
	BasePath             SkinPart
	HeadPath             SkinPart
}

type SkinPart struct {
	RootPath   string
	Layer0Path string
	Layer1Path string
}

func (sm *SkinManager) setupSkinParts() error {
	necessaryParts := []SkinPart{sm.BasePath, sm.HeadPath}

	for _, part := range necessaryParts {
		part.Layer0Path = fmt.Sprintf("%s/%s", part.RootPath, Layer0SkinFolder)
		part.Layer1Path = fmt.Sprintf("%s/%s", part.RootPath, Layer1SkinFolder)

		if _, err := os.Stat(part.RootPath); os.IsNotExist(err) {
			err := os.Mkdir(part.RootPath, 0755)
			if err != nil {
				return fmt.Errorf("Error creating folder: %w\n", err)
			}
		}

		if _, err := os.Stat(part.Layer0Path); os.IsNotExist(err) {
			err := os.Mkdir(part.Layer0Path, 0755)
			if err != nil {
				return fmt.Errorf("Error creating folder: %w\n", err)
			}
		}

		if _, err := os.Stat(part.Layer1Path); os.IsNotExist(err) {
			err := os.Mkdir(part.Layer1Path, 0755)
			if err != nil {
				return fmt.Errorf("Error creating folder: %w\n", err)
			}
		}
	}

	return nil
}

func (sm *SkinManager) GenerateSkin() {

}

func NewSkinManager() (skinManager SkinManager, err error) {

	// the most important thing is the config file,
	// add flag to config file json
	// - stop asking for configuration via user input?

	// if config file does not exist users can generate one
	// just use the flags and it will generate the config filepath
	//
	// flags:
	// config, skin-dir, randomize-dir => skinDir and randomizeDir you can create a config folder
	// if the config is set you can overwrite it with skindir and randomizerDir, this should
	// save the new configuration in the config file
	//
	// recolor generation configs are experimental, so I will not add flags for it for now

	skinDir := flag.String("skin-dir", "", "Path to the skin directory")
	randomizeDir := flag.String("randomize-dir", "", "Path to the randomize directory")
	flag.Parse()

	var randomizerFolderPath string
	var minecraftSkinPath string

	if randomizeDir == nil {
		fmt.Println("Please paste the path of your randomizer skin folder:")
		// - remove those scans??
		fmt.Scan(&randomizerFolderPath)
		randomizerFolderPath, err = SelectPath(randomizerFolderPath)
		if err != nil {
			return skinManager, err
		}
	}

	if skinDir == nil {
		fmt.Println("Please paste the path of the skin on the minecraft folder to be customized:")
		fmt.Scan(&minecraftSkinPath)
		minecraftSkinPath, err = SelectPath(minecraftSkinPath)
		if err != nil {
			return skinManager, err
		}
	}

	skinManager = SkinManager{
		SkinPath:             minecraftSkinPath,
		RandomizerFolderPath: randomizerFolderPath,
	}

	skinManager.BasePath.RootPath = fmt.Sprintf("%s/%s/%s", randomizerFolderPath, SkinsFolder, BaseSkinFolder)
	skinManager.HeadPath.RootPath = fmt.Sprintf("%s/%s/%s", randomizerFolderPath, SkinsFolder, HeadSkinFolder)

	err = skinManager.setupSkinParts()
	if err != nil {
		fmt.Println("error", err)
	}

	configFile := ConfigFile{
		ConfigVersion:    Version,
		EditableSkin:     minecraftSkinPath,
		RandomizerFolder: randomizerFolderPath,
	}

	jsonData, err := json.MarshalIndent(configFile, "", "  ")
	if err != nil {
		return skinManager, err
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", randomizerFolderPath, ConfigurationFile))
	if err != nil {
		return skinManager, fmt.Errorf("Error creating file: %w\n", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return skinManager, fmt.Errorf("Error writing to file: %w\n", err)
	}

	return skinManager, nil
}

func main() {
	NewSkinManager()
}
