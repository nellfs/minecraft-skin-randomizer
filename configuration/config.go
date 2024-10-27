package configuration

type ConfigFile struct {
	ConfigVersion    string `json:"version"`
	EditableSkin     string `json:"edit_skin"`
	RandomizerFolder string `json:"randomizer_folder"`
	LastGeneration   string `json:"last_generation"`
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
