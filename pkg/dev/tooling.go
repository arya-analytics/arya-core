package dev

import (
	"github.com/arya-analytics/aryacore/pkg/util/emoji"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type Tools []string

// RequiredTools returns a slice of the required tools for provisioning dev clusters.
func RequiredTools() Tools {
	return Tools{
		"multipass",
		"kubernetes-cli",
		"krew",
		"yq",
		"helm",
		"gh",
	}

}

// || REQUIRED TOOL INSTALLS ||

// InstallRequiredTools installs tools required required for provisioning dev clusters.
func InstallRequiredTools() {
	log.Infof("%s Installing dev tools", emoji.Tools)
	t := NewTooling()
	for _, k := range RequiredTools() {
		if t.Installed(k) {
			log.Infof("%s %s already installed", emoji.Frog, k)
		} else {
			log.Infof("%s Installing %s", emoji.Flame, k)
			if err := t.Install(k); err != nil {
				log.Fatalln(err)
			}
			log.Infof("%s  Installed %s", emoji.Check, k)
		}
	}
	log.Infof("%s  All Done!", emoji.Check)
}

// UninstallRequiredTools uninstalls tools required for provisioning dev clusters.
func UninstallRequiredTools() {
	log.Infof("%s Uninstalling dev tools", emoji.Tools)
	t := NewTooling()
	for i := range RequiredTools() {
		k := RequiredTools()[len(RequiredTools())-1-i]
		if t.Installed(k) {
			log.Infof("%s Uninstalling %s", emoji.Flame, k)
			if err := t.Uninstall(k); err != nil {
				log.Fatalln(err)
			}
			log.Infof("%s  Uninstalled %s", emoji.Check, k)

		} else {
			log.Infof("%s %s is not installed. \n", emoji.Check, k)
		}
	}
	log.Infof("%s  All Done!", emoji.Check)
}

// RequiredToolsInstalled checks to see if all mandatory dev tools are installed.
func RequiredToolsInstalled() bool {
	log.Infof("%s Checking if required dev tools are installed", emoji.Tools)
	t := NewTooling()
	for _, v := range RequiredTools() {
		if !t.Installed(v) {
			log.Infof("%s Missing required dev tools\n", emoji.Frog)
			return false
		}
	}
	log.Infof("%s  All required dev tools installed\n", emoji.Check)
	return true
}

// || GENERAL TOOLING ||

// Tooling provides a generic interface for installing dev tools such as kubectl,
// multipass, yq, etc.
type Tooling interface {
	// Install installs a dev tool based on its name.
	Install(tool string) error
	// Uninstall uninstalls a dev tool based on its name.
	Uninstall(tool string) error
	// Installed checks if a package has already been installed.
	Installed(tool string) bool
}

// NewTooling creates and returns the correct OS specific tools manager.
func NewTooling() Tooling {
	t := BrewTooling{RequiredTools()}
	t.preReqs()
	return &t
}

// || BREW TOOLING ||

const brewCmd = "brew"

var requiredBrewVersion, _ = version.NewVersion("3.3.8")

type BrewTooling struct {
	tools Tools
}

func (t BrewTooling) Install(tool string) error {
	return t.command("install", tool).Run()
}

func (t BrewTooling) Uninstall(tool string) error {
	return t.command("uninstall", tool).Run()
}

func (t BrewTooling) Installed(tool string) bool {

	out, err := t.command("list").Output()
	if err != nil {
		panic(err)
	}
	outStr := string(out[:])
	return strings.Contains(outStr, tool)
}

func (t BrewTooling) command(args ...string) *exec.Cmd {
	return exec.Command(brewCmd, args...)
}

func (t BrewTooling) preReqs() {
	cmdString := "brew --version | grep \"Homebrew \" | awk '{print $2}'"
	// Need to manually create exec.Command here in order to use bash pipes
	o, err := exec.Command("bash", "-c", cmdString).Output()
	if err != nil {
		log.Fatalf("%s", err)
	}
	bv, err := version.NewVersion(strings.TrimSpace(string(o[:])))
	if err != nil {
		log.Fatalf("%s", err)
	}
	if bv.LessThan(requiredBrewVersion) {
		log.Fatalf("Brew Version %s is less than the required version %s",
			bv.String(), requiredBrewVersion.String())
	}
}
