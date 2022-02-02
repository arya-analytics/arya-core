package dev

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	dockerComposeRelPath = "./deploy/docker/docker-compose.yaml"
	dockerBaseCmd        = "docker"
	nameTagSeparator     = ":"
	nameTagFormat        = "Name:tag"
)

type ImageCfg struct {
	Repository   string
	Tag          string
	BuildCtxPath string
}

type DockerImage struct {
	cfg ImageCfg
}

// NewDockerImage creates a new docker image based off of the provided config.
func NewDockerImage(cfg ImageCfg) *DockerImage {
	di := DockerImage{
		cfg: cfg,
	}
	return &di
}

func dockerCommand(args ...string) *exec.Cmd {
	cmd := exec.Command(dockerBaseCmd, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd
}

// NameTag returns the Name tag for the docker image.
func (d DockerImage) NameTag() string {
	return createNameTag(d.cfg.Repository, d.cfg.Tag)
}

// Build builds and tags the docker image based on d.ImageCfg.
func (d DockerImage) Build() error {
	return dockerCommand("build", d.cfg.BuildCtxPath, "-t", d.NameTag()).Run()
}

// Push pushes the docker image to the locally authenticated repository.
func (d DockerImage) Push() error {
	return dockerCommand("push", d.NameTag()).Run()
}

// createNameTag creates the given Name tag based on the repository and tag.
func createNameTag(repository, tag string) string {
	return repository + nameTagSeparator + tag
}

func StartDockerCompose() error {
	cmd := dockerCommand("compose", "-f", dockerComposeRelPath, "up")
	fmt.Printf(cmd.String())
	return cmd.Run()
}
