package dev

import "fmt"

type ClusterConfig struct {
	PodCidr     string
	ServiceCidr string
}

const (
	k3sAddr             = "https://get.k3s.io"
	writeKubeConfigMode = "777"
)

type K3sCluster struct {
	vm  VM
	cfg ClusterConfig
}

func NewK3sCluster(vm VM, cfg ClusterConfig) *K3sCluster {
	return &K3sCluster{vm, cfg}
}

func (p K3sCluster) Provision() error {
	curlCmd := fmt.Sprintf("curl -sfL %s", k3sAddr)
	k3sEnv := fmt.Sprintf("INSTALL_K3S_EXEC=\"--cluster-cidr %s --service-cidr %s"+
		" --write-kubeconfig-mode %s \"",
		p.cfg.PodCidr, p.cfg.ServiceCidr, writeKubeConfigMode)
	k3sInitCmd := "sh -s -"
	cmdString := fmt.Sprintf("%s | %s %s", curlCmd, k3sEnv, k3sInitCmd)
	_, err := p.vm.Exec(cmdString)
	return err
}
