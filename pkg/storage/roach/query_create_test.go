package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	BeforeEach(migrate)
	BeforeEach(deleteMockModel)
	AfterEach(deleteMockModel)
	Describe("Create a new Channel Config", func() {
		var errChan error
		BeforeEach(func() {
			errChan = mockEngine.NewCreate(mockAdapter).Model(mockModel).Exec(mockCtx)
		})
		AfterEach(func() {
			errChan = mockEngine.NewDelete(mockAdapter).Model(mockModel).WherePK(
				mockModel.ID,
			).Exec(mockCtx)
		})
		It("Should create it without error", func() {
			Expect(errChan).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			m := &storage.ChannelConfig{
				Name:   "Channel Config",
				NodeID: 1,
			}
			err := mockEngine.NewRetrieve(mockAdapter).Model(m).WherePK(mockModel.
				ID).Exec(mockCtx)
			Expect(err).To(BeNil())
			Expect(m.Name).To(Equal(mockModel.Name))
			Expect(m.ID).To(Equal(mockModel.ID))
		})
	})
	Describe("Edge cases + errors", func() {
		Context("No PK provided", func() {
			Context("Single Model", func() {
				It("Should return the correct error type", func() {
					m := &storage.Node{}
					err := mockEngine.NewCreate(mockAdapter).Model(m).Exec(mockCtx)
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
					err := mockEngine.NewCreate(mockAdapter).Model(&m).Exec(mockCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeNoPK))
				})
			})
		})
		Context("Unique Violation", func() {
			It("Should return the correct error type", func() {
				commonPk := uuid.New()
				mOne := &storage.ChannelConfig{
					ID:     commonPk,
					NodeID: mockNode.ID,
				}
				err := mockEngine.NewCreate(mockAdapter).Model(mOne).Exec(mockCtx)
				Expect(err).To(BeNil())
				mTwo := &storage.ChannelConfig{
					ID:     commonPk,
					NodeID: mockNode.ID,
				}
				err = mockEngine.NewCreate(mockAdapter).Model(mTwo).Exec(mockCtx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeUniqueViolation))
			})
		})
	})
})
