package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nellfs/minecraft-skin-randomizer/utils"
)

const (
	CurrentConfigVersion = "0.0.1"
	ConfigurationFile    = "config.json"
)

type ConfigFile struct {
	ConfigVersion    string `json:"version"`
	EditableSkin     string `json:"edit_skin"`
	RandomizerFolder string `json:"randomizer_folder"`
	LastGeneration   string `json:"last_generation"`
	//TODO: implement that:
	GenerationConfig GenerationConfig `json:"gen_config"`
}

type GenerationConfig struct {
	BaseGenConfig     PartConfig `json:"base"`
	HeadGenConfig     PartConfig `json:"head"`
	BodyGenConfig     PartConfig `json:"body"`
	LeftArmGenConfig  PartConfig `json:"left_arm"`
	RightArmGenConfig PartConfig `json:"right_arm"`
	LeftLegGenConfig  PartConfig `json:"left_leg"`
	RightLegGenConfig PartConfig `json:"right_leg"`
}

type PartConfig struct {
	Layer1Enabled bool `json:"layer_0"`
	Layer2Enabled bool `json:"layer_1"`
}

func CreateConfigFile() (configFile ConfigFile, err error) {
	var randomizerFolder string
	var editSkin string

	fmt.Println("> Please select a skin template folder (this folder contains the different parts and layers for the random generation; check the project's README if you have any questions).")
	fmt.Scan(&randomizerFolder)
	randomizerFolder, err = utils.FormatPath(randomizerFolder)
	if err != nil {
		return configFile, err
	}
	fmt.Println("> Please provide the path of an existing skin to be overwritten by the generated skin.")
	fmt.Scan(&editSkin)
	editSkin, err = utils.FormatPath(editSkin)
	if err != nil {
		return configFile, err
	}

	configFile = ConfigFile{
		ConfigVersion:    CurrentConfigVersion,
		EditableSkin:     editSkin,
		RandomizerFolder: randomizerFolder,
	}

	file, err := os.Create(ConfigurationFile)
	if err != nil {
		return configFile, fmt.Errorf("Could not create the config file: %v\n", err)
	}

	defer file.Close()

	jsonConfig, err := json.MarshalIndent(&configFile, "", "  ")
	if err != nil {
		return configFile, fmt.Errorf("Could not unmarshal the config file: %v\n", err)
	}

	_, err = file.Write(jsonConfig)
	if err != nil {
		return configFile, fmt.Errorf("Error writing the config file: %v\n", err)
	}

	fmt.Printf("\nConfiguration file saved in the current directory (config.json), run the generator again with --config=/path/to/your/config.json\n\n")

	return configFile, nil
}
