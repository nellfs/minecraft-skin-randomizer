package main

import (
	"log"

	"github.com/nellfs/minecraft-skin-randomizer/skin"
)

//TODO:
//make the first release
//create config.json if it does not exist or the flag is not set, use an interactive setup
//follow the config.json layer enabling
//- skin preview in terminal?

func main() {
	skinManager, err := skin.SetupSkinManager()
	if err != nil {
		log.Fatalf("Could not set up a new skin manager: %v", err)
	}

	err = skinManager.GenerateSkin()
	if err != nil {
		log.Fatalf("Could not generate a new skin: %v", err)
	}
}
