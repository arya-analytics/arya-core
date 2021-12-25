package dev

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/kubectl"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// || KUBECTL ||


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
	if err := kubectl.SwitchContext(info.Name); err != nil {
		return err
	}

	if err := kubectl.Exec("delete", "secret", authSecretName); err != nil {
		fmt.Println("We're fine here")
	}

	if err := kubectl.Exec(
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
	if err := a.ClearClusterConfig(c); err != nil {
		log.Fatalln(err)
	}
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

	if err := kubectl.Exec("krew", "install", "konfig"); err != nil {
		log.Warn("krew plugin already installed. Skipping reinstallation.")
	}

	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl konfig import -s %s", hostPath)).Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

var clearCfgCmdChain = []string{"delete-cluster","delete-user","delete-context"}

func (a AryaConfig) ClearClusterConfig(c K3sCluster) error {
	for _, v := range clearCfgCmdChain {
		if err := kubectl.Exec("config", v, c.VM.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (a AryaConfig) LabelOrchestrator(nodeName string) error {
	return kubectl.Exec("label","nodes",nodeName, "aryaRole=orchestrator")
}