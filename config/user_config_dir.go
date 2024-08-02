//go:build !darwin

package config

import "os"

func userConfigDir() (string, error) {
	return os.UserConfigDir()
}
