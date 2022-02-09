# tasks
--
    import "."

Package tasks holds utilities for executing actions at specified intervals. It
includes a Task struct that defines a name, interval, and action. This Task can
be provided to a Scheduler that manages its execution. SchedulerBatch can be
used to compose multiple Scheduler together and run them as a set of independent
goroutines.


### Error Handling

Task.Action returns an error. This error can be used to pass issues encountered
while running your task to the Scheduler error handling mechanisms. When a
Scheduler encounters an error, it will log the error to standard output and pipe
the error to the channel returned by Scheduler.Errors(). It's probably a good
idea to listen to this channel and handle the error accordingly.

When the Scheduler encounters an error while executing a task, it will keep
executing the task on the specified interval, even if the task repeatedly fails.

## Usage

#### type Scheduler

```go
type Scheduler interface {
	// Start starts the scheduler. See implementation specific details for more information on the behavior of Start.
	// The provided context can be used to stop the scheduler (by calling cancel).
	Start(ctx context.Context)
	// Stop stops the scheduler. It's common to call Stop within a defer statement.
	Stop()
	// Errors returns a channel that pipes errors encountered during task scheduling and
	// execution.
	Errors() chan error
}
```

Scheduler aggregates and executes tasks on their specified interval. The tasks
package provides two schedulers. SchedulerSimple executes a base set of provided
tasks. SchedulerBatch composes multiple schedulers and runs them concurrently.

#### func  NewBatchScheduler

```go
func NewBatchScheduler(schedulers ...Scheduler) Scheduler
```
NewBatchScheduler creates a new SchedulerBatch from a set of Scheduler.

#### func  NewSimpleScheduler

```go
func NewSimpleScheduler(tasks []Task, opts ...SchedulerOpt) Scheduler
```
NewSimpleScheduler creates a new Scheduler with the specified tasks. Optional
SchedulerOpts can be provided to modify the scheduler's behavior.

#### type SchedulerBatch

```go
type SchedulerBatch struct {
}
```

SchedulerBatch composes a set of Scheduler and runs each of them in a separate
go -routine. You can construct hierarchies of tasks schedulers by composing
multiple SchedulerBatch and SchedulerSimple.

#### func (*SchedulerBatch) Errors

```go
func (sb *SchedulerBatch) Errors() chan error
```
Errors implements Scheduler.

#### func (*SchedulerBatch) Start

```go
func (sb *SchedulerBatch) Start(ctx context.Context)
```
Start implements Scheduler. NOTE: this method is non-blocking, and starts a new
goroutine for each child Scheduler.

#### func (*SchedulerBatch) Stop

```go
func (sb *SchedulerBatch) Stop()
```
Stop implements Scheduler.

#### type SchedulerConfig

```go
type SchedulerConfig struct {
	Name   string
	Accel  float64
	Silent bool
}
```

SchedulerConfig holds configuration information for a Scheduler. It shouldn't be
instantiated directly, and should instead configured by passing SchedulerOpt to
one of the Scheduler constructors.

#### type SchedulerOpt

```go
type SchedulerOpt func(opts *SchedulerConfig)
```

SchedulerOpt is a modifier function that allows arbitrary configuration of a
scheduler.

#### func  ScheduleWithAccel

```go
func ScheduleWithAccel(accel float64) SchedulerOpt
```
ScheduleWithAccel accelerates the scheduler by the specified multiple. A utility
mostly for testing tasks with long Task.Interval values. It's probably not a
good idea to use this in production, as it modifies the sense of time, which can
be confusing.

#### func  ScheduleWithName

```go
func ScheduleWithName(name string) SchedulerOpt
```
ScheduleWithName assigns a name to the scheduler. This is primarily for logging
purposes.

#### func  ScheduleWithSilence

```go
func ScheduleWithSilence() SchedulerOpt
```
ScheduleWithSilence disables all logging internal to the scheduler. NOTE: This
won't disable logging from Tasks, unless the task accesses SchedulerConfig and
implements the appropriate logic.

#### type SchedulerSimple

```go
type SchedulerSimple struct {
	Tasks []Task
}
```

SchedulerSimple executes a set of tasks at their specified interval. To create a
new scheduler, use NewSimpleScheduler. SchedulerSimple doesn't start any
goroutines.

#### func (*SchedulerSimple) Errors

```go
func (s *SchedulerSimple) Errors() chan error
```
Errors implements Scheduler.

#### func (*SchedulerSimple) Start

```go
func (s *SchedulerSimple) Start(ctx context.Context)
```
Start implements Scheduler. NOTE: This is a blocking operation. It's common to
run Start as a new goroutine.

#### func (*SchedulerSimple) Stop

```go
func (s *SchedulerSimple) Stop()
```

#### type Task

```go
type Task struct {
	// Name the tasks for tracing purposes
	Name string
	// Interval defines how often Task.Action gets called.
	Interval time.Duration
	// Action is the function executed on the specified Task.Interval. The error returned by a task will be sent to the
	// Scheduler error handling mechanism. A context is provided to the action. This context is the same context as
	// provided to Scheduler.Start. The SchedulerConfig is the config of the scheduler managing the Task. This config is
	// particularly useful for handling logging conditions (such as ScheduleWithSilence).
	// NOTE: SchedulerConfig is provided as a value, so editing the config won't make any changes to the Scheduler config.
	Action func(ctx context.Context, cfg SchedulerConfig) error
}
```

Task is an action called on a specified interval. Task isn't useful on its own,
and should be provided to a Scheduler so that it can be executed.
