package kubectl

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
)

const kubectlCmd = "kubectl"

func Command(args ...string) *exec.Cmd {
	return exec.Command(kubectlCmd, args...)
}

func Exec(args ...string) error {
	return Command(args...).Run()
}

func CurrentContext() (string, error) {
	o, err := Command("config","view","-o","template",
		"--template='{{index . \"current-context\"}}").Output()
	return string(o[:]), err
}

func SwitchContext(ctx string) error {
	err := Command("config", "use-context", ctx).Run()
	if err != nil {
		return err
	}
	log.Infof("Successfully switched to kubectl context %s", ctx)
	return nil
}