package main

import "github.com/magefile/mage/mg"

type Webui mg.Namespace

// Build builds for WebUI.
func (Webui) Build() (err error) {
	mg.Deps(Client.Build)
	return run("webui", "go", "build", "-o", "webui", ".")
}
