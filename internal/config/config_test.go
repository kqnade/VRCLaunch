package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNextProfileIndex(t *testing.T) {
	tests := []struct {
		name     string
		profiles []Profile
		want     int
	}{
		{"empty", nil, 1},
		{"sequential", []Profile{{Index: 1}, {Index: 2}}, 3},
		{"gap at start", []Profile{{Index: 2}, {Index: 3}}, 1},
		{"gap in middle", []Profile{{Index: 1}, {Index: 3}}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{Profiles: tt.profiles}
			if got := c.NextProfileIndex(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	c := &Config{Profiles: []Profile{
		{ID: "a", Name: "A"},
		{ID: "b", Name: "B"},
	}}
	if got := c.FindByID("b"); got == nil || got.Name != "B" {
		t.Errorf("FindByID b = %+v", got)
	}
	if got := c.FindByID("missing"); got != nil {
		t.Errorf("FindByID missing should be nil, got %+v", got)
	}
}

func TestNewIDUnique(t *testing.T) {
	seen := make(map[string]bool)
	for range 100 {
		id := NewID()
		if seen[id] {
			t.Fatalf("duplicate id: %s", id)
		}
		seen[id] = true
	}
}

func TestDefault(t *testing.T) {
	c := Default()
	if c.Version != SchemaVersion {
		t.Errorf("Version = %d, want %d", c.Version, SchemaVersion)
	}
	if c.Profiles == nil {
		t.Error("Profiles should be empty slice, got nil")
	}
	if len(c.Profiles) != 0 {
		t.Errorf("Profiles len = %d, want 0", len(c.Profiles))
	}
}

func TestLoadFrom_FileNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	c, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Version != SchemaVersion || len(c.Profiles) != 0 {
		t.Errorf("expected default config, got %+v", c)
	}
}

func TestLoadFrom_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{not valid"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFrom(path); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoadFrom_MissingVersionFilledIn(t *testing.T) {
	path := filepath.Join(t.TempDir(), "novers.json")
	if err := os.WriteFile(path, []byte(`{"launch_path":"/x"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}
	if c.Version != SchemaVersion {
		t.Errorf("Version = %d, want %d", c.Version, SchemaVersion)
	}
	if c.Profiles == nil {
		t.Error("Profiles should be non-nil even when omitted")
	}
}

func TestSaveTo_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "config.json")
	c := &Config{Version: 1, LaunchPath: "/path/launch.exe"}

	if err := SaveTo(path, c); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not written: %v", err)
	}
}

func TestConfigDirAndFile(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG_CONFIG_HOME-based test only runs on Linux")
	}
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	dir, err := ConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != filepath.Join(tmp, appDirName) {
		t.Errorf("ConfigDir = %q, want %q", dir, filepath.Join(tmp, appDirName))
	}

	file, err := ConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	if file != filepath.Join(dir, "config.json") {
		t.Errorf("ConfigFile = %q, want %q", file, filepath.Join(dir, "config.json"))
	}
}

func TestLoadSave_TopLevel(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG_CONFIG_HOME-based test only runs on Linux")
	}
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	c := &Config{Version: 1, LaunchPath: "/x/launch.exe"}
	if err := Save(c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.LaunchPath != c.LaunchPath {
		t.Errorf("LaunchPath roundtrip: got %q, want %q", loaded.LaunchPath, c.LaunchPath)
	}
}

func TestLoad_DefaultsWhenMissing(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG_CONFIG_HOME-based test only runs on Linux")
	}
	tmp := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmp)

	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.Version != SchemaVersion {
		t.Errorf("Version = %d, want %d", c.Version, SchemaVersion)
	}
}

func TestSaveTo_RenameFailure(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config.json")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := SaveTo(target, &Config{}); err == nil {
		t.Error("expected error when target path is a directory")
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	original := &Config{
		Version:      1,
		LaunchPath:   "/games/launch.exe",
		LastSelected: "abc",
		Profiles: []Profile{
			{
				ID:    "abc",
				Name:  "Main",
				Index: 1,
				Options: ProfileOptions{
					FPS:              90,
					ScreenWidth:      1920,
					ScreenHeight:     1080,
					ScreenFullscreen: true,
					CustomArgs:       "--foo --bar",
				},
			},
		},
	}

	if err := SaveTo(path, original); err != nil {
		t.Fatalf("SaveTo: %v", err)
	}
	loaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom: %v", err)
	}

	if loaded.LaunchPath != original.LaunchPath {
		t.Errorf("LaunchPath: got %q, want %q", loaded.LaunchPath, original.LaunchPath)
	}
	if loaded.LastSelected != original.LastSelected {
		t.Errorf("LastSelected: got %q, want %q", loaded.LastSelected, original.LastSelected)
	}
	if len(loaded.Profiles) != 1 {
		t.Fatalf("Profiles len = %d, want 1", len(loaded.Profiles))
	}
	got := loaded.Profiles[0]
	want := original.Profiles[0]
	if got != want {
		t.Errorf("profile mismatch:\n got: %+v\nwant: %+v", got, want)
	}
}
