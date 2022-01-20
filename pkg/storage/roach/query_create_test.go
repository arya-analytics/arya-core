package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	BeforeEach(migrate)
	BeforeEach(deleteDummyModel)
	AfterEach(deleteDummyModel)
	Describe("Create a new Channel Config", func() {
		var errNode, errChan error
		BeforeEach(func() {
			errNode = dummyEngine.NewCreate(dummyAdapter).Model(dummyNode).Exec(
				dummyCtx)
			errChan = dummyEngine.NewCreate(dummyAdapter).Model(dummyModel).Exec(dummyCtx)
		})
		AfterEach(func() {
			errNode = dummyEngine.NewDelete(dummyAdapter).Model(dummyNode).WherePK(
				dummyNode.ID).
				Exec(
					dummyCtx)
			errChan = dummyEngine.NewDelete(dummyAdapter).Model(dummyModel).WherePK(
				dummyModel.ID,
			).Exec(dummyCtx)
		})
		It("Should create it without error", func() {
			Expect(errNode).To(BeNil())
			Expect(errChan).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			m := &storage.ChannelConfig{
				Name: "Channel Config",
			}
			err := dummyEngine.NewRetrieve(dummyAdapter).Model(m).WherePK(dummyModel.
				ID).Exec(dummyCtx)
			Expect(err).To(BeNil())
			Expect(m.Name).To(Equal(dummyModel.Name))
			Expect(m.ID).To(Equal(dummyModel.ID))
		})
	})
	Describe("Edge cases + errors", func() {
		Context("No PK provided", func() {
			Context("Single Model", func() {
				It("Should return the correct error type", func() {
					m := &storage.Node{}
					err := dummyEngine.NewCreate(dummyAdapter).Model(m).Exec(dummyCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeNoPK))
				})
			})
			Context("Chain of models", func() {
				It("Should return the correct error type", func() {
					m := []*storage.Node{
						&storage.Node{
							ID: 125,
						},
						&storage.Node{},
					}
					err := dummyEngine.NewCreate(dummyAdapter).Model(&m).Exec(dummyCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeNoPK))
				})
			})
		})
		Context("Unique Violation", func() {
			It("Should return the correct error type", func() {
				commonPk := uuid.New()
				mOne := &storage.ChannelConfig{
					ID: commonPk,
				}
				err := dummyEngine.NewCreate(dummyAdapter).Model(mOne).Exec(dummyCtx)
				Expect(err).To(BeNil())
				mTwo := &storage.ChannelConfig{
					ID: commonPk,
				}
				err = dummyEngine.NewCreate(dummyAdapter).Model(mTwo).Exec(dummyCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeUniqueViolation))
			})
		})
	})
})
