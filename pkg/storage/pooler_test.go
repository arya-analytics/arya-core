package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/stub"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var dummyEngines = []storage.Engine{
	&stub.MDEngine{},
}

var _ = Describe("Pooler", func() {
	var p *storage.Pooler
	BeforeEach(func() {
		log.SetReportCaller(true)
		p = storage.NewPooler()
	})
	Describe("Retrieving a new adapter", func() {
		Context("The config was provided", func() {
			It("Should retrieve an adapter", func() {
				a, err := p.Retrieve(&stub.MDEngine{})
				Expect(err).To(BeNil())
				Expect(len(a.ID().String())).To(Equal(len(uuid.New().String())))
			})
			It("Should retrieve the same adapter if queried twice", func() {
				aOne, err := p.Retrieve(&stub.MDEngine{})
				Expect(err).To(BeNil())
				aTwo, err := p.Retrieve(&stub.MDEngine{})
				Expect(err).To(BeNil())
				Expect(aOne.ID()).To(Equal(aTwo.ID()))
			})
		})
	})
})
