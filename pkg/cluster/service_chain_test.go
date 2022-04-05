package cluster_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceChain", func() {
	It("Should execute the query on the correct service", func() {
		s := cluster.NewStorageService(store)
		svc := cluster.ServiceChain{s}
		Expect(svc.Exec(context.Background(), query.NewCreate().Model(&models.Node{ID: 1}).Pack())).To(BeNil())
	})

})
