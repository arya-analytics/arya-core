package dev

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/emoji"
	"github.com/hashicorp/go-version"
	"log"
	"os/exec"
	"strings"
)

// || REQUIRED TOOL INSTALLS ||

// InstallRequired installs all mandatory dev tools necessary for developing aryacore.
// Takes a verbosity flag
func InstallRequired() error {
	fmt.Printf("%s Installing dev tools \n", emoji.Tools)
	t := NewTooling()
	for _, k := range RequiredTools {
		if t.Installed(k) {
			fmt.Printf("%s %s already installed \n", emoji.Frog, k)
		} else {
			fmt.Printf("%s Installing %s \n", emoji.Flame, k)
			if err := t.Install(k); err != nil {
				return err
			}
			fmt.Printf("%s  Installed %s \n", emoji.Check, k)
		}
	}
	fmt.Printf("%s  All Done! \n", emoji.Check)
	return nil
}

// UninstallRequired uninstalls all mandatory dev tools necessary for developing
//aryacore. Takes a verbosity flag.
func UninstallRequired() error {
	fmt.Printf("%s Uninstalling dev tools \n", emoji.Tools)
	t := NewTooling()
	for i := range RequiredTools {
		k := RequiredTools[len(RequiredTools)-1-i]
		if t.Installed(k) {
			fmt.Printf("%s Uninstalling %s \n", emoji.Flame, k)
			if err := t.Uninstall(k); err != nil {
				return err
			}
			fmt.Printf("%s  Uninstalled %s \n", emoji.Check, k)

		} else {
			fmt.Printf("%s %s is not installed. \n", emoji.Check, k)
		}
	}
	fmt.Printf("%s  All Done! \n", emoji.Check)
	return nil
}

// RequiredInstalled checks to see if all mandatory dev tools are installed.
func RequiredInstalled() bool {
	fmt.Printf("%s Checking if required dev tools are installed\n", emoji.Tools)
	t := NewTooling()
	for _, v := range RequiredTools {
		if !t.Installed(v) {
			fmt.Printf("%s Missing required dev tools\n", emoji.Frog)
			return false
		}
	}
	fmt.Printf("%s  All required dev tools installed\n", emoji.Check)
	return true
}

// || GENERAL TOOLING ||

type Tooling interface {
	Install(tool string) error
	Uninstall(tool string) error
	Installed(tool string) bool
}

// NewTooling creates and returns the correct OS specific tooling manager.
func NewTooling() Tooling {
	t := BrewTooling{RequiredTools}
	t.checkPreReqs()
	return &t
}

// || BREW TOOLING ||

const brewCommand = "brew"

var requiredBrewVersion, _ = version.NewVersion("3.3.8")

type BrewTooling struct {
	tooling ToolingConfig
}

// Install installs a dev tool based on its name.
func (t BrewTooling) Install(tool string) error {
	_, err := t.command("install", tool)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return nil
}

// Uninstall uninstalls a dev tool based on its name.
func (t BrewTooling) Uninstall(tool string) error {
	_, err := t.command("uninstall", tool)
	return err
}

// Installed checks if a package has already been installed.
func (t BrewTooling) Installed(tool string) bool {
	out, err := t.command("list")
	if err != nil {
		panic(err)
	}
	outStr := string(out[:])
	return strings.Contains(outStr, tool)
}

/// command wraps exec.command to add brew specific functionality.
func (t BrewTooling) command(args ...string) ([]byte, error) {
	cmd := exec.Command(brewCommand, args...)
	return cmd.Output()
}

// checkPreReqs checks for necessary pre-requisites for the installer to run correctly
func (t BrewTooling) checkPreReqs() {
	cmdString := "brew --version | grep \"Homebrew \" | awk '{print $2}'"
	// Need to manually create exec.Command here in order to use bash pipes
	out, err := exec.Command("bash", "-c", cmdString).Output()
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
