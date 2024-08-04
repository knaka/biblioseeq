//go:build ignore
// +build ignore

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "github.com/knaka/go-utils"
)

func main() {
	cmd := exec.Command("go", "run", ".")
	_, thisFilePath, _, _ := runtime.Caller(0)
	cmd.Dir = filepath.Dir(thisFilePath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	V0(cmd.Run())
}
