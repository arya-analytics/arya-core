package chanstream_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanstream"
	"github.com/arya-analytics/aryacore/pkg/models"
	modelMock "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Describe("CanHandle", func() {
		var (
			svc *chanstream.Service
		)
		BeforeEach(func() {
			svc = chanstream.NewService(store.Exec, struct{}{})
		})
		It("Should return false for a query it can't handle", func() {
			c := make(chan *modelMock.ModelA)
			p := query.NewRetrieve().Model(&c).Pack()
			Expect(svc.CanHandle(p)).To(BeFalse())
		})
		It("Should return true for a query it can handle", func() {
			c := make(chan *models.ChannelSample)
			p := query.NewRetrieve().Model(&c).Pack()
			Expect(svc.CanHandle(p)).To(BeTrue())
		})
	})
	Describe("Node Is Local", func() {
		BeforeEach(func() {})
	})
})
