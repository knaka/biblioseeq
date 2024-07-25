package conf

import (
	"bytes"
	_ "embed"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	. "github.com/knaka/go-utils"
)

func getConfigFilePath() (string, error) {
	return filepath.Join(V(userConfigDir()), "biblioseeq.toml"), nil
}

type Directory struct {
	Path           string   `toml:"path"`
	FileExtensions []string `toml:"file_extensions"`
}

type Config struct {
	Directories []*Directory `toml:"directories"`
}

var reHomeVariable = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`\$HOME\b`) })

//go:embed config-default.toml
var defaultConfigToml []byte

func ReadConfig() (config *Config, err error) {
	configFilePath := V(getConfigFilePath())
	V0(os.MkdirAll(filepath.Dir(configFilePath), 0755))
	if _, err := os.Stat(configFilePath); err != nil {
		func() {
			reader := bytes.NewReader(defaultConfigToml)
			writer := V(os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0644))
			defer (func() { Ignore(writer.Close()) })()
			V0(io.Copy(writer, reader))
		}()
	}
	configToml := V(os.ReadFile(configFilePath))
	config = &Config{}
	V0(toml.Unmarshal(configToml, config))
	homeDir := V(os.UserHomeDir())
	config.Directories = lo.Map(config.Directories, func(dir *Directory, index int) *Directory {
		dir.Path = V(filepath.EvalSymlinks(
			V(filepath.Abs(
				reHomeVariable().ReplaceAllString(dir.Path, homeDir),
			)),
		))
		dir.FileExtensions = lo.Map(dir.FileExtensions, func(ext string, index int) string {
			// Todo: Use sync value.
			ext = regexp.MustCompile(`^\*\.`).ReplaceAllString(ext, ".")
			if len(ext) > 0 && ext[0] != '.' {
				ext = "." + ext
			}
			return ext
		})
		return dir
	})
	return
}
