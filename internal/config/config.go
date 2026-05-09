package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const SchemaVersion = 1

type ProfileOptions struct {
	FPS              int    `json:"fps,omitempty"`
	ScreenWidth      int    `json:"screen_width,omitempty"`
	ScreenHeight     int    `json:"screen_height,omitempty"`
	ScreenFullscreen bool   `json:"screen_fullscreen,omitempty"`
	CustomArgs       string `json:"custom_args,omitempty"`
}

type Profile struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Index   int            `json:"index"`
	Options ProfileOptions `json:"options"`
}

type Config struct {
	Version      int       `json:"version"`
	LaunchPath   string    `json:"launch_path"`
	LastSelected string    `json:"last_selected,omitempty"`
	Profiles     []Profile `json:"profiles"`
}

func Default() *Config {
	return &Config{Version: SchemaVersion, Profiles: []Profile{}}
}

func NewID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "id-fallback"
	}
	return hex.EncodeToString(b)
}

func Load() (*Config, error) {
	path, err := ConfigFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Default(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if c.Version == 0 {
		c.Version = SchemaVersion
	}
	if c.Profiles == nil {
		c.Profiles = []Profile{}
	}
	return &c, nil
}

func Save(c *Config) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if err := ensureDir(dir); err != nil {
		return fmt.Errorf("mkdir config: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := filepath.Join(dir, "config.json.tmp")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	final := filepath.Join(dir, "config.json")
	if err := os.Rename(tmp, final); err != nil {
		return fmt.Errorf("commit config: %w", err)
	}
	return nil
}

func (c *Config) NextProfileIndex() int {
	used := make(map[int]bool, len(c.Profiles))
	for _, p := range c.Profiles {
		used[p.Index] = true
	}
	for i := 1; ; i++ {
		if !used[i] {
			return i
		}
	}
}

func (c *Config) FindByID(id string) *Profile {
	for i := range c.Profiles {
		if c.Profiles[i].ID == id {
			return &c.Profiles[i]
		}
	}
	return nil
}
