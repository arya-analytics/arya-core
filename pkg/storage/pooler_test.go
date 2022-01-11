package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var cfgChain = storage.ConfigChain{
	storage.MDStubConfig{},
}


var _ = Describe("Pooler", func() {
	var p *storage.Pooler
	BeforeEach(func() {
		log.SetReportCaller(true)
		p = storage.NewPooler(cfgChain)
	})
	Describe("Retrieving a new adapter", func() {
		It("Should retrieve an adapter", func() {
			a := p.Retrieve(storage.EngineRoleMetaData)
			Expect(a.Status()).To(Equal(storage.ConnStatusReady))
		})
		It("Should retrieve the same adapter if queried twice", func() {
			aOne := p.Retrieve(storage.EngineRoleMetaData)
			aTwo := p.Retrieve(storage.EngineRoleMetaData)
			Expect(aOne.ID()).To(Equal(aTwo.ID()))
		})
	})
})
