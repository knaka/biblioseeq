//go:build !darwin

package conf

import "os"

func userConfigDir() (string, error) {
	return os.UserConfigDir()
}
