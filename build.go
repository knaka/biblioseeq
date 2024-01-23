package common

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var MainPackage string

type AirInfoParams struct {
	BuildPackage string
	TargetEnv    string
	ExcludeDirs  []string
}

var AirInfo *AirInfoParams

func baseName() string {
	return filepath.Base(MainPackage)
}

var buildDir string

func init() {
	wd, _ := os.Getwd()
	buildDir = filepath.Join(wd, ".build")
	DirsToCleanUp = append(DirsToCleanUp, buildDir)
}

//goland:noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Build mg.Namespace

// Native builds the main binary for the current environment.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Build) Native() error {
	mg.Deps(Gen)
	// Do not Deps this together because this chroot's
	mg.Deps(Client.Build)
	return sh.Run(mg.GoCmd(), "build", "-o", filepath.Join(buildDir, baseName()), MainPackage)
}

func makeBinName(baseName, targetEnv, goos, goarch string) string {
	return fmt.Sprintf("%s-%s-%s-%s", baseName, targetEnv, goos, goarch)
}

// Cross (goarch string) builds a binary for a specified architecture.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Build) Cross(goarch string) error {
	mg.Deps(Gen)
	//Ensure0(os.Symlink(
	//	fmt.Sprintf("prebuilt-%s.Dockerfile", goarch),
	//	"Dockerfile",
	//))
	// Do not Deps because this chroot's
	mg.Deps(Client.Build)
	const goos = "linux"
	if goarch == "native" {
		goarch = runtime.GOARCH
	}
	return sh.RunWith(
		map[string]string{
			"GOOS":   goos,
			"GOARCH": goarch,
		},
		mg.GoCmd(), "build",
		// command line arguments - Passing an optimization flag to a Go compiler? - Stack Overflow https://stackoverflow.com/questions/45003259/passing-an-optimization-flag-to-a-go-compiler
		// -s
		//   Omit the symbol table and debug information.
		// -w
		//   Omit the DWARF symbol table.
		"-ldflags=-s -w",
		"-o", filepath.Join(buildDir, makeBinName(baseName(), "prod", goos, goarch)),
		MainPackage,
	)
}

// Air launches the server process and keeps it live-reloading.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Air() error {
	envMap := make(map[string]string)
	binPath := filepath.Join(buildDir, makeBinName(baseName(), AirInfo.TargetEnv, runtime.GOOS, runtime.GOARCH))
	return RunWith(envMap,
		"air",
		"--build.cmd", fmt.Sprintf("go build -gcflags 'all=-N -l' -o %s %s", binPath, AirInfo.BuildPackage),
		"--build.bin", binPath,
		"--build.stop_on_error", "true",
		"--build.exclude_dir", strings.Join(AirInfo.ExcludeDirs, ","),
	)
}

// Dev is an alias for 'air'.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Dev() { mg.Deps(Air) }

// Start is an alias for 'air'.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Start() { mg.Deps(Air) }
