package dev

import (
	"fmt"
	"strconv"
)

const baseNodeName = "ad"

// || ARYA CLUSTER ||

type AryaClusterConfig struct {
	Name     string
	NumNodes int
	Cores    int
	Memory   int
	Storage  int
	ReInit   bool
	CidrOffset int
}

var BaseAryaClusterCfg = AryaClusterConfig{
	Name:     "ad",
	NumNodes: 3,
	Cores:    2,
	Memory:   4,
	Storage:  15,
	ReInit:   true,
	CidrOffset: 10,
}

type AryaCluster struct {
	cfg   AryaClusterConfig
	nodes []*K3sCluster
}

func NewAryaCluster(cfg AryaClusterConfig) *AryaCluster {
	return &AryaCluster{cfg: cfg}
}

func (a *AryaCluster) Provision() ([]*K3sCluster, error) {
	var k3sClusters []*K3sCluster
	for i := 1; i <= a.cfg.NumNodes; i++ {
		nodeName := baseNodeName + strconv.Itoa(i)
		fmt.Printf("Bootstrapping node %v  with name %s \n", i, nodeName)
		fmt.Printf("Memory: %v, Cores: %v, Storage: %v \n", a.cfg.Memory, a.cfg.Cores,
			a.cfg.Storage)
		vm, err := a.ProvisionVM(nodeName)
		if err != nil {
			return k3sClusters, err
		}
		podCidrNum := a.cfg.CidrOffset + i * 2
		k3s, err := a.ProvisionK3s(vm, podCidrNum)
		if err != nil {
			return k3sClusters, err
		}
		k3sClusters = append(k3sClusters, k3s)
		fmt.Printf("Successfully started k3s cluster on node %s \n", nodeName)
	}
	return k3sClusters, nil
}

func (a *AryaCluster) ProvisionK3s(vm VM, podCidrNum int) (*K3sCluster, error) {
	cfg := K3sClusterConfig{
		PodCidr:     Cidr(podCidrNum),
		ServiceCidr: Cidr(podCidrNum + 1),
	}
	c := NewK3sCluster(vm, cfg)
	if err := c.Provision(); err != nil {
		return c, err
	}
	return c, nil
}

func (a *AryaCluster) ProvisionVM(nodeName string) (VM, error) {
	cfg := VMConfig{
		Name:    nodeName,
		Memory:  a.cfg.Memory,
		Cores:   a.cfg.Cores,
		Storage: a.cfg.Storage,
	}	
	vm := NewVM(cfg)
	if !vm.Exists() {
		fmt.Printf("Launching new VM named %s \n ", nodeName)
		if err := vm.Provision(); err != nil {
			return vm, err
		}
	} else if a.cfg.ReInit {
		fmt.Printf("VM %s already existed, tearing down and re-launching \n",
			nodeName)
		if err := vm.Delete(); err != nil {
			return vm, err
		}
		if err := vm.Provision(); err != nil {
			return vm, err
		}
	} else {
		return vm, fmt.Errorf("vm %s already exists and ReInit is false \n", nodeName)
	}
	return vm, nil
}

// || K3S CLUSTER ||

const (
	k3sAddr             = "https://get.k3s.io"
	writeKubeConfigMode = "777"
)

type K3sClusterConfig struct {
	PodCidr     string
	ServiceCidr string
}

type K3sCluster struct {
	vm  VM
	cfg K3sClusterConfig
}

func NewK3sCluster(vm VM, cfg K3sClusterConfig) *K3sCluster {
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

// || CLUSTER NETWORKING ||

const baseCidrPrefix = "10."
const baseCidrSuffix = ".0.0/16"

func Cidr(ip int) string {
	return baseCidrPrefix + strconv.Itoa(ip) + baseCidrSuffix
}