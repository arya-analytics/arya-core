package git

import (
	"os/exec"
	"strings"
)

func CurrentCommitHash() string {
	cmd := exec.Command("bash", "-c", "git log -1 --format=\"%H\"")
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(o[:]))
}

func Username() string {
	cmd := exec.Command("bash", "-c", "git config user.email")
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(o[:]))
}