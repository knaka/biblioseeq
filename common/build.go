package common

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/knaka/go-utils"
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

var buildDirPath string

func init() {
	wd, _ := os.Getwd()
	buildDirPath = filepath.Join(wd, ".build")
	DirsToCleanUp = append(DirsToCleanUp, buildDirPath)
}

func makeBinName(baseName, targetEnv, goos, goarch string) string {
	return fmt.Sprintf("%s-%s-%s-%s", baseName, targetEnv, goos, goarch)
}

// Build builds the binary for the current platform.
// Known architectures: https://github.com/golang/go/blob/105ac94486f243fc478c3a146d836302a95cdbbc/src/go/build/syslist.go#L54
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Build() error {
	goos := Ternary(os.Getenv("GOOS") != "",
		os.Getenv("GOOS"),
		runtime.GOOS,
	)
	goarch := Ternary(os.Getenv("GOARCH") != "",
		os.Getenv("GOARCH"),
		runtime.GOARCH,
	)
	binName := makeBinName(baseName(), "prod", goos, goarch)
	V0(sh.RunWith(
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
		"-o", filepath.Join(buildDirPath, binName),
		MainPackage,
	))
	if goos == runtime.GOOS && goarch == runtime.GOARCH {
		V0(sh.Copy(filepath.Join(buildDirPath, baseName()), filepath.Join(buildDirPath, binName)))
		// Add executable permission to the binary
		V0(sh.RunV("chmod", "+x", filepath.Join(buildDirPath, baseName())))
	}
	return nil
}

// Air launches the server process and keeps it live-reloading.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Air() error {
	envMap := make(map[string]string)
	binPath := filepath.Join(buildDirPath, makeBinName(baseName(), AirInfo.TargetEnv, runtime.GOOS, runtime.GOARCH))
	return RunWith("", envMap,
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
