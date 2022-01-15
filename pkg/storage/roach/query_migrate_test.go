package roach_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Migrator", func() {
	BeforeEach(migrate)
	Describe("Init Migrations", func() {
		log.SetReportCaller(true)
		It("Should create all of the tables correctly", func() {
			err := dummyEngine.NewMigrate(dummyAdapter).Verify(dummyCtx)
			Expect(err).To(BeNil())
		})
	})
})