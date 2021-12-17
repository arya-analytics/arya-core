package dev

import (
	"github.com/hashicorp/go-version"
	"log"
	"os"
	"os/exec"
	"strings"
)


type  ToolingConfig []string

var RequiredTools = ToolingConfig{
	"multipass",
	"kubernetes-cli",
	"krew",
	"yq",
	"helm",
}

type Tooling interface {
	Install(tool string) error
	Uninstall(tool string) error
	CheckInstalled(tool string) bool
}

func NewTooling() Tooling {
	t := BrewTooling{RequiredTools}
	t.checkPreReqs()
	return &t
}

var requiredBrewVersion, _ = version.NewVersion("3.3.8")

type BrewTooling struct {
	tooling ToolingConfig
}

// Install installs development tools into the users environment.
// By default, installs all tools and versions in tooling.config.json.
// Receives an option array of arguments specifying which tools to skip install.
func (t BrewTooling) Install (tool string) error {
	_, err := t.command("brew install " + tool)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return nil
}

func (t BrewTooling) command (cmdString string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", cmdString)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func (t BrewTooling) Uninstall (tool string) error {
	_, err := t.command("brew uninstall " + tool)
	return err
}

func (t BrewTooling) CheckInstalled (tool string) bool {
	out, err := t.command("brew list")
	if err != nil {
		panic(err)
	}
	outStr := string(out[:])
	return strings.Contains(outStr, tool)
}

// checkPreReqs checks for necessary pre-requisites for the installer to run correctly
func (t BrewTooling) checkPreReqs() {
	cmd := "brew --version | grep \"Homebrew \" | awk '{print $2}'"
	out, err := t.command(cmd)
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