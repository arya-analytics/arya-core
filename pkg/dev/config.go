package dev

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// || KUBECTL ||

type KubeCtl int

const kubectlCmd = "kubectl"

func NewKubeCtl() *KubeCtl {
	var k KubeCtl
	return &k
}

func (k *KubeCtl) Command(args ...string) *exec.Cmd {
	return exec.Command(kubectlCmd, args...)
}

func (k *KubeCtl) Exec(args ...string) error {
	return k.Command(args...).Run()
}

func (k *KubeCtl) CurrentContext() (string, error) {
	o, err := k.Command("config","view","-o","template",
		"--template='{{index . \"current-context\"}}").Output()
	return string(o[:]), err
}

func (k *KubeCtl) SwitchContext(ctx string) error {
	err := k.Command("config", "use-context", ctx).Run()
	if err != nil {
		return err
	}
	log.Infof("Successfully switched to kubectl context %s", ctx)
	return nil
}

// || ARYA CONFIG ||

type AryaConfig struct {
	cfgPath string
}

const k3sKubeCfgPath = "/etc/rancher/k3s/k3s.yaml"

var hostKubeCfgPathBase = os.ExpandEnv("$HOME") + "/.kube/"

var aryaCfgPath = os.ExpandEnv("$HOME") + "/.arya/config.json"

var aryaCfgType = "kubernetes.io/dockerconfigjson"

const authSecretName = "regcred"

func NewAryaConfig(cfgPath string) *AryaConfig {
	return &AryaConfig{cfgPath: cfgPath}
}

func (a AryaConfig) AuthenticateCluster(c K3sCluster) error {
	info, err := c.VM.Info()
	if err != nil {
		return err
	}
	kubeCtl := NewKubeCtl()
	if err := kubeCtl.SwitchContext(info.Name); err != nil {
		return err
	}

	if err := kubeCtl.Exec("delete", "secret", authSecretName); err != nil {
		fmt.Println("We're fine here")
	}

	if err := kubeCtl.Exec(
		"create",
		"secret",
		"generic",
		authSecretName,
		"--from-file=.dockerconfigjson="+a.cfgPath,
		"--type="+aryaCfgType,
	); err != nil {
		return err
	}
	return nil
}

func (a AryaConfig) MergeClusterConfig(c K3sCluster) error {
	vmInfo, err := c.VM.Info()
	if err != nil {
		return err
	}
	name := "kubeconfig." + vmInfo.Name
	hostPath := hostKubeCfgPathBase + name
	fmt.Printf("Copying kubeconfig from %s to host path %s", name, hostPath)
	if err := c.VM.Transfer(TransferFrom, k3sKubeCfgPath,
		hostPath); err != nil {
		return err
	}
	fmt.Printf("Modifying kubeconfig to bind to correct IP")
	cmd := fmt.Sprintf(
		"yq -i eval '.clusters[].cluster.server |= sub(\"127.0.0.1\", \"%s\")"+
			" | .contexts[].name = \"%s\""+
			" | .current-context = \"%s\""+
			" | .clusters[].name = \"%s\""+
			" | .contexts[].context.cluster=\"%s\""+
			" | .users[].name = \"%s\""+
			" | .contexts[].context.user = \"%s\"' "+
			"%s",
		vmInfo.IPv4, vmInfo.Name, vmInfo.Name, vmInfo.Name, vmInfo.Name,
		vmInfo.Name, vmInfo.Name, hostPath,
	)
	if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("bash", "-c", "kubectl krew install konfig").Run(); err != nil {
		fmt.Println("We're fine here 2")
	}

	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl config delete-cluster %s", vmInfo.Name)).Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl config delete-context %s", vmInfo.Name)).Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl config delete-user %s", vmInfo.Name)).Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl konfig import -s %s", hostPath)).Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (a AryaConfig) LabelOrchestrator(nodeName string) error {
	k := NewKubeCtl()
	return k.Exec("label","nodes",nodeName, "aryaRole=orchestrator")
}