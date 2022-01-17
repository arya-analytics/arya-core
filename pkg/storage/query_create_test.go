package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Create", func() {
	Describe("Create a  new item", func() {
		AfterEach(deleteDummyModel)
		It("Should create it without error", func() {
			err := dummyStorage.NewCreate().Model(dummyModel).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			err := dummyStorage.NewCreate().Model(dummyModel).Exec(dummyCtx)
			Expect(err).To(BeNil())
			m := &storage.ChannelConfig{}
			err = dummyStorage.NewRetrieve().Model(m).WherePK(dummyModel.ID).Exec(
				dummyCtx)
			Expect(err).To(BeNil())
			Expect(m.ID).To(Equal(dummyModel.ID))
		})
	})
	Describe("Bulk create items", func() {
		models := []*storage.ChannelConfig{
			&storage.ChannelConfig{
				ID:   9621,
				Name: "Cool Name 1",
			},
			&storage.ChannelConfig{
				ID:   9622,
				Name: "Cool Name 2",
			},
		}
		modelPks := []int{models[0].ID, models[1].ID}
		AfterEach(func() {
			if err := dummyStorage.NewDelete().Model(&models).WherePKs(modelPks).Exec(
				dummyCtx); err != nil {
				log.Fatalln(err)
			}
		})
		It("Should create without error", func() {
			err := dummyStorage.NewCreate().Model(&models).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			err := dummyStorage.NewCreate().Model(&models).Exec(dummyCtx)
			Expect(err).To(BeNil())
			var m []*storage.ChannelConfig
			err = dummyStorage.NewRetrieve().Model(&m).WherePKs(modelPks).Exec(
				dummyCtx)
			Expect(err).To(BeNil())
			Expect(m).To(HaveLen(2))
			Expect(m[0].ID).To(Equal(9621))
			Expect(m[1].ID).To(Equal(9622))

		})
	})
})
