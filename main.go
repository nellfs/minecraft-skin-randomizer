package main

import (
	"log"

	"github.com/nellfs/minecraft-skin-randomizer/skin"
)

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
