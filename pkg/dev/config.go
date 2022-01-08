package dev

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/kubectl"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// || ARYA CONFIG ||

const (
	k3sKubeCfgPath    = "/etc/rancher/k3s/k3s.yaml"
	authSecretName    = "regcred"
	kubeCfgBaseName   = "kubeconfig."
	aryaCfgType       = "kubernetes.io/dockerconfigjson"
	orchestratorLabel = "aryaRole=orchestrator"
)

var (
	hostKubeCfgPathBase = os.ExpandEnv("$HOME") + "/.kube/"
	aryaCfgPath         = os.ExpandEnv("$HOME") + "/.arya/config.json"
)

// AuthenticateCluster authenticate a k3s cluster by pulling auth credentials from
// the arya config.json file and creating an auth secret that resources can access.
func AuthenticateCluster(c K3sCluster) {
	info, err := c.VM.Info()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infof("Authenticating cluster on node %s", info.Name)
	if err := kubectl.SwitchContext(info.Name); err != nil {
		log.Fatalln(err)
	}

	// Ok to skip error check as will get caught on next command
	_ = kubectl.Exec("delete", "secret", authSecretName)

	if err := kubectl.Exec(
		"create",
		"secret",
		"generic",
		authSecretName,
		"--from-file=.dockerconfigjson="+aryaCfgPath,
		"--type="+aryaCfgType,
	); err != nil {
		log.Fatalln(err)
	}
}

// MergeClusterConfig pulls the kubeconfig file from cluster c,
// transfers it to the host machine, and merges it into the host kubeconfig.
func MergeClusterConfig(c K3sCluster) {
	ClearClusterConfig(c)
	vmInfo, err := c.VM.Info()
	if err != nil {
		log.Fatal(err)
	}
	name := kubeCfgBaseName + vmInfo.Name
	cfgPath := hostKubeCfgPathBase + name

	if err := transferKubeConfig(c.VM, cfgPath); err != nil {
		log.Fatal(err)
	}

	BindConfigToCluster(c, cfgPath)

	_ = kubectl.Exec("krew", "install", "konfig")

	if err := exec.Command("bash", "-c",
		fmt.Sprintf("kubectl konfig import -s %s", cfgPath)).Run(); err != nil {
		log.Fatal(err)
	}
}

// BindConfigToCluster binds a remote kubeconfig to the host machine by adding the
// remote IP and setting the correct context.
func BindConfigToCluster(c K3sCluster, cfgPath string) {
	n := c.VM.Name()
	info, err := c.VM.Info()
	if err != nil {
		log.Fatalln(err)
	}
	cmd := fmt.Sprintf(
		"yq -i eval '.clusters[].cluster.server |= sub(\"127.0.0.1\", \"%s\")"+
			" | .contexts[].name = \"%s\""+
			" | .current-context = \"%s\""+
			" | .clusters[].name = \"%s\""+
			" | .contexts[].context.cluster=\"%s\""+
			" | .users[].name = \"%s\""+
			" | .contexts[].context.user = \"%s\"' "+
			"%s",
		info.IPv4, n, n, n, n, n, n, cfgPath,
	)
	if err := exec.Command("bash", "-c", cmd).Run(); err != nil {
		log.Fatal(err)
	}
}

var clearCfgCmdChain = []string{"delete-cluster", "delete-user", "delete-context"}

// ClearClusterConfig clears a clusters kubeconfig information from the host
// kubeconfig.
func ClearClusterConfig(c K3sCluster) {
	name := c.VM.Name()
	for _, v := range clearCfgCmdChain {
		if err := kubectl.Exec("config", v, name); err != nil {
			log.Warnf("Unable to delete config for cluster %s", name)
		}
	}
}

// LabelOrchestrator labels a specific kubernetes node with the Arya role of
// orchestrator.
func LabelOrchestrator(nodeName string) {
	if err := kubectl.Exec("label", "nodes", nodeName, orchestratorLabel); err != nil {
		log.Fatalf("Failed to label orchestrator node %s", nodeName)
	}
}

func transferKubeConfig(vm VM, hostPath string) error {
	return vm.Transfer(TransferFrom, k3sKubeCfgPath, hostPath)
}
