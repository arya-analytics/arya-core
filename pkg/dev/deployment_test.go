package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo/v2"	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var dummyCfg = dev.DeploymentConfig{
	Name:      "mydummydeployment",
	ChartPath: "../../" + dev.DefaultChartRelPath,
	ImageCfg: dev.ImageCfg{Repository: dev.DefaultImageRepo,
		BuildCtxPath: ".",
		Tag:          dev.GitImageTag(),
	},
}

var _ = Describe("Deployment", func() {
	var c *dev.AryaCluster
	var d *dev.Deployment
	BeforeEach(func() {
		var cErr error
		c, cErr = provisionDummyAryaClusterIfNotExists()
		if cErr != nil {
			log.Fatalln(cErr)
		}
		dummyCfg.Cluster = c
		var dErr error
		d, dErr = dev.NewDeployment(dummyCfg)
		if dErr != nil {
			log.Fatalln(dErr)
		}

	})
	Describe("Creating a New Deployment", func() {
		It("Should create a new deployment without error", func() {
			_, err := dev.NewDeployment(dummyCfg)
			Expect(err).To(BeNil())
		})
	})
	Describe("Installing the Deployment", func() {
		It("Should install the deployment without error", func() {
			err := d.Install()
			Expect(err).To(BeNil())
		})
	})
	Describe("Re-deploying arya", func() {
		It("Should redeploy arya without error", func() {
			err := d.RedeployArya()
			Expect(err).To(BeNil())
		})
	})

})
