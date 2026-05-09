package launcher

import (
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/kqnade/VRCLaunch/internal/config"
)

func TestLaunch_EmptyLaunchPath(t *testing.T) {
	err := Launch("", config.Profile{Index: 1}, ModeDesktop)
	if !errors.Is(err, ErrLaunchPathNotSet) {
		t.Errorf("expected ErrLaunchPathNotSet, got %v", err)
	}
}

func TestLaunch_NonexistentExecutable(t *testing.T) {
	err := Launch("/nonexistent/path/to/launch.exe", config.Profile{Index: 1}, ModeDesktop)
	if err == nil {
		t.Error("expected error for nonexistent executable, got nil")
	}
	if errors.Is(err, ErrLaunchPathNotSet) {
		t.Error("should not be ErrLaunchPathNotSet for nonexistent path")
	}
}

func TestLaunch_RealExecutable(t *testing.T) {
	// Use /bin/true (or similar) as a stand-in: succeeds immediately
	truePath, err := exec.LookPath("true")
	if err != nil {
		t.Skip("'true' not available on this system")
	}
	if err := Launch(truePath, config.Profile{Index: 1}, ModeDesktop); err != nil {
		t.Errorf("expected successful spawn, got %v", err)
	}
}

func TestErrLaunchPathNotSet_Message(t *testing.T) {
	if !strings.Contains(ErrLaunchPathNotSet.Error(), "launch") {
		t.Errorf("error message should mention launch path: %s", ErrLaunchPathNotSet)
	}
}
