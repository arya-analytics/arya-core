package git

import (
	"os/exec"
	"strings"
)

const gitCmd = "git"

// CurrentCommitHash returns the current commit hash of the repository.
func CurrentCommitHash() string {
	cmd := exec.Command(gitCmd, "log", "-1", "--format=\"%H\"")
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.Trim(strings.TrimSpace(string(o[:])), "\"")
}

// Username returns the git username registered on the host machine.
func Username() string {
	cmd := exec.Command(gitCmd,"config", "user.email")
	o, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(o[:]))
}