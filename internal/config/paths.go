package config

import (
	"os"
	"path/filepath"
)

const appDirName = "VRCLaunch"

func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appDirName), nil
}

func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
