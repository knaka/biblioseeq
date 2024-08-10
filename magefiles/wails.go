package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/knaka/go-utils"
)

type Wails mg.Namespace

// Build builds for Wails.
func (Wails) Build() (err error) {
	mg.Deps(Client.Build)
	return gobinRun("wails", "wails", "build")
}

// Dev starts Wails development server.
func (Wails) Dev() (err error) {
	mg.Deps(Client.Build)
	return gobinRun("wails", "wails", "dev")
}

// Start starts Wails development server.
func (Wails) Start() (err error) {
	mg.Deps(Client.Build)
	return gobinRun("wails", "wails", "dev")
}

// Run runs Wails application.
func (Wails) Run() (err error) {
	switch goos := runtime.GOOS; goos {
	case "darwin":
		appPath := filepath.Join(V(os.Getwd()), "wails", "build", "bin", "BiblioSeeQ.app")
		// Should I always build?
		//if stat, err_ := os.Stat(appPath); err_ != nil || !stat.IsDir() {
		mg.Deps(Wails.Build)
		//}
		return run("wails", "open", appPath)
	default:
		return fmt.Errorf("unsupported OS: %s", goos)
	}
}
