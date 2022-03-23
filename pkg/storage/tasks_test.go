package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Tasks", func() {
	It("Should start and stop the task runner correctly", func() {
		var storeTwo = mock.NewStorage()
		mErr := storeTwo.NewMigrate().Exec(ctx)
		Expect(mErr).To(BeNil())
		Expect(storeTwo.Start(ctx, tasks.ScheduleWithAccel(100), tasks.ScheduleWithSilence())).To(BeNil())
		var err error
		go func() {
			err = <-storeTwo.Errors()
		}()
		time.Sleep(300 * time.Millisecond)
		storeTwo.Stop()
		Expect(err).To(BeNil())
	})
})
