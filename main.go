package main

import (
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/unit"

	"github.com/kqnade/VRCLaunch/internal/config"
	"github.com/kqnade/VRCLaunch/internal/launcher"
	"github.com/kqnade/VRCLaunch/internal/ui"
	"github.com/kqnade/VRCLaunch/internal/uistate"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("config load: %v (continuing with defaults)", err)
		cfg = config.Default()
	}

	state := uistate.NewState(cfg)
	theme := ui.NewDarkTheme()
	uiApp := ui.NewApp(state, theme, config.Save, launcher.Launch)

	go func() {
		w := new(gioapp.Window)
		w.Option(
			gioapp.Title("VRCLaunch"),
			gioapp.Size(unit.Dp(720), unit.Dp(560)),
			gioapp.MinSize(unit.Dp(480), unit.Dp(400)),
		)
		if err := uiApp.Run(w); err != nil {
			log.Printf("app exited: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	gioapp.Main()
}
