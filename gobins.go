package common

import (
	"errors"
	"fmt"
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mattn/go-shellwords"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var goBinDir string

type GobinPkgParams struct {
	Name    string
	Version string
	Tags    string
}

func init() {
	// The `mage` command only runs in directories containing magefiles or a `magefiles/` directory.
	wd := Ensure(os.Getwd())
	goBinDir = filepath.Join(wd, ".gobin")
	DirsToCleanUp = append(DirsToCleanUp, goBinDir)
}

var GobinPkgs []*GobinPkgParams

// RunWith runs the given command prioritizing binaries in .gobin/ directory. It is not a task function.
func RunWith(env map[string]string, cmd string, args ...string) error {
	_, err := exec.LookPath(cmd)
	if err != nil {
		err := Gobin(cmd)
		if err != nil {
			return err
		}
	}
	Assert(os.Setenv("PATH", goBinDir+string(os.PathListSeparator)+os.Getenv("PATH")))
	return sh.RunWith(env, cmd, args...)
}

// Exec (command string) executes the given command prioritizing binaries in .gobin/ directory.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Exec(command string) error {
	envs, args, err := shellwords.ParseWithEnvs(command)
	if err != nil {
		log.Panicf("Failed to parse commandline (%v)", err)
	}
	envMap := make(map[string]string)
	for _, env := range envs {
		fields := strings.Split(env, "=")
		envMap[fields[0]] = fields[1]
	}
	return RunWith(envMap, args[0], args[1:]...)
}

func ensureGobinInstalled(pkgName, version, tags string) error {
	cmdName := filepath.Base(pkgName)
	binPath := filepath.Join(goBinDir, fmt.Sprintf("%s@%s", cmdName, version))
	if _, err := os.Stat(binPath); err == nil {
		return nil
	}
	pkgNameWithVer := pkgName + "@" + version
	Assert(fmt.Fprintf(os.Stderr, "Building %s\n", pkgNameWithVer))
	linkTgtPath := filepath.Join(goBinDir, cmdName)
	Assert(os.Remove(linkTgtPath))
	Assert(sh.RunWith(
		map[string]string{"GOBIN": goBinDir},
		mg.GoCmd(), "install", "-tags", tags, pkgNameWithVer,
	))
	Assert(os.Rename(linkTgtPath, binPath))
	Assert(os.Symlink(filepath.Base(binPath), linkTgtPath))
	return nil
}

// Gobins installs necessary Go binaries.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gobins() {
	Assert(os.MkdirAll(goBinDir, 0755))
	for _, pkg := range GobinPkgs {
		Assert(ensureGobinInstalled(pkg.Name, pkg.Version, pkg.Tags))
	}
}

// Gobin installs the Go binary with the given command name or package name.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gobin(name string) error {
	for _, pkg := range GobinPkgs {
		if name == pkg.Name || name == filepath.Base(pkg.Name) {
			return ensureGobinInstalled(pkg.Name, pkg.Version, pkg.Tags)
		}
	}
	return errors.New("no such package")
}
