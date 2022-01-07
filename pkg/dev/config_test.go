package dev_test

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/dev"
	"github.com/arya-analytics/aryacore/pkg/util/kubectl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"os/exec"
)

var _ = Describe("Config", func() {
	log.SetReportCaller(true)
	var c *dev.AryaCluster
	var cErr error
	BeforeEach(func() {
		c, cErr = provisionDummyAryaClusterIfNotExists()
		if cErr != nil {
			log.Fatalln(cErr)
		}
	})
	Describe("Authenticating a cluster", func() {
		It("Should generate the correct authentication secret", func() {
			nodes := c.Nodes()
			fmt.Println(nodes[0].VM.Name())
			if err := kubectl.SwitchContext(nodes[0].VM.Name()); err != nil {
				log.Fatalln(err)
			}
			o, err := exec.Command("bash", "-c",
				"kubectl get secret | grep regcred | awk '{print $1}'").Output()
			if err != nil {
				log.Fatalln(err)
			}
			Expect(string(o[:])).To(Equal("regcred\n"))
		})
	})
})
