package common

import (
	"fmt"
	"github.com/magefile/mage/target"
	"os"
	"path/filepath"
	"text/template"

	. "github.com/knaka/go-utils"
)

const dir = "."

type Dockerfile struct {
	Filename string
	Goarch   string
	Prebuilt bool
}

var dockerfileTemplate = "Dockerfile.tmpl"

var dockerfiles []*Dockerfile

func SetDockerfiles(dockerfiles_ ...*Dockerfile) {
	dockerfiles = dockerfiles_
}

func AddDockerfile(filename, goarch string, prebuilt bool) {
	dockerfiles = append(dockerfiles, &Dockerfile{filename, goarch, prebuilt})
}

func generateDockerfile(data *Dockerfile) (err error) {
	defer Catch(&err)
	// Struct “Template” can contain multiple templates (“defined templates”).
	tmplTree := template.New("")
	// “It must be called before the template is parsed.”
	tmplTree = tmplTree.Funcs(template.FuncMap{
		"GeneratedContentNotice": func() string {
			return fmt.Sprintf("THIS FILE IS GENERATED FROM `%s`. DO NOT EDIT DIRECTLY.", dockerfileTemplate)
		},
	})
	// The template parsed with `ParseFiles` is treated as a defined template and associated with the base name of the file.
	tmplTree = V(tmplTree.ParseFiles(filepath.Join(dir, dockerfileTemplate)))
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, data.Filename)) // Create or truncate a file
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	// Execute by the template name.
	V0(tmplTree.ExecuteTemplate(file, filepath.Base(dockerfileTemplate), data))
	return nil
}

// Dockerfiles generates dockerfiles from the template.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Dockerfiles() error {
	for _, dockerfileParams := range dockerfiles {
		if updated, _ := target.Path(dockerfileParams.Filename, dockerfileTemplate); updated {
			err := generateDockerfile(dockerfileParams)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	AddGenFn(Dockerfiles)
}
