package git

import (
	"os/exec"
	"strings"
)

const gitCmd = "git"

func CurrentCommitHash() string {
	cmd := exec.Command(gitCmd, "log", "-1", "--format=\"%H\"")
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(o[:]))
}

func Username() string {
	cmd := exec.Command(gitCmd,"git","config", "user.email")
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(o[:]))
}