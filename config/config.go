package config

const (
	Version           = "0.0.1"
	ConfigurationFile = "config.json"
)

type ConfigFile struct {
	ConfigVersion    string           `json:"version"`
	EditableSkin     string           `json:"edit_skin"`
	RandomizerFolder string           `json:"randomizer_folder"`
	LastGeneration   string           `json:"last_generation"`
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
