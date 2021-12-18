package dev

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"log"
	"os"
	"os/exec"
	"strings"
)

// || REQUIRED TOOL INSTALLS ||

// InstallRequired installs all mandatory dev tools necessary for developing aryacore.
func InstallRequired() error {
	fmt.Printf("%s Installing dev tools \n", emoji("\\U0001F6E0"))
	t := NewTooling()
	for _, k := range RequiredTools {
		if t.CheckInstalled(k) {
			fmt.Printf("%s %s already installed \n", emoji("\\U0001F438"), k)
		} else {
			fmt.Printf("%s Installing %s \n", emoji("\\U0001f525"), k)
			if err := t.Install(k); err != nil {
				return err
			}
			fmt.Printf("%s  Installed %s \n", emoji("\\U0002705"), k)
		}
	}
	fmt.Printf("%s  All Done! \n", emoji("\\U0002705"))
	return nil
}

// UninstallRequired uninstalls all mandatory dev tools necessary for developing
//aryacore.
func UninstallRequired() error {
	fmt.Printf("%s Uninstalling dev tools \n", emoji("\\U0001F6E0"))
	t := NewTooling()
	for i := range RequiredTools {
		k := RequiredTools[len(RequiredTools)-1-i]
		if t.CheckInstalled(k) {
			fmt.Printf("%s Uninstalling %s \n", emoji("\\U0001f525"), k)
			if err := t.Uninstall(k); err != nil {
				return err
			}
			fmt.Printf("%s  Uninstalled %s \n", emoji("\\U0002705"), k)

		} else {
			fmt.Printf("%s %s is not installed. \n", emoji("\\U0001F438"), k)
		}
	}
	fmt.Printf("%s  All Done! \n", emoji("\\U0002705"))
	return nil
}

// || GENERAL TOOLING ||

type Tooling interface {
	Install(tool string) error
	Uninstall(tool string) error
	CheckInstalled(tool string) bool
}

// NewTooling creates and returns the correct OS specific tooling manager.
func NewTooling() Tooling {
	t := BrewTooling{RequiredTools}
	t.checkPreReqs()
	return &t
}

// || BREW TOOLING ||

var requiredBrewVersion, _ = version.NewVersion("3.3.8")

type BrewTooling struct {
	tooling ToolingConfig
}

// Install installs a dev tool based on its name.
func (t BrewTooling) Install(tool string) error {
	_, err := t.command("install " + tool)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return nil
}

// Uninstall uninstalls a dev tool based on its name.
func (t BrewTooling) Uninstall(tool string) error {
	_, err := t.command("uninstall " + tool)
	return err
}

// CheckInstalled checks if a package has already been installed.
func (t BrewTooling) CheckInstalled(tool string) bool {
	out, err := t.command("list")
	if err != nil {
		panic(err)
	}
	outStr := string(out[:])
	return strings.Contains(outStr, tool)
}

/// command wraps exec.command to add brew specific functionality.
func (t BrewTooling) command(cmdString string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", "brew "+cmdString)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// checkPreReqs checks for necessary pre-requisites for the installer to run correctly
func (t BrewTooling) checkPreReqs() {
	cmd := "--version | grep \"Homebrew \" | awk '{print $2}'"
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
