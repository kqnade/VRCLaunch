package launcher

import (
	"slices"
	"testing"

	"github.com/kqnade/VRCLaunch/internal/config"
)

func argsValueAfter(args []string, key string) (string, bool) {
	for i, a := range args {
		if a == key && i+1 < len(args) {
			return args[i+1], true
		}
	}
	return "", false
}

func TestBuildArgs_AlwaysIncludesProfileIndex(t *testing.T) {
	p := config.Profile{Index: 3}
	got := BuildArgs(p, ModeDesktop)
	if !slices.Contains(got, "--profile=3") {
		t.Errorf("expected --profile=3 in args, got %v", got)
	}
}

func TestBuildArgs_VRModeIncludesVRFlag(t *testing.T) {
	got := BuildArgs(config.Profile{Index: 1}, ModeVR)
	if !slices.Contains(got, "--vr") {
		t.Errorf("expected --vr in args, got %v", got)
	}
	if slices.Contains(got, "--no-vr") {
		t.Errorf("--no-vr should not appear in VR mode, got %v", got)
	}
}

func TestBuildArgs_DesktopModeIncludesNoVRFlag(t *testing.T) {
	got := BuildArgs(config.Profile{Index: 1}, ModeDesktop)
	if !slices.Contains(got, "--no-vr") {
		t.Errorf("expected --no-vr in args, got %v", got)
	}
	if slices.Contains(got, "--vr") {
		t.Errorf("--vr should not appear in Desktop mode, got %v", got)
	}
}

func TestBuildArgs_VRModeForcesMinimumScreen(t *testing.T) {
	p := config.Profile{
		Index: 1,
		Options: config.ProfileOptions{
			ScreenWidth:      1920,
			ScreenHeight:     1080,
			ScreenFullscreen: true,
		},
	}
	got := BuildArgs(p, ModeVR)

	w, ok := argsValueAfter(got, "-screen-width")
	if !ok || w != "320" {
		t.Errorf("VR -screen-width: got %q (ok=%v), want 320", w, ok)
	}
	h, ok := argsValueAfter(got, "-screen-height")
	if !ok || h != "240" {
		t.Errorf("VR -screen-height: got %q (ok=%v), want 240", h, ok)
	}
	if slices.Contains(got, "-screen-fullscreen") {
		t.Errorf("VR mode should not use fullscreen, got %v", got)
	}
}

func TestBuildArgs_DesktopHonorsScreenOptions(t *testing.T) {
	p := config.Profile{
		Index: 2,
		Options: config.ProfileOptions{
			ScreenWidth:      1280,
			ScreenHeight:     720,
			ScreenFullscreen: true,
		},
	}
	got := BuildArgs(p, ModeDesktop)

	if w, ok := argsValueAfter(got, "-screen-width"); !ok || w != "1280" {
		t.Errorf("Desktop -screen-width: got %q (ok=%v), want 1280", w, ok)
	}
	if h, ok := argsValueAfter(got, "-screen-height"); !ok || h != "720" {
		t.Errorf("Desktop -screen-height: got %q (ok=%v), want 720", h, ok)
	}
	if fs, ok := argsValueAfter(got, "-screen-fullscreen"); !ok || fs != "1" {
		t.Errorf("Desktop -screen-fullscreen: got %q (ok=%v), want 1", fs, ok)
	}
}

func TestBuildArgs_DesktopOmitsScreenWhenUnset(t *testing.T) {
	got := BuildArgs(config.Profile{Index: 1}, ModeDesktop)
	if slices.Contains(got, "-screen-width") {
		t.Errorf("expected no -screen-width when unset, got %v", got)
	}
	if slices.Contains(got, "-screen-height") {
		t.Errorf("expected no -screen-height when unset, got %v", got)
	}
}

func TestBuildArgs_FPSIncludedWhenPositive(t *testing.T) {
	p := config.Profile{Index: 1, Options: config.ProfileOptions{FPS: 90}}
	got := BuildArgs(p, ModeDesktop)
	if !slices.Contains(got, "--fps=90") {
		t.Errorf("expected --fps=90, got %v", got)
	}
}

func TestBuildArgs_FPSOmittedWhenZero(t *testing.T) {
	got := BuildArgs(config.Profile{Index: 1}, ModeDesktop)
	for _, a := range got {
		if len(a) >= 6 && a[:6] == "--fps=" {
			t.Errorf("expected no --fps when 0, got %v", got)
		}
	}
}

func TestBuildArgs_CustomArgsSplitOnWhitespace(t *testing.T) {
	p := config.Profile{
		Index: 1,
		Options: config.ProfileOptions{CustomArgs: "  --foo   --bar=baz  "},
	}
	got := BuildArgs(p, ModeDesktop)
	if !slices.Contains(got, "--foo") || !slices.Contains(got, "--bar=baz") {
		t.Errorf("expected --foo and --bar=baz from custom args, got %v", got)
	}
}

func TestBuildArgs_CustomArgsEmpty(t *testing.T) {
	p := config.Profile{
		Index:   1,
		Options: config.ProfileOptions{CustomArgs: "   "},
	}
	got := BuildArgs(p, ModeDesktop)
	want := []string{"--profile=1", "--no-vr"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
