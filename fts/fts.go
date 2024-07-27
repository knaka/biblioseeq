package fts

import (
	"fmt"
	. "github.com/knaka/go-utils"
	"os"
	"path/filepath"
)

func GetDefaultDBFilePath() (string, error) {
	return filepath.Join(V(os.UserCacheDir()), fmt.Sprintf("%s.sqlite3"), appName), nil
}
