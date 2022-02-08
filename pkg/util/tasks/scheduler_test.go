package tasks_test

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Scheduler", func() {
	Describe("Standard usage", func() {
		Context("Basics", func() {
			It("Should execute tasks at the correct interval", func() {
				count := 0
				s := tasks.NewScheduler([]tasks.Task{
					{
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context) error {
							count += 1
							return nil
						},
					},
				}, 125*time.Millisecond,
					tasks.ScheduleWithName("test tasks"),
				)
				go s.Start(ctx)
				time.Sleep(625 * time.Millisecond)
				Expect(count).To(Equal(2))
			})
			It("Should pipe to the errors channel when a task fails", func() {
				s := tasks.NewScheduler([]tasks.Task{
					{
						Name:     "bad task",
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context) error {
							return errors.New("a terrible error")
						},
					},
				}, 125*time.Millisecond, tasks.ScheduleWithSilence())
				go s.Start(ctx)
				Expect(<-s.Errors).ToNot(BeNil())
			})
			It("Should break out of the scheduler when the context is cancelled", func() {
				ctxWithCancel, cancel := context.WithCancel(ctx)
				count := 0
				s := tasks.NewScheduler([]tasks.Task{
					{
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context) error {
							count += 1
							return nil
						},
					},
				}, 125*time.Millisecond, tasks.ScheduleWithSilence())
				go s.Start(ctxWithCancel)
				time.Sleep(375 * time.Millisecond)
				cancel()
				Expect(count).To(Equal(1))
			})
		})
		Context("Acceleration", func() {
			It("Should accelerate the scheduler correctly", func() {
				count := 0
				s := tasks.NewScheduler([]tasks.Task{
					{
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context) error {
							count += 1
							return nil
						},
					},
				}, 125*time.Millisecond,
					tasks.ScheduleWithSilence(),
					tasks.ScheduleWithAccel(5),
				)
				go s.Start(ctx)
				time.Sleep(625 * time.Millisecond)
				Expect(count).To(Equal(25))
			})
		})
	})
})
