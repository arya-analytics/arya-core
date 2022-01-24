package dev

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/emoji"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// || LOCAL DEV CLUSTER ||

// ProvisionLocalDevCluster provisions a local development Cluster based on the
// specified parameters.
func ProvisionLocalDevCluster(numNodes int, name string, cores int, memory int,
	storage int, reInit bool, cidrOffset int) (*AryaCluster, error) {
	var cluster *AryaCluster
	InstallRequiredTools()
	log.Infof("%s Provisioning an Arya Cluster named %s with %v nodes",
		emoji.Bolt, name, numNodes)
	cfg := AryaClusterConfig{
		Name:       name,
		NumNodes:   numNodes,
		Cores:      cores,
		Memory:     memory,
		Storage:    storage,
		ReInit:     reInit,
		CidrOffset: cidrOffset,
	}
	cluster = NewAryaCluster(cfg)
	err := cluster.Provision()
	if err != nil {
		log.Fatalln(err)
	}
	for i, c := range cluster.Nodes() {
		nodeName := c.VM.Name()
		log.Infof("Merging kubeconfig for node %s", nodeName)
		MergeClusterConfig(*c)
		log.Infof("Authenticating Cluster %s", nodeName)
		AuthenticateCluster(*c)
		if i == 0 {
			log.Infof("Marking node %s as the Cluster orchestrator", nodeName)
			LabelOrchestrator(nodeName)
		}
	}
	log.Infof("%s  Successfully initialized Arya Cluster %s", emoji.Check, name)
	return cluster, nil
}

// DeleteLocalDevCluster deletes a local development Cluster based on its Name.
func DeleteLocalDevCluster(name string) error {
	log.Infof("%s Deleting Arya Cluster %s", emoji.Flame, name)
	cfg := AryaClusterConfig{Name: name}
	c := NewAryaCluster(cfg)
	c.Bind()
	return c.Delete()
}

// || ARYA CLUSTER ||
// Utilities for provisioning and managing development Arya Clusters.

type AryaClusterConfig struct {
	Name       string
	NumNodes   int
	Cores      int
	Memory     int
	Storage    int
	ReInit     bool
	CidrOffset int
}

var BaseAryaClusterCfg = AryaClusterConfig{
	Name:       "ad",
	NumNodes:   3,
	Cores:      2,
	Memory:     4,
	Storage:    15,
	ReInit:     true,
	CidrOffset: 10,
}

type AryaCluster struct {
	cfg   AryaClusterConfig
	nodes []*K3sCluster
}

// NewAryaCluster creates a new Arya Cluster based off of a config.
// NOTE: For binding to an existing Cluster, only Cfg.Name is necessary.
func NewAryaCluster(cfg AryaClusterConfig) *AryaCluster {
	return &AryaCluster{cfg: cfg}
}

// Provision provisions a new Cluster base off of a.Cfg.
func (a *AryaCluster) Provision() error {
	for i := 1; i <= a.cfg.NumNodes; i++ {
		nodeName := a.cfg.Name + strconv.Itoa(i)
		log.Infof("Bootstrapping node %v  with name %s", i, nodeName)
		log.Infof("Memory: %v Gb, Cores: %v, Storage: %v Gb", a.cfg.Memory,
			a.cfg.Cores,
			a.cfg.Storage)
		vm, err := a.provisionVM(nodeName)
		if err != nil {
			return err
		}
		podCidrID := a.cfg.CidrOffset + i*2
		k3s, err := a.provisionK3s(vm, podCidrID)
		if err != nil {
			return err
		}
		a.nodes = append(a.nodes, k3s)
	}
	return nil
}

// provisionK3s provisions a K3s Cluster on a VM
// Needs to receive a pod cidr ID (ex. 44 would result in the call Cidr(
// 44) for the pod Cidr and Cidr(45) for the service Cidr).
func (a *AryaCluster) provisionK3s(vm VM, podCidrID int) (*K3sCluster, error) {
	log.Infof("Provisioning new K3s Cluster on VM %s", vm.Name())
	cfg := K3sClusterConfig{
		PodCidr:     Cidr(podCidrID),
		ServiceCidr: Cidr(podCidrID + 1),
	}
	c := NewK3sCluster(vm, cfg)
	if err := c.Provision(); err != nil {
		return c, err
	}
	log.Infof("Succesffully provisioned k3s Cluster on VM %s", vm.Name())
	return c, nil
}

// provisionVM provisions a virtual machine for the Cluster based off a node Name
// and internal config information.
// NOTE: If a.Cfg.reInit is set to true, will tear down existing VM's.
func (a *AryaCluster) provisionVM(nodeName string) (VM, error) {
	cfg := VMConfig{
		Name:    nodeName,
		Memory:  a.cfg.Memory,
		Cores:   a.cfg.Cores,
		Storage: a.cfg.Storage,
	}
	vm := NewVM(cfg)
	if !vm.Exists() {
		log.Infof("Provisioning new VM named %s ", nodeName)
		if err := vm.Provision(); err != nil {
			return vm, err
		}
	} else if a.cfg.ReInit {
		log.Infof("VM %s already existed, tearing down and re-launching",
			nodeName)
		if err := vm.Delete(); err != nil {
			return vm, err
		}
		if err := vm.Provision(); err != nil {
			return vm, err
		}
	} else {
		return vm, fmt.Errorf("VM %s already exists and ReInit is false", nodeName)
	}
	log.Infof("Successfully provisioned VM %s", nodeName)
	return vm, nil
}

// Bind binds to an existing arya Cluster based on its Name.
func (a *AryaCluster) Bind() {
	for i := 1; i > 0; i++ {
		cfg := VMConfig{
			Name: a.cfg.Name + strconv.Itoa(i),
		}
		vm := NewVM(cfg)
		if !vm.Exists() {
			break
		}
		cluster := NewK3sCluster(vm, K3sClusterConfig{})
		a.nodes = append(a.nodes, cluster)
	}
}

// Nodes returns the nodes in the Cluster.
func (a *AryaCluster) Nodes() []*K3sCluster {
	return a.nodes
}

// Exists checks if a Cluster with a.Cfg.Name already exists.
func (a *AryaCluster) Exists() bool {
	if len(a.Nodes()) > 0 {
		return true
	}
	a.Bind()
	if len(a.Nodes()) > 0 {
		return true
	}
	return false
}

// Delete Deletes an existing Cluster and purges all of its data.
func (a *AryaCluster) Delete() error {
	for _, node := range a.Nodes() {
		if err := node.VM.Delete(); err != nil {
			return err
		}
		ClearClusterConfig(*node)
	}
	return nil
}

// || K3S CLUSTER ||
// Utilities for provisioning k3S (https://k3s.io/) clusters on VM's.

const (
	k3sAddr             = "https://get.k3s.io"
	writeKubeConfigMode = "777"
)

type K3sClusterConfig struct {
	PodCidr     string
	ServiceCidr string
}

type K3sCluster struct {
	VM  VM
	Cfg K3sClusterConfig
}

// NewK3sCluster creates a new k3s Cluster.
func NewK3sCluster(vm VM, cfg K3sClusterConfig) *K3sCluster {
	return &K3sCluster{vm, cfg}
}

// Provision provisions a new k3s Cluster on p.VM.
func (p K3sCluster) Provision() error {
	curlCmd := fmt.Sprintf("curl -sfL %s", k3sAddr)
	k3sEnv := fmt.Sprintf("INSTALL_K3S_EXEC=\"--cluster-cidr %s --service-cidr %s"+
		" --write-kubeconfig-mode %s \"",
		p.Cfg.PodCidr, p.Cfg.ServiceCidr, writeKubeConfigMode)
	k3sInitCmd := "sh -s -"
	cmdString := fmt.Sprintf("%s | %s %s", curlCmd, k3sEnv, k3sInitCmd)
	_, err := p.VM.Exec(cmdString)
	return err
}

// || CLUSTER NETWORKING ||
// Utilities for networking inside of and between clusters.

const (
	baseCidrPrefix = "10."
	baseCidrSuffix = ".0.0/16"
)

// Cidr is a utility for generating kubernetes resource IP ranges.
// Generates an IPv4 address from a two digit unique ID (00-99).
func Cidr(ID int) string {
	return baseCidrPrefix + strconv.Itoa(ID) + baseCidrSuffix
}
