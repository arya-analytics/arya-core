package dev

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	dockerCmd        = "docker"
	nameTagSeparator = ":"
	nameTagFormat    = "name:tag"
)

type ImageCfg struct {
	Repository   string
	Tag          string
	BuildCtxPath string
}

type DockerImage struct {
	cfg ImageCfg
}

// NewDockerImage creates a new docker image based off of the provided config
func NewDockerImage(cfg ImageCfg) *DockerImage {
	di := DockerImage{
		cfg: cfg,
	}
	return &di
}

func (d DockerImage) command(args ...string) *exec.Cmd {
	cmd := exec.Command(dockerCmd, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd
}

// NameTag returns the name tag for the docker image.
func (d DockerImage) NameTag() string {
	return createNameTag(d.cfg.Repository, d.cfg.Tag)
}

// Build builds and tags the docker image based on d.ImageCfg.
func (d DockerImage) Build() error {
	return d.command("build", d.cfg.BuildCtxPath, "-t", d.NameTag()).Run()
}


// Push pushes the docker image to the locally authenticated repository.
func (d DockerImage) Push() error {
	return d.command("push", d.NameTag()).Run()
}

// parseNameTag parses the given name tag for the repository and string parts.
func parseNameTag(nameTag string) (repository string, tag string, err error) {
	split := strings.Split(nameTag, nameTagSeparator)
	if len(split) != 2 {
		return "", "", fmt.Errorf("nameTag %s has the incorrect format. "+
			"The correct format is %s", nameTag, nameTagFormat)
	}
	return split[0], split[1], nil
}

// createNameTag creates the given name tag based on the repository and tag
func createNameTag(repository, tag string) string {
	return repository + nameTagSeparator + tag
}
