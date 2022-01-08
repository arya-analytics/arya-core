package kubectl

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

const kubectlCmd = "kubectl"

func Command(args ...string) *exec.Cmd {
	return exec.Command(kubectlCmd, args...)
}

// Exec executes a kubectl command.
func Exec(args ...string) error {
	return Command(args...).Run()
}

// CurrentContext returns the active kubectl config context.
func CurrentContext() (string, error) {
	o, err := Command("config", "view", "-o", "template",
		"--template='{{index . \"current-context\"}}").Output()
	return string(o[:]), err
}

// SwitchContext switches the kubectl context to the provided ctx string.
func SwitchContext(ctx string) error {
	err := Command("config", "use-context", ctx).Run()
	if err != nil {
		return err
	}
	log.Trace("Switched to kubecontext %s", ctx)
	return nil
}
