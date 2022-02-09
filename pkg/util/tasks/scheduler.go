// Package tasks holds utilities for executing actions at specified intervals.
// It includes a Task struct that defines a name, interval, and action. This Task can be provided
// to a Scheduler that manages its execution. SchedulerBatch can be used to compose multiple Scheduler together and run
// them as a set of independent goroutines.
//
// Error Handling
//
// Task.Action returns an error. This error can be used to pass issues encountered while running your task to the
// Scheduler error handling mechanisms. When a Scheduler encounters an error, it will log the error to standard output
// and pipe the error to the channel returned by Scheduler.Errors(). It's probably a good idea to listen to this channel
// and handle the error accordingly.
//
// When the Scheduler encounters an error while executing a task, it will keep executing the task on the specified
// interval, even if the task repeatedly fails.
package tasks

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

// |||| SCHEDULER ||||

// Scheduler aggregates and executes tasks on their specified interval.
// The tasks package provides two schedulers.
// SchedulerSimple executes a base set of provided tasks.
// SchedulerBatch composes multiple schedulers and runs them concurrently.
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

const defaultAccel = 1

// SchedulerSimple executes a set of tasks at their specified interval.
// To create a new scheduler, use NewSimpleScheduler.
// SchedulerSimple doesn't start any goroutines.
type SchedulerSimple struct {
	Tasks    []Task
	cfg      *SchedulerConfig
	chanErr  chan error
	chanStop chan bool
}

// NewSimpleScheduler creates a new Scheduler with the specified tasks.
// Optional SchedulerOpts can be provided to modify the scheduler's behavior.
func NewSimpleScheduler(tasks []Task, opts ...SchedulerOpt) Scheduler {
	s := &SchedulerSimple{
		Tasks:    tasks,
		chanErr:  make(chan error),
		chanStop: make(chan bool),
		cfg: &SchedulerConfig{
			Accel: defaultAccel,
		},
	}
	s.bindOpts(opts...)
	return s
}

// Start implements Scheduler.
// NOTE: This is a blocking operation. It's common to run Start as a new goroutine.
func (s *SchedulerSimple) Start(ctx context.Context) {
	s.logStart()
	t, t0 := time.NewTicker(s.tickInterval()), time.Now()
	defer t.Stop()
	for {
		select {
		case ct := <-t.C:
			s.exec(ctx, ct.Sub(t0))
		case <-s.chanStop:
			s.logStop()
			return
		case <-ctx.Done():
			s.logStop()
			return
		}
	}
}

func (s *SchedulerSimple) Stop() {
	s.chanStop <- true
}

// Errors implements Scheduler.
func (s *SchedulerSimple) Errors() chan error {
	return s.chanErr
}

func (s *SchedulerSimple) tickInterval() time.Duration {
	var intervals []time.Duration
	for _, t := range s.Tasks {
		intervals = append(intervals, t.Interval)
	}
	i := s.accelerate(durationGCD(intervals...))
	return i
}

const taskExecThreshold = 50 * time.Millisecond

func (s *SchedulerSimple) exec(ctx context.Context, t time.Duration) {
	for _, task := range s.Tasks {
		if t%s.accelerate(task.Interval) < taskExecThreshold {
			if err := task.Action(ctx, *s.cfg); err != nil {
				s.logTaskFailure(task, err)
				s.chanErr <- err
			}
		}
	}
}

const minToNanoSec = 1000000000 * 60

func (s *SchedulerSimple) accelerate(t time.Duration) time.Duration {
	return time.Duration((t.Minutes() / s.cfg.Accel) * minToNanoSec)
}

func (s *SchedulerSimple) bindOpts(opts ...SchedulerOpt) {
	for _, opt := range opts {
		opt(s.cfg)
	}
}

func (s *SchedulerSimple) logStart() {
	if s.cfg.Name != "" && !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name":          s.cfg.Name,
			"task_count":    len(s.Tasks),
			"tick_interval": s.tickInterval(),
			"Accel":         s.cfg.Accel,
		}).Infof("Starting %s", s.cfg.Name)
	}
}

func (s *SchedulerSimple) logStop() {
	if s.cfg.Name != "" && !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name": s.cfg.Name,
		}).Infof("Stopping %s", s.cfg.Name)

	}
}

func (s *SchedulerSimple) logTaskFailure(task Task, err error) {
	if !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name":      s.cfg.Name,
			"task_name": task.Name,
		}).Errorf("Task failed! %s", err)
	}
}

// |||| BATCH ||||

// SchedulerBatch composes a set of Scheduler and runs each of them in a separate
// go -routine. You can construct hierarchies of tasks schedulers by composing
// multiple SchedulerBatch and SchedulerSimple.
type SchedulerBatch struct {
	schedulers []Scheduler
	chanErr    chan error
}

// NewBatchScheduler creates a new SchedulerBatch from a set of Scheduler.
func NewBatchScheduler(schedulers ...Scheduler) Scheduler {
	return &SchedulerBatch{schedulers: schedulers}
}

// Start implements Scheduler.
// NOTE: this method is non-blocking, and starts a new goroutine for each child Scheduler.
func (sb *SchedulerBatch) Start(ctx context.Context) {
	for _, s := range sb.schedulers {
		go s.Start(ctx)
		go sb.pipeErrors(s.Errors())
	}
}

// Stop implements Scheduler.
func (sb *SchedulerBatch) Stop() {
	for _, s := range sb.schedulers {
		s.Stop()
	}
}

// Errors implements Scheduler.
func (sb *SchedulerBatch) Errors() chan error {
	return sb.chanErr
}

func (sb *SchedulerBatch) pipeErrors(chanErr chan error) {
	err := <-chanErr
	sb.chanErr <- err
}

// |||| SCHEDULER OPTIONS ||||

// SchedulerConfig holds configuration information for a Scheduler. It shouldn't be instantiated directly, and should
// instead be configured by passing SchedulerOpt to one of the Scheduler constructors.
type SchedulerConfig struct {
	Name   string
	Accel  float64
	Silent bool
}

// SchedulerOpt is a modifier function that allows arbitrary configuration of a scheduler.
type SchedulerOpt func(opts *SchedulerConfig)

// ScheduleWithName assigns a name to the scheduler. This is primarily for logging purposes.
func ScheduleWithName(name string) SchedulerOpt {
	return func(opts *SchedulerConfig) { opts.Name = name }
}

// ScheduleWithAccel accelerates the scheduler by the specified multiple. A utility mostly for testing tasks with long
// Task.Interval values. It's probably not a good idea to use this in production, as it modifies the sense of time, which
// can be confusing.
func ScheduleWithAccel(accel float64) SchedulerOpt {
	return func(opts *SchedulerConfig) { opts.Accel = accel }
}

// ScheduleWithSilence disables all logging internal to the scheduler. NOTE: This won't disable logging from Tasks,
// unless the task accesses SchedulerConfig and implements the appropriate logic.
func ScheduleWithSilence() SchedulerOpt {
	return func(opts *SchedulerConfig) { opts.Silent = true }
}

// |||| MATH UTILITIES ||||

func durationGCD(durs ...time.Duration) time.Duration {
	if len(durs) == 0 {
		panic("cannot get the duration g with no arguments")
	}
	if len(durs) < 2 {
		return durs[0]
	}
	g := gcd(durs[0].Nanoseconds(), durs[1].Nanoseconds())
	if len(durs) == 2 {
		return time.Duration(g)
	}
	for _, dur := range durs[2:] {
		g = gcd(dur.Nanoseconds(), g)
	}
	return time.Duration(g)
}

func gcd(a, b int64) int64 {
	bigA, bigB := big.NewInt(a), big.NewInt(b)
	gcd := big.NewInt(0)
	return gcd.GCD(nil, nil, bigA, bigB).Int64()
}
