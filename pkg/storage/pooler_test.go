package storage_test

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var cfgChain = storage.ConfigChain{
	storage.MDStubConfig{},
	storage.CacheStubConfig{},
}

var _ = Describe("Pooler", func() {
	var p *storage.Pooler
	BeforeEach(func() {
		log.SetReportCaller(true)
		p = storage.NewPooler(cfgChain)
	})
	Describe("Retrieving a new adapter", func() {
		Context("The config was provided", func() {
			It("Should retrieve an adapter", func() {
				a, err := p.Retrieve(storage.EngineTypeMDStub)
				Expect(err).To(BeNil())
				Expect(a.Status()).To(Equal(storage.ConnStatusReady))
			})
			It("Should retrieve the same adapter if queried twice", func() {
				aOne, err := p.Retrieve(storage.EngineTypeMDStub)
				Expect(err).To(BeNil())
				aTwo, err := p.Retrieve(storage.EngineTypeMDStub)
				Expect(err).To(BeNil())
				Expect(aOne.ID()).To(Equal(aTwo.ID()))
			})
		})
		Context("The config was not provided", func() {
			It("Should return a config error", func() {
				_, err := p.Retrieve(storage.EngineTypeBulkStub)
				Expect(err).ToNot(BeNil())
				cfgErr, ok := err.(storage.ConfigError)
				Expect(ok).To(BeTrue())
				Expect(cfgErr.Error()).To(Equal(
					fmt.Sprintf("config not found in config chain %v", storage.
						EngineTypeBulkStub)))
			})
		})
		Context("The adapter does not exist", func() {
			It("Should return a pooler error", func() {
				_, err := p.Retrieve(storage.EngineTypeCacheStub)
				Expect(err).ToNot(BeNil())
				cfgErr, ok := err.(storage.PoolerError)
				Expect(ok).To(BeTrue())
				Expect(cfgErr.Error()).To(Equal(
					fmt.Sprintf("adapter type does not exist %v", storage.
						EngineTypeCacheStub)))
			})
		})
	})
})
