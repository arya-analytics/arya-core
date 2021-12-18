package dev

import (
	"fmt"
	"strconv"
)

type ClusterConfig struct {
	PodCidr     string
	ServiceCidr string
}

const baseNodeName = "ad"

const baseCidrPrefix = "10."
const baseCidrSuffix = ".0.0/16"

func Cidr(ip int) string {
	return baseCidrPrefix + strconv.Itoa(ip) + baseCidrSuffix
}

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
	return &AryaCluster{cfg: cfg, nodes: []*K3sCluster{}}
}

func (a *AryaCluster) Provision() error {
	for i := 1; i <= a.cfg.NumNodes; i++ {
		nodeName := baseNodeName + strconv.Itoa(i)
		fmt.Printf("Bootstrapping node %v  with name %s \n", i, nodeName)
		fmt.Printf("Memory: %v, Cores: %v, Storage: %v \n", a.cfg.Memory, a.cfg.Cores,
			a.cfg.Storage)
		vmCfg := VMConfig{
			Name:    nodeName,
			Memory:  a.cfg.Memory,
			Cores:   a.cfg.Cores,
			Storage: a.cfg.Storage,
		}
		vm := NewVM(vmCfg)
		if !vm.Exists() {
			fmt.Printf("Launching new VM named %s \n ", nodeName)
			if err := vm.Launch(); err != nil {
				return err
			}
		} else if a.cfg.ReInit {
			fmt.Printf("VM %s already existed, tearing down and re-launching \n",
				nodeName)
			if err := vm.Delete(); err != nil {
				return err
			}
			if err := vm.Launch(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("vm %s already exists and ReInit is false \n", nodeName)
		}
		podCidrNum := a.cfg.CidrOffset + i * 2
		k3sCfg := ClusterConfig{
			PodCidr:     Cidr(podCidrNum),
			ServiceCidr: Cidr(podCidrNum + 1),
		}
		fmt.Printf("Assigning pod cidr %s and service cidr %s to node %s \n",
			k3sCfg.PodCidr, k3sCfg.ServiceCidr, nodeName)
		k3sCluster := NewK3sCluster(vm, k3sCfg)
		if err := k3sCluster.Provision(); err != nil {
			return err
		}
		fmt.Printf("Successfully started k3s cluster on node %s \n", nodeName)
		a.nodes = append(a.nodes, k3sCluster)
	}
	return nil
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
