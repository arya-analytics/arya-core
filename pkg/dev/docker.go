package dev

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ImageCfg struct {
	Repository string
	Tag string
	BuildCtxPath string
}

type DockerImage struct {
	cfg ImageCfg
}

func NewDockerImage(cfg ImageCfg) *DockerImage {
	di := DockerImage{
		cfg: cfg,
	}
	return &di
}

const dockerCmd = "docker"

func (d DockerImage) command(args ...string) *exec.Cmd {
	cmd := exec.Command(dockerCmd, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd
}

func (d DockerImage) NameTag() string {
	return createNameTag(d.cfg.Repository, d.cfg.Tag)
}

func (d DockerImage) Build() error {
	return d.command("build", d.cfg.BuildCtxPath, "-t", d.NameTag()).Run()
}

func (d DockerImage) Push() error {
	return d.command("push", d.NameTag()).Run()
}

const imageNameTagSeparator = ":"
const nameTagFormat = "name:tag"

func parseNameTag(nameTag string) (repository  string, tag string, err error) {
	split := strings.Split(nameTag, imageNameTagSeparator)
	if len(split) != 2 {
		return "", "", fmt.Errorf("nameTag %s has the incorrect format. " +
			"The correct format is %s", nameTag, nameTagFormat)
	}
	return split[0], split[1], nil
}

func createNameTag(repository, tag string) string {
	return repository + imageNameTagSeparator + tag
}