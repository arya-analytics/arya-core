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
		Describe("Provisioning a new cluster", func() {
			var c *dev.AryaCluster
			var cErr error
			BeforeEach(func() {
				cfg := dev.BaseAryaClusterCfg
				cfg.Name = dummyAryaClusterName
				cfg.NumNodes = 1
				cfg.Memory = 2
				cfg.Storage = 3
				log.Warn("Integrity check %s", dev.BaseAryaClusterCfg)
				c = dev.NewAryaCluster(cfg)
				c.Bind()
				if !c.Exists() {
					cErr = c.Provision()
				}
			})
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
			It("Should provision the correct number k3s clusters", func() {
				clusters := c.Nodes()
				Expect(len(clusters)).To(Equal(1))
			})
			It("It should provision the correct k3s clusters", func() {
				clusters := c.Nodes()
				Expect(clusters[0].Cfg.PodCidr).To(Equal("10.11.0.0/16"))
				Expect(clusters[0].Cfg.ServiceCidr).To(Equal("10.12.0.0/16"))
			})

		})
	})
})
