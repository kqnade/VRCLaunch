package launcher

import (
	"fmt"
	"strings"

	"github.com/kqnade/VRCLaunch/internal/config"
)

type Mode int

const (
	ModeDesktop Mode = iota
	ModeVR
)

const (
	vrMinScreenWidth  = 320
	vrMinScreenHeight = 240
)

func BuildArgs(p config.Profile, mode Mode) []string {
	args := []string{fmt.Sprintf("--profile=%d", p.Index)}

	switch mode {
	case ModeVR:
		args = append(args, "--vr")
		args = append(args, "-screen-width", fmt.Sprintf("%d", vrMinScreenWidth))
		args = append(args, "-screen-height", fmt.Sprintf("%d", vrMinScreenHeight))
	case ModeDesktop:
		args = append(args, "--no-vr")
		if p.Options.ScreenWidth > 0 {
			args = append(args, "-screen-width", fmt.Sprintf("%d", p.Options.ScreenWidth))
		}
		if p.Options.ScreenHeight > 0 {
			args = append(args, "-screen-height", fmt.Sprintf("%d", p.Options.ScreenHeight))
		}
		if p.Options.ScreenFullscreen {
			args = append(args, "-screen-fullscreen", "1")
		}
	}

	if p.Options.FPS > 0 {
		args = append(args, fmt.Sprintf("--fps=%d", p.Options.FPS))
	}

	if custom := strings.TrimSpace(p.Options.CustomArgs); custom != "" {
		args = append(args, strings.Fields(custom)...)
	}

	return args
}
