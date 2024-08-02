//go:build darwin

package config

import (
	"os"
	"path/filepath"

	. "github.com/knaka/go-utils"
)

func userConfigDir() (string, error) {
	return Elvis(os.Getenv("XDG_CONFIG_HOME"), filepath.Join(V(os.UserHomeDir()), ".config")), nil
}
