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
			svc = chanstream.NewService(store, struct{}{})
		})
		It("Should return false for a query it can't handle", func() {
			p := query.NewRetrieve().Model(&modelMock.ModelA{}).Pack()
			Expect(svc.CanHandle(p)).To(BeFalse())
		})
		It("Should return true for a query it can handle", func() {
			p := query.NewRetrieve().Model(&models.ChannelSample{}).Pack()
			Expect(svc.CanHandle(p)).To(BeTrue())
		})
	})
	Describe("Node Is Local", func() {
		BeforeEach(func() {})
	})
})
