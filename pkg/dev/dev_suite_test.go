package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDev(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dev Suite")
}

const testTool = "aamath"

var tooling dev.Tooling

var vm dev.VM
var vmInfo dev.VMInfo

var vmCfg = dev.VMConfig{
	Name:    "adtst1",
	Memory:  2,
	Cores:   3,
	Storage: 4,
}

func provisionDummyAryaClusterIfNotExists() (*dev.AryaCluster, error) {
	cfg := dev.BaseAryaClusterCfg
	cfg.Name = dummyAryaClusterName
	cfg.NumNodes = 1
	cfg.Memory = 2
	cfg.Storage = 3
	c := dev.NewAryaCluster(cfg)
	c.Bind()
	var cErr error
	if !c.Exists() {
		cErr = c.Provision()
		for _, c := range c.Nodes() {
			dev.MergeClusterConfig(*c)
			dev.AuthenticateCluster(*c)
		}
	}
	return c, cErr
}

var _ = BeforeSuite(func() {
	tooling = dev.NewTooling()
	if !tooling.Installed(testTool) {
		if err := tooling.Install(testTool); err != nil {
			log.Fatalln("Unable to uninstall test tool")
		}
	}

	vm = dev.NewVM(vmCfg)
	if vm.Exists() {
		if err := vm.Delete(); err != nil {
			log.Fatalln("Failed to delete test VM")
		}
	}
	if err := vm.Provision(); err != nil {
		log.Fatalln("Failed to provision test VM")
	}
	vi, err := vm.Info()
	vmInfo = vi
	if err != nil {
		log.Fatalln("Failed to pull info from test VM")
	}
})

var _ = AfterSuite(func() {
	vm = dev.NewVM(vmCfg)
	if vm.Exists() {
		if err := vm.Delete(); err != nil {
			log.Fatalln("Failed to delete test VM")
		}
	}
	aryaCluster := dev.NewAryaCluster(dev.AryaClusterConfig{Name: dummyAryaClusterName})
	aryaCluster.Bind()
	if err := aryaCluster.Delete(); err != nil {
		log.Fatalln(err)
	}
	if err := dev.InstallRequiredTools(); err != nil {
		log.Fatalln(err)
	}
})
