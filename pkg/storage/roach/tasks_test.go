package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	. "github.com/onsi/ginkgo/v2"
	log "github.com/sirupsen/logrus"
	"time"

	. "github.com/onsi/gomega"
)

var _ = Describe("Tasks", func() {
	AfterEach(func() {
		err := engine.NewDelete(adapter).Model(&storage.Node{}).WherePK(1).Exec(ctx)
		Expect(err).To(BeNil())
	})
	It("Should create the missing nodes in the database", func() {
		tasks := engine.Tasks(
			adapter,
			tasks.ScheduleWithAccel(100),
			tasks.ScheduleWithSilence(),
		)
		go tasks.Start(ctx)
		go func() {
			log.Fatalln(<-tasks.Errors)
		}()
		time.Sleep(150 * time.Millisecond)
		count, err := engine.NewRetrieve(adapter).Model(&storage.Node{}).Count(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(Equal(1))
	})
})
