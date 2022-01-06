package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = BeforeSuite(func() {
	log.Info("Bootstrapping test suite")
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
	log.Info("Test suite bootstrapped successfully")
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