package skin

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nellfs/minecraft-skin-randomizer/config"
	"github.com/nellfs/minecraft-skin-randomizer/utils"
)

const (
	SkinsFolder = "/skins"

	Layer0SkinFolder = "layer_0"
	Layer1SkinFolder = "layer_1"

	BaseSkinFolder     = "/base"
	HeadSkinFolder     = "/head"
	BodySkinFolder     = "/body"
	LeftArmSkinFolder  = "/left_arm"
	RightArmSkinFolder = "/righ_arm"
	LeftLegSkinFolder  = "/left_leg"
	RightLegSkinFolder = "/right_leg"
)

type skinManager struct {
	Config config.ConfigFile

	RandomizerFolderPath string
	SkinPath             string

	BasePart     SkinPart
	HeadPart     SkinPart
	BodyPart     SkinPart
	LeftArmPart  SkinPart
	RightArmPart SkinPart
	LeftLegPart  SkinPart
	RightLegPart SkinPart
}

type SkinPart struct {
	RootPath   string
	Layer0Path string
	Layer1Path string
}

func (sm *skinManager) setupSkinParts() error {
	necessaryParts := []*SkinPart{&sm.BasePart, &sm.HeadPart, &sm.BodyPart, &sm.LeftArmPart, &sm.RightArmPart, &sm.LeftLegPart, &sm.RightLegPart}

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

	configFlag := flag.String("config", "", "JSON Config file path")
	flag.Parse()

	if configFlag == nil || *configFlag == "" {
		return skinManager, fmt.Errorf("The config file is not set, create one with skin-dir and randomize-dir flags")
	}

	configPath, err := utils.FormatPath(*configFlag)
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

	err = json.Unmarshal(bytes, &skinManager.Config)
	if err != nil {
		return skinManager, fmt.Errorf("Your config file is not correct.")
	}

	if skinManager.Config.RandomizerFolder == "" {
		return skinManager, fmt.Errorf("Randomizer folder path is missing.")
	}

	if skinManager.Config.EditableSkin == "" {
		return skinManager, fmt.Errorf("Editable skin path is missing.")
	}

	skinManager.RandomizerFolderPath, err = utils.FormatPath(skinManager.Config.RandomizerFolder)
	if err != nil {
		return skinManager, err
	}

	skinManager.SkinPath, err = utils.FormatPath(skinManager.Config.EditableSkin)
	if err != nil {
		return skinManager, err
	}

	skinManager.BasePart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, BaseSkinFolder)
	skinManager.HeadPart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, HeadSkinFolder)
	skinManager.BodyPart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, BodySkinFolder)
	skinManager.LeftArmPart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, LeftArmSkinFolder)
	skinManager.RightArmPart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, RightArmSkinFolder)
	skinManager.LeftLegPart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, LeftLegSkinFolder)
	skinManager.RightLegPart.RootPath = fmt.Sprintf("%s%s%s", skinManager.RandomizerFolderPath, SkinsFolder, RightLegSkinFolder)

	err = skinManager.setupSkinParts()
	if err != nil {
		return skinManager, err
	}

	return skinManager, nil
}

func (sm *skinManager) GenerateSkin() error {
	//TODO: change to use only parts in the config

	//Order
	parts := []SkinPart{sm.BasePart, sm.HeadPart, sm.BodyPart, sm.LeftArmPart, sm.RightArmPart, sm.LeftLegPart, sm.RightLegPart}

	finalMix := []string{}

	for _, part := range parts {
		dirLayer0 := part.Layer0Path
		dirLayer1 := part.Layer1Path

		filesLayer0, err := os.ReadDir(dirLayer0)
		if err != nil {
			return err
		}
		filesLayer1, err := os.ReadDir(dirLayer1)
		if err != nil {
			return err
		}

		var layer0Files []string
		var layer1Files []string

		for _, file := range filesLayer0 {
			if filepath.Ext(file.Name()) == ".png" {
				skinPath := dirLayer0 + "/" + file.Name()
				err = utils.VerifySkin(skinPath)
				if err != nil {
					return fmt.Errorf("Layer 0 is in the wrong size: %v (%s)", err, file.Name())
				}
				layer0Files = append(layer0Files, skinPath)
			}
		}
		for _, file := range filesLayer1 {
			if filepath.Ext(file.Name()) == ".png" {
				skinPath := dirLayer1 + "/" + file.Name()
				err = utils.VerifySkin(skinPath)
				if err != nil {
					return fmt.Errorf("Layer 1 is in the wrong size: %v (%s)", err, file.Name())
				}
				layer1Files = append(layer1Files, skinPath)
			}
		}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)

		fmt.Println("Mix:")
		if len(layer0Files) != 0 {
			randomIndexLayer0 := r.Intn(len(layer0Files))
			randomLayer0Part := layer0Files[randomIndexLayer0]

			finalMix = append(finalMix, randomLayer0Part)
			fmt.Println(randomLayer0Part)
		}
		if len(layer1Files) != 0 {
			randomIndexLayer1 := r.Intn(len(layer1Files))
			randomLayer1Part := layer1Files[randomIndexLayer1]
			finalMix = append(finalMix, randomLayer1Part)
			fmt.Println(randomLayer1Part)
		}

		fmt.Println("--------------------")

	}

	skin, err := utils.LoadImage(finalMix[0]) // -> base skin layer 0
	if err != nil {
		return fmt.Errorf("Error loading base skin image: %v\n", err)
	}

	for _, layerPath := range finalMix {
		currentSkinPart, err := utils.LoadImage(layerPath)
		if err != nil {
			return err
		}

		tempImage := image.NewRGBA(skin.Bounds())
		draw.Draw(tempImage, skin.Bounds(), skin, image.Point{}, draw.Over)
		draw.Draw(tempImage, currentSkinPart.Bounds().Add(image.Point{0, 0}), currentSkinPart, image.Point{}, draw.Over)

		skin = tempImage
	}

	outputFile, err := os.Create(sm.Config.EditableSkin)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v\n", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, skin)
	if err != nil {
		return fmt.Errorf("Error saving final image: %v\n", err)
	}

	return nil
}
