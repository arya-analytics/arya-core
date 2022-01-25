package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Create", func() {
	Describe("Create a  new item", func() {
		AfterEach(deleteMockChannelCfg)
		It("Should create it without error", func() {
			err := mockStorage.NewCreate().Model(mockChannelCfg).Exec(mockCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			err := mockStorage.NewCreate().Model(mockChannelCfg).Exec(mockCtx)
			Expect(err).To(BeNil())
			m := &storage.ChannelConfig{}
			err = mockStorage.NewRetrieve().Model(m).WherePK(mockChannelCfg.ID).Exec(
				mockCtx)
			Expect(err).To(BeNil())
			Expect(m.ID).To(Equal(mockChannelCfg.ID))
		})
	})
	Describe("Object create items", func() {
		models := []*storage.ChannelConfig{
			&storage.ChannelConfig{
				ID:     uuid.New(),
				Name:   "Cool Name 1",
				NodeID: 1,
			},
			&storage.ChannelConfig{
				ID:     uuid.New(),
				Name:   "Cool Name 2",
				NodeID: 1,
			},
		}
		modelPks := []uuid.UUID{models[0].ID, models[1].ID}
		AfterEach(func() {
			if err := mockStorage.NewDelete().Model(&models).WherePKs(modelPks).Exec(
				mockCtx); err != nil {
				log.Fatalln(err)
			}
		})
		It("Should create without error", func() {
			err := mockStorage.NewCreate().Model(&models).Exec(mockCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			err := mockStorage.NewCreate().Model(&models).Exec(mockCtx)
			Expect(err).To(BeNil())
			var m []*storage.ChannelConfig
			err = mockStorage.NewRetrieve().Model(&m).WherePKs(modelPks).Exec(
				mockCtx)
			Expect(err).To(BeNil())
			Expect(m).To(HaveLen(2))
			Expect(m[1].ID == models[1].ID || m[1].ID == models[0].ID).To(BeTrue())
		})
	})
})
