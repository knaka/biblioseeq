package common

import (
	"fmt"
	"github.com/magefile/mage/target"
	"os"
	"path/filepath"
	"text/template"
)

const dir = "."

type DockerfileParams struct {
	Filename string
	Goarch   string
	Prebuilt bool
}

var dockerfileTemplate = "template.dockerfile"

var DockerfileSettings []*DockerfileParams

func generateDockerfile(data *DockerfileParams) error {
	tmpl, err := template.ParseFiles(filepath.Join(dir, dockerfileTemplate))
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, data.Filename)) // Create or truncate a file
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	_, _ = file.WriteString(
		fmt.Sprintf("# THIS FILE IS GENERATED FROM `%s`. DO NOT EDIT DIRECTLY.\n", dockerfileTemplate),
	)
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}

// Dockerfiles generates dockerfiles from the template.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Dockerfiles() error {
	for _, dockerfileParams := range DockerfileSettings {
		if updated, _ := target.Path(dockerfileParams.Filename, dockerfileTemplate); updated {
			err := generateDockerfile(dockerfileParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
