package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

const dummyAryaClusterName = "mytestcluster"

var _ = Describe("Cluster", func() {
	Describe("AryaCluster", func() {
		var c *dev.AryaCluster
		var cErr error
		BeforeEach(func() {
			cfg := dev.BaseAryaClusterCfg
			cfg.Name = dummyAryaClusterName
			cfg.NumNodes = 1
			cfg.Memory = 2
			cfg.Storage = 3
			c = dev.NewAryaCluster(cfg)
			c.Bind()
			if !c.Exists() {
				log.Warn("Cluster does not exist")
				cErr = c.Provision()
			}
		})
		Describe("Provisioning a new cluster", func() {
			It("Shouldn't encounter an error while provisioning the cluster", func() {
				Expect(cErr).To(BeNil())
			})
			It("Should provision the correct vms", func() {
				checkVMCfg := dev.VMConfig{Name: c.Nodes()[0].VM.Name()}
				checkVM := dev.NewVM(checkVMCfg)
				Expect(checkVM.Exists()).To(BeTrue())
			})
			It("Shouldn't provision any extraneous vms", func() {
				extVMCfg := dev.VMConfig{Name: dummyAryaClusterName + "2"}
				extVm := dev.NewVM(extVMCfg)
				Expect(extVm.Exists()).To(BeFalse())
			})
			It("Should provision the correct number of k3s clusters", func() {
				Expect(len(c.Nodes())).To(Equal(1))
			})
		})
		Describe("Checking if a cluster exists", func() {
			It("Should return true when the cluster exists", func() {
				cfg := dev.AryaClusterConfig{Name: dummyAryaClusterName}
				existingCluster := dev.NewAryaCluster(cfg)
				Expect(existingCluster.Exists()).To(BeTrue())
			})
			It("Should return false when the cluster doesn't exist", func() {
				cfg := dev.AryaClusterConfig{Name: "randomclustername12414"}
				nonExistentCluster := dev.NewAryaCluster(cfg)
				Expect(nonExistentCluster.Exists()).To(BeFalse())
			})
		})
		Describe("Binding to an existing cluster", func() {
			It("Should effectively bind to the correct number of nodes", func() {
				bindCluster := dev.NewAryaCluster(dev.AryaClusterConfig{Name: dummyAryaClusterName})
				bindCluster.Bind()
				Expect(len(bindCluster.Nodes())).To(Equal(1))
			})
		})
	})
})
