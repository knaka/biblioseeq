package conf

import (
	"bytes"
	_ "embed"
	. "github.com/knaka/go-utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func getConfigFilePath() (string, error) {
	return filepath.Join(V(userConfigDir()), "biblioseeq.toml"), nil
}

type directory struct {
	Path           string   `toml:"path"`
	FileExtensions []string `toml:"file_extensions"`
}

type data struct {
	Directories []*directory `toml:"directories"`
}

//go:embed config-default.toml
var defaultConfigToml []byte

var reHomeVariable = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`\$HOME\b`) })
var reTrailingSlashes = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`[/\\]+$`) })
var reExtWildcard = sync.OnceValue(func() *regexp.Regexp { return regexp.MustCompile(`^\*\.`) })

//goland:noinspection GoExportedFuncWithUnexportedType
func Read(configFilePath string) (dat *data, err error) {
	if configFilePath == "" {
		configFilePath = V(getConfigFilePath())
	}
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
	dat = &data{}
	V0(toml.Unmarshal(configToml, dat))
	homeDir := V(os.UserHomeDir())
	dat.Directories = lo.Map(dat.Directories, func(dir *directory, index int) *directory {
		dir.Path = filepath.Clean(V(filepath.EvalSymlinks(
			V(filepath.Abs(
				reHomeVariable().ReplaceAllString(
					reTrailingSlashes().ReplaceAllString(dir.Path, ""),
					homeDir,
				),
			)),
		)))
		dir.FileExtensions = lo.Map(dir.FileExtensions, func(ext string, index int) string {
			ext = reExtWildcard().ReplaceAllString(ext, ".")
			if len(ext) > 0 && ext[0] != '.' {
				ext = "." + ext
			}
			return ext
		})
		return dir
	})
	return
}
