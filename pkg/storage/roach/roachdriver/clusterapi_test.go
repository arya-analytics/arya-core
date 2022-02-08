package roachdriver_test

import (
	"crypto/tls"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/storage/roach/roachdriver"
	. "github.com/onsi/ginkgo/v2"
	"log"
	"net/http"
)

var _ = Describe("Clusterapi", func() {
	It("Should authenticate correctly", func() {
		d := &mock.DriverRoach{}
		_, err := d.Connect()
		if err != nil {
			log.Fatalln(err)
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		c := &roachdriver.ClusterAPI{
			Port:     d.HTTPPort,
			Host:     d.Host,
			Username: d.Username,
			Password: d.Password,
		}
		c.Connect()
		c.Nodes()
	})
})
