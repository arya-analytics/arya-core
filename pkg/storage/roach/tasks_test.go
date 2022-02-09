package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	. "github.com/onsi/ginkgo/v2"
	log "github.com/sirupsen/logrus"
	"time"

	. "github.com/onsi/gomega"
)

var _ = Describe("Tasks", func() {
	Describe("Node Synchronization", func() {
		AfterEach(func() {
			err := engine.NewDelete(adapter).Model(&storage.Node{}).WherePK(1).Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should create the missing nodes", func() {
			tasks := engine.Tasks(
				adapter,
				tasks.ScheduleWithAccel(100),
				tasks.ScheduleWithName("roach tasks"),
			)
			go tasks.Start(ctx)
			go func() {
				log.Fatalln(<-tasks.Errors)
			}()
			tasks.Stop()
			time.Sleep(150 * time.Millisecond)
			count, err := engine.NewRetrieve(adapter).Model(&storage.Node{}).Count(ctx)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})
		Context("Extra nodes", func() {
			bunDB := roach.UnsafeConn(adapter)
			var extraNode *storage.Node
			BeforeEach(func() {
				extraNode = &storage.Node{ID: 2}
			})
			JustBeforeEach(func() {
				err := engine.NewCreate(adapter).Model(extraNode).Exec(ctx)
				Expect(err).To(BeNil())
				count, err := bunDB.NewSelect().Table("nodes").
					Count(ctx)
				Expect(count).To(Equal(1))
			})
			It("Should remove the extra nodes", func() {
				tasks := engine.Tasks(
					adapter,
					tasks.ScheduleWithAccel(100),
					tasks.ScheduleWithName("roach tasks"),
				)
				go tasks.Start(ctx)
				go func() {
					log.Fatalln(<-tasks.Errors)
				}()
				time.Sleep(150 * time.Millisecond)
				tasks.Stop()
				count, err := bunDB.NewSelect().Table("nodes").
					Count(ctx)
				Expect(count).To(Equal(1))
				Expect(err).To(BeNil())
				rErr := engine.NewRetrieve(adapter).Model(extraNode).WherePK(
					extraNode.ID).Exec(ctx)
				Expect(rErr).ToNot(BeNil())
				Expect(rErr.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
			})
			// Just in case we don't delete the extra node
			JustAfterEach(func() {
				err := engine.NewDelete(adapter).Model(extraNode).WherePK(extraNode.ID).Exec(ctx)
				if err != nil {
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
				}
				Expect(err).To(BeNil())
			})
		})
	})
})
