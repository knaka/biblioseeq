//go:build ignore
// +build ignore

package main

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/knaka/go-utils"
)

func main() {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = filepath.Join(V(os.Getwd()), "gen-from-mod")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	V0(cmd.Run())
}
