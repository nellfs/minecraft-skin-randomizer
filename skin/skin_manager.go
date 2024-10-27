package skin

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nellfs/minecraft-skin-randomizer/configuration"
	"github.com/nellfs/minecraft-skin-randomizer/utils"
)

const (
	SkinsFolder      = "/skins"
	BaseSkinFolder   = "/base"
	HeadSkinFolder   = "/head"
	Layer0SkinFolder = "layer_0"
	Layer1SkinFolder = "layer_1"

	Version           = "0.0.1"
	ConfigurationFile = "config.json"
)

type skinManager struct {
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

func (sm *skinManager) setupSkinParts() error {
	necessaryParts := []*SkinPart{&sm.BasePath, &sm.HeadPath}

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

func (sm *skinManager) GenerateSkin() error {
	dir := sm.BasePath.Layer0Path

	fmt.Println(sm.BasePath.Layer0Path)

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var pngFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {

			filePath := dir + "/" + file.Name()
			inputFile, err := os.Open(filePath)
			if err != nil {
				continue
			}

			img, err := png.Decode(inputFile)
			inputFile.Close()
			if err != nil {
				continue
			}

			// it's a valid skin
			if img.Bounds().Dx() <= 64 && img.Bounds().Dy() <= 64 {
				pngFiles = append(pngFiles, filePath)
			}
		}
	}

	if len(pngFiles) == 0 {
		return fmt.Errorf("No PNG files in the BasePath directory\n")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	selectedFile := pngFiles[rng.Intn(len(pngFiles))]

	inputFile, err := os.Open(selectedFile)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Define a color transformation (example: grayscale)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := img.At(x, y)
			r, g, b, a := oldColor.RGBA()

			// Convert to grayscale
			gray := (r + g + b) / 3
			newColor := color.RGBA{
				R: uint8(gray >> 8),
				G: uint8(gray >> 8),
				B: uint8(gray >> 8),
				A: uint8(a >> 8),
			}

			newImg.Set(x, y, newColor)
		}
	}

	outputFile, err := os.Create("recolored_image.png")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, newImg); err != nil {
		panic(err)
	}

	return nil
}

func SetupSkinManager() (skinManager skinManager, err error) {
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

	config := flag.String("config", "", "JSON Config file path")
	flag.Parse()

	if config == nil || *config == "" {
		return skinManager, fmt.Errorf("The config file is not set, create one with skin-dir and randomize-dir flags")
	}

	configPath, err := utils.FormatPath(*config)
	if err != nil {
		return skinManager, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return skinManager, err
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return skinManager, err
	}

	var configFile configuration.ConfigFile

	err = json.Unmarshal(bytes, &configFile)
	if err != nil {
		return skinManager, fmt.Errorf("Your config file is not correct.")
	}

	if configFile.RandomizerFolder == "" {
		return skinManager, fmt.Errorf("Randomizer folder path is missing.")
	}

	if configFile.EditableSkin == "" {
		return skinManager, fmt.Errorf("Editable skin path is missing.")
	}

	skinManager.RandomizerFolderPath = configFile.RandomizerFolder
	skinManager.SkinPath = configFile.EditableSkin

	skinManager.RandomizerFolderPath, err = utils.FormatPath(skinManager.RandomizerFolderPath)
	if err != nil {
		return skinManager, err
	}

	skinManager.SkinPath, err = utils.FormatPath(skinManager.SkinPath)
	if err != nil {
		return skinManager, err
	}

	skinManager.BasePath.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, BaseSkinFolder)
	skinManager.HeadPath.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, HeadSkinFolder)

	err = skinManager.setupSkinParts()
	if err != nil {
		fmt.Println("error", err)
	}

	return skinManager, nil
}
