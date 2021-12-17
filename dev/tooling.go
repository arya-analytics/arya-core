package dev

import (
	"github.com/hashicorp/go-version"
	"log"
	"os/exec"
	"strings"
)


type  ToolingConfig map[string]string

var RequiredTools = ToolingConfig{
	"multipass": "^1.8.1",
	"kubernetes-cli": "^v1.23.0",
	"krew": "^v0.4.2",
	"yq": "^4.16.1",
	"helm": "^3.7.2",
}

type Tooling interface {
	Install(tool string) error
	CheckInstalled(tool string) bool
}

func NewTooling() Tooling {
	return &BrewTooling{RequiredTools}
}

var requiredBrewVersion, _ = version.NewVersion("3.3.8")

type BrewTooling struct {
	tooling ToolingConfig
}

// Install installs development tools into the users environment.
// By default, installs all tools and versions in tooling.config.json.
// Receives an option array of arguments specifying which tools to skip install.
func (t BrewTooling) Install (tool string) error {
	cmd := "brew install " + tool
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("%s", err)
	}
	return nil
}

func (t BrewTooling) CheckInstalled (tool string) bool {
	cmd := exec.Command("brew", "list")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	outStr := string(out[:])
	return strings.Contains(outStr, tool)
}

// checkPreReqs checks for necessary pre-requisites for the installer to run correctly
func (t BrewTooling) checkPreReqs () {
	cmd := "brew --version | grep \"Homebrew \" | awk '{print $2}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("%s", err)
	}
	outString := strings.TrimSpace(string(out[:]))
	brewVersion, err := version.NewVersion(outString)
	if err != nil {
		log.Fatalf("%s", err)
	}
	if brewVersion.LessThan(requiredBrewVersion) {
		log.Fatalf("Brew Version %s is less than the required version %s",
			outString, requiredBrewVersion.String())
	}
}