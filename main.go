package main

import (
	"fmt"
	"log/slog"

	"github.com/nellfs/minecraft-skin-randomizer/skin"
)

func main() {
	skinManager, err := skin.SetupSkinManager()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not set up a new skin manager: %v", err))
	}

	err = skinManager.GenerateSkin()
	if err != nil {
		slog.Error(fmt.Sprintf("Could not generate a new skin: %v", err))
	}
}
