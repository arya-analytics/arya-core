package roachdriver_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/roach/roachdriver"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Clusterapi", func() {
	var clusterAPI *roachdriver.ClusterAPI
	BeforeEach(func() {
		clusterAPI = &roachdriver.ClusterAPI{
			Port:     mockDriver.HTTPPort,
			Host:     mockDriver.Host,
			Username: mockDriver.Username,
			Password: mockDriver.Password,
		}
	})
	It("Should authenticate correctly", func() {
		err := clusterAPI.Connect()
		Expect(err).To(BeNil())
	})
	Describe("Querying the API", func() {
		BeforeEach(func() {
			Expect(clusterAPI.Connect()).To(BeNil())
		})
		It("Should return the correct number of cluster API nodes", func() {
			nodes, err := clusterAPI.Nodes()
			Expect(err).To(BeNil())
			Expect(nodes[0].ID).To(Equal(1))
		})
	})
})
