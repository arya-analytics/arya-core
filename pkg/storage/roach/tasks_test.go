package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	. "github.com/onsi/ginkgo/v2"
	log "github.com/sirupsen/logrus"
	"time"

	. "github.com/onsi/gomega"
)

const (
	taskAccel     = 50
	sleepDuration = 300 * time.Millisecond
)

var _ = Describe("NewTasks", func() {
	Describe("Node Synchronization", func() {
		AfterEach(func() {
			err := engine.NewDelete().Model(&models.Node{}).WherePK(1).Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should create the missing nodes", func() {
			tasks := engine.NewTasks(
				tasks.ScheduleWithSilence(),
				tasks.ScheduleWithAccel(taskAccel),
				tasks.ScheduleWithName("roach tasks"),
			)
			go tasks.Start(ctx)
			go func() {
				err := <-tasks.Errors()
				if err != nil {
					log.Fatalln(err)
				}
			}()
			time.Sleep(sleepDuration)
			tasks.Stop()
			var count int
			err := engine.NewRetrieve().Model(&models.Node{}).Calc(query.CalcCount, "id", &count).Exec(ctx)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(1))
		})
		Context("Extra nodes", func() {
			bunDB := roach.UnsafeConn(pool.Retrieve(engine))
			var extraNode *models.Node
			BeforeEach(func() {
				extraNode = &models.Node{ID: 2}
			})
			JustBeforeEach(func() {
				cErr := engine.NewCreate().Model(extraNode).Exec(ctx)
				Expect(cErr).To(BeNil())
				count, sErr := bunDB.NewSelect().Table("nodes").
					Count(ctx)
				Expect(sErr).To(BeNil())
				Expect(count).To(Equal(1))
			})
			It("Should remove the extra nodes", func() {
				tasks := engine.NewTasks(
					tasks.ScheduleWithSilence(),
					tasks.ScheduleWithAccel(taskAccel),
					tasks.ScheduleWithName("roach tasks"),
				)
				go tasks.Start(ctx)
				go func() {
					err := <-tasks.Errors()
					if err != nil {
						log.Fatalln(err)
					}
				}()
				time.Sleep(sleepDuration)
				tasks.Stop()
				count, err := bunDB.NewSelect().Table("nodes").
					Count(ctx)
				Expect(count).To(Equal(1))
				Expect(err).To(BeNil())
				rErr := engine.NewRetrieve().Model(extraNode).WherePK(
					extraNode.ID).Exec(ctx)
				Expect(rErr).ToNot(BeNil())
				Expect(rErr.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
			})
			// Just in case we don't del the extra node
			JustAfterEach(func() {
				err := engine.NewDelete().Model(extraNode).WherePK(extraNode.ID).Exec(ctx)
				if err != nil {
					Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
				}
				Expect(err).To(BeNil())
			})
		})
	})
})
