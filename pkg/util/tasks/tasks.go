package tasks

import (
	"context"
	"time"
)

// Action is the function executed on the specified Task.Interval. The error returned by a task will be sent to the
// Scheduler error handling mechanism. A context is provided to the action. This context is the same context as
// provided to Scheduler.Start. The SchedulerConfig is the config of the scheduler managing the Task. This config is
// particularly useful for handling logging conditions (such as ScheduleWithSilence).
// NOTE: SchedulerConfig is provided as a value, so editing the config won't make any changes to the Scheduler config.
type Action func(ctx context.Context, cfg SchedulerConfig) error

// Task is an action called on a specified interval. Task isn't useful on its own, and should be provided to a
// Scheduler so that it can be executed.
type Task struct {
	// Name the tasks for tracing purposes
	Name string
	// Interval defines how often Task.Action gets called.
	Interval time.Duration
	// See Action
	Action Action
}
