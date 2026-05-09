package launcher

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/kqnade/VRCLaunch/internal/config"
)

var ErrLaunchPathNotSet = errors.New("launch path is not set")

func Launch(launchPath string, p config.Profile, mode Mode) error {
	if launchPath == "" {
		return ErrLaunchPathNotSet
	}
	args := BuildArgs(p, mode)
	cmd := exec.Command(launchPath, args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("spawn launch: %w", err)
	}
	// Detach: do not wait. Release the process so it survives independently.
	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("release process: %w", err)
	}
	return nil
}
