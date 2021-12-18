
package dev

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type AryaConfig struct {
	kubeConfigs []string
}

const k3sKubeCfgPath = "/etc/rancher/k3s/k3s.yaml"

var hostKubeCfgPathBase = os.ExpandEnv("$HOME") + "/.kube/"

func NewAryaConfig() *AryaConfig {
	return &AryaConfig{}
}

func (a AryaConfig) MergeRemoteKubeConfig(c K3sCluster) error {
	vmInfo, err := c.vm.Info()
	if err != nil {
		return err
	}
	name := "kubeconfig." + vmInfo.Name
	hostPath := hostKubeCfgPathBase + name
	fmt.Printf("Copying kubeconfig from %s to host path %s", name, hostPath)
	if err := c.vm.Transfer(TransferFrom, k3sKubeCfgPath,
		hostPath); err != nil {
		return err
	}
	fmt.Printf("Modifying kubeconfig to bind to correct IP")
	cmd := fmt.Sprintf(
		"yq -i eval '.clusters[].cluster.server |= sub(\"127.0.0.1\", \"%s\")" +
			" | .contexts[].name = \"%s\"" +
			" | .current-context = \"%s\"" +
			" | .clusters[].name = \"%s\"" +
			" | .contexts[].context.cluster=\"%s\"" +
			" | .users[].name = \"%s\"" +
			" | .contexts[].context.user = \"%s\"' " +
			"%s",
			vmInfo.IPv4, vmInfo.Name, vmInfo.Name, vmInfo.Name, vmInfo.Name,
			vmInfo.Name, vmInfo.Name, hostPath,
	)
	if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c", "kubectl krew install konfig").Run(
		); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl config delete-cluster %s", vmInfo.Name)).Run(
	); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl config delete-context %s", vmInfo.Name)).Run(
	); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl config delete-user %s", vmInfo.Name)).Run(
	); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl konfig import -s %s", hostPath)).Run(
	); err != nil {
		log.Fatal(err)
	}


	a.kubeConfigs = append(a.kubeConfigs, name)
	return nil
}
