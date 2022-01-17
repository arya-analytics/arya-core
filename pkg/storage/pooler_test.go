package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("pooler", func() {
	BeforeEach(func() {
		log.SetReportCaller(true)
	})
	Describe("Retrieving a new adapter", func() {
		It("Should retrieveQuery an adapter", func() {
			p := storage.UnsafeNewPooler()
			a := p.Retrieve(&mock.MDEngine{})
			Expect(len(a.ID().String())).To(Equal(len(uuid.New().String())))
		})
		It("Should retrieve the same adapter if queried twice", func() {
			p := storage.UnsafeNewPooler()
			aOne := p.Retrieve(&mock.MDEngine{})
			aTwo := p.Retrieve(&mock.MDEngine{})
			Expect(aOne.ID()).To(Equal(aTwo.ID()))
		})
	})
})
