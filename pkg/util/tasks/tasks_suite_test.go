package tasks_test

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ctx = context.Background()

func TestTasks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewTasks Suite")
}

// A SimpleScheduler that increments a counter.
func ExampleNewSimpleScheduler() {
	ctx := context.Background()
	count := 0
	s := tasks.NewSimpleScheduler(
		[]tasks.Task{
			{
				Name:     "increment counter",
				Interval: 250 * time.Millisecond,
				Action: func(ctx context.Context, cfg tasks.SchedulerConfig) error {
					count += 1
					return nil
				},
			},
		},
	)
	go s.Start(ctx)
	defer s.Stop()
	time.Sleep(550 * time.Millisecond)
	fmt.Println(count)
	// Output:
	// 2
}

// A BatchScheduler that compose two simple schedulers together.
func ExampleNewBatchScheduler() {
	ctx := context.Background()
	countOne := 0
	sOne := tasks.NewSimpleScheduler(
		[]tasks.Task{
			{
				Name:     "increment counter one",
				Interval: 250 * time.Millisecond,
				Action: func(ctx context.Context, cfg tasks.SchedulerConfig) error {
					countOne += 1
					return nil
				},
			},
		},
	)
	countTwo := 0
	sTwo := tasks.NewSimpleScheduler(
		[]tasks.Task{
			{
				Name:     "increment counter two",
				Interval: 500 * time.Millisecond,
				Action: func(ctx context.Context, cfg tasks.SchedulerConfig) error {
					countTwo += 1
					return nil
				},
			},
		},
	)
	s := tasks.NewBatchScheduler(sOne, sTwo)
	s.Start(ctx)
	defer s.Stop()
	time.Sleep(550 * time.Millisecond)
	fmt.Printf("counter one: %d, counter two: %d", countOne, countTwo)
	// Output:
	// counter one: 2, counter two: 1
}
