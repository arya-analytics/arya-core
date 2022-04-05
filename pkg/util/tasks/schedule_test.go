package tasks_test

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ScheduleSimple", func() {
	Describe("Standard usage", func() {
		Context("base Schedule", func() {
			It("Should execute tasks at the correct interval", func() {
				count := 0
				s := tasks.NewScheduleSimple([]tasks.Task{
					{
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
							count += 1
							return nil
						},
					},
				},
					tasks.ScheduleWithName("test tasks"),
				)
				go s.Start(ctx)
				time.Sleep(625 * time.Millisecond)
				s.Stop()
				Expect(count).To(Equal(2))
			})
			It("Should pipe to the errors channel when a task fails", func() {
				s := tasks.NewScheduleSimple([]tasks.Task{
					{
						Name:     "bad task",
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
							return errors.New("a terrible error")
						},
					},
				}, tasks.ScheduleWithSilence())
				go s.Start(ctx)
				time.Sleep(300 * time.Millisecond)
				Expect(<-s.Errors()).ToNot(BeNil())
				s.Stop()
			})
			It("Should break out of the scheduler when the context is cancelled", func() {
				ctxWithCancel, cancel := context.WithCancel(ctx)
				count := 0
				s := tasks.NewScheduleSimple([]tasks.Task{
					{
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
							count += 1
							return nil
						},
					},
				}, tasks.ScheduleWithSilence())
				go s.Start(ctxWithCancel)
				time.Sleep(375 * time.Millisecond)
				cancel()
				Expect(count).To(Equal(1))
			})
			Context("Acceleration", func() {
				It("Should accelerate the scheduler correctly", func() {
					count := 0
					s := tasks.NewScheduleSimple([]tasks.Task{
						{
							Interval: 250 * time.Millisecond,
							Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
								count += 1
								return nil
							},
						},
					},
						tasks.ScheduleWithSilence(),
						tasks.ScheduleWithAccel(5),
					)
					go s.Start(ctx)
					time.Sleep(625 * time.Millisecond)
					s.Stop()
					Expect(count).To(Equal(12))
				})
			})
			Context("Multiple Tasks", func() {
				It("Should execute the tasks at the correct interval", func() {
					countOne, countTwo := 0, 0
					s := tasks.NewScheduleSimple([]tasks.Task{
						{
							Interval: 250 * time.Millisecond,
							Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
								countOne += 1
								return nil
							},
						},
						{
							Interval: 500 * time.Millisecond,
							Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
								countTwo += 1
								return nil
							},
						},
					},
						tasks.ScheduleWithSilence(),
						tasks.ScheduleWithAccel(5),
					)
					go s.Start(ctx)
					time.Sleep(625 * time.Millisecond)
					s.Stop()
					Expect(countOne).To(Equal(12))
					Expect(countTwo).To(Equal(6))
				})
			})
		})
		Context("Batch Schedule", func() {
			It("Should execute tasks at the correct interval", func() {
				count := 0
				s := tasks.NewScheduleSimple([]tasks.Task{
					{
						Interval: 250 * time.Millisecond,
						Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
							count += 1
							return nil
						},
					},
					{
						Interval: 300 * time.Millisecond,
						Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
							return nil
						},
					},
					{
						Interval: 350 * time.Millisecond,
						Action: func(ctx context.Context, cfg tasks.ScheduleConfig) error {
							return nil
						},
					},
				},
					tasks.ScheduleWithName("test tasks"),
					tasks.ScheduleWithSilence(),
				)
				batchScheduler := tasks.NewScheduleBatch(s)
				go batchScheduler.Start(ctx)
				time.Sleep(625 * time.Millisecond)
				var err error
				go func() {
					err = <-batchScheduler.Errors()
				}()
				batchScheduler.Stop()
				Expect(count).To(Equal(2))
				Expect(err).To(BeNil())
			})
		})
	})
})
