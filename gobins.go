package common

import (
	"errors"
	"fmt"
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

type GobinPackage struct {
	Name    string
	Version string
	Tags    string
}

func init() {
	// The `mage` command only runs in directories containing magefiles or a `magefiles/` directory.
	wd := ensure(os.Getwd())
	goBinDir = filepath.Join(wd, ".gobin")
	DirsToCleanUp = append(DirsToCleanUp, goBinDir)
}

var gobinPackages []GobinPackage

func InitGobins(packages ...GobinPackage) {
	gobinPackages = packages
}

func ExecWith(env map[string]string, cmd string, args ...string) error {
	mg.Deps(Gobins)
	_, err := exec.LookPath(cmd)
	if err != nil {
		err = Gobin(cmd)
		if err != nil {
			return err
		}
	}
	_ = os.Setenv("PATH", goBinDir+string(os.PathListSeparator)+os.Getenv("PATH"))
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
	return ExecWith(envMap, args[0], args[1:]...)
}

func ensureGobinInstalled(pkgName, version, tags string) error {
	cmdName := filepath.Base(pkgName)
	binPath := filepath.Join(goBinDir, fmt.Sprintf("%s@%s", cmdName, version))
	if _, err := os.Stat(binPath); err == nil {
		return nil
	}
	pkgNameWithVer := pkgName + "@" + version
	assert(fmt.Fprintf(os.Stderr, "Building %s\n", pkgNameWithVer))
	linkTgtPath := filepath.Join(goBinDir, cmdName)
	assert(os.Remove(linkTgtPath))
	assert(sh.RunWith(
		map[string]string{"GOBIN": goBinDir},
		mg.GoCmd(), "install", "-tags", tags, pkgNameWithVer,
	))
	assert(os.Rename(linkTgtPath, binPath))
	assert(os.Symlink(filepath.Base(binPath), linkTgtPath))
	return nil
}

// Gobins installs necessary Go binaries.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gobins() {
	_ = os.MkdirAll(goBinDir, 0755)
	log.Println(gobinPackages)
	for _, pkg := range gobinPackages {
		assert(ensureGobinInstalled(pkg.Name, pkg.Version, pkg.Tags))
	}
}

// Gobin installs the Go binary with the given command name or package name.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gobin(name string) error {
	log.Println("1ff5f77", name)
	for _, pkg := range gobinPackages {
		if name == pkg.Name || name == filepath.Base(pkg.Name) {
			return ensureGobinInstalled(pkg.Name, pkg.Version, pkg.Tags)
		}
	}
	return errors.New("no such package")
}
