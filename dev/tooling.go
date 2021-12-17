package dev

import (
	"github.com/hashicorp/go-version"
	"log"
	"os/exec"
	"strings"
)


type  ToolingConfig map[string]string

var requiredTools = ToolingConfig{
	"multipass": "^1.8.1",
	"kubernetes-cli": "^v1.23.0",
	"krew": "v0.4.2",
}

type Tooling interface {
	Install() error
}

func NewInstaller () Tooling {
	return &BrewTooling{requiredTools}
}

var requiredBrewVersion, _ = version.NewVersion("3.3.8")

type BrewTooling struct {
	tooling ToolingConfig
}

// Install installs development tools into the users environment.
// By default, installs all tools and versions in tooling.config.json.
// Receives an option array of arguments specifying which tools to skip install.
func (i BrewTooling) Install () error {
	i.checkPreReqs()
	for k := range i.tooling {
		cmd := "brew install " + k
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatalf("%s", err)
		}
		log.Println(string(out[:]))
	}
	return nil
}

// checkPreReqs checks for necessary pre-requisites for the installer to run correctly
func (i BrewTooling) checkPreReqs () {
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