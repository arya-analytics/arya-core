package cluster_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StorageService", func() {
	It("Should create a new storage service", func() {
		Expect(func() { cluster.NewStorageService(store) }).ToNot(Panic())
	})
	It("Should be able to handle a request", func() {
		s := cluster.NewStorageService(store)
		Expect(s.CanHandle(query.NewRetrieve().Model(&models.ChannelConfig{}).Pack())).To(BeTrue())
	})
})
