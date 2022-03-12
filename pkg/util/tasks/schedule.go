// Package tasks holds utilities for executing actions at specified intervals.
// It includes a Task struct that defines a name, interval, and action. This Task can be provided
// to a Schedule that manages its execution. ScheduleBatch can be used to compose multiple Schedule together and run
// them as a set of independent goroutines.
//
// Error Handling
//
// Task.Action returns an error. This error can be used to pass issues encountered while running your task to the
// Schedule error handling mechanisms. When a Schedule encounters an error, it will log the error to standard output
// and pipe the error to the channel returned by Schedule.Errors(). It's probably a good idea to listen to this channel
// and handle the error accordingly.
//
// When the Schedule encounters an error while executing a task, it will keep executing the task on the specified
// interval, even if the task repeatedly fails.
package tasks

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

// |||| SCHEDULER ||||

// Schedule aggregates and executes tasks on their specified interval.
// The tasks package provides two schedulers.
// ScheduleSimple executes a base set of provided tasks.
// ScheduleBatch composes multiple schedulers and runs them concurrently.
type Schedule interface {
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

// ScheduleSimple executes a set of tasks at their specified interval.
// To create a new scheduler, use NewScheduleSimple.
// ScheduleSimple doesn't start any goroutines.
type ScheduleSimple struct {
	Tasks    []Task
	cfg      *ScheduleConfig
	chanErr  chan error
	chanStop chan bool
}

// NewScheduleSimple creates a new Schedule with the specified tasks.
// Optional SchedulerOpts can be provided to modify the scheduler's behavior.
func NewScheduleSimple(tasks []Task, opts ...ScheduleOpt) Schedule {
	s := &ScheduleSimple{
		Tasks:    tasks,
		chanErr:  make(chan error),
		chanStop: make(chan bool),
		cfg: &ScheduleConfig{
			Accel: defaultAccel,
		},
	}
	s.bindOpts(opts...)
	return s
}

// Start implements Schedule.
// NOTE: This is a blocking operation. It's common to run Start as a new goroutine.
func (s *ScheduleSimple) Start(ctx context.Context) {
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

// Stop implements Schedule.
func (s *ScheduleSimple) Stop() {
	s.chanStop <- true
}

// Errors implements Schedule.
func (s *ScheduleSimple) Errors() chan error {
	return s.chanErr
}

func (s *ScheduleSimple) tickInterval() time.Duration {
	var intervals []time.Duration
	for _, t := range s.Tasks {
		intervals = append(intervals, t.Interval)
	}
	i := s.accelerate(durationGCD(intervals...))
	return i
}

const taskExecThreshold = 50 * time.Millisecond

func (s *ScheduleSimple) exec(ctx context.Context, t time.Duration) {
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

func (s *ScheduleSimple) accelerate(t time.Duration) time.Duration {
	return time.Duration((t.Minutes() / s.cfg.Accel) * minToNanoSec)
}

func (s *ScheduleSimple) bindOpts(opts ...ScheduleOpt) {
	for _, opt := range opts {
		opt(s.cfg)
	}
}

func (s *ScheduleSimple) logStart() {
	if s.cfg.Name != "" && !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name":          s.cfg.Name,
			"task_count":    len(s.Tasks),
			"tick_interval": s.tickInterval(),
			"Accel":         s.cfg.Accel,
		}).Infof("Starting %s", s.cfg.Name)
	}
}

func (s *ScheduleSimple) logStop() {
	if s.cfg.Name != "" && !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name": s.cfg.Name,
		}).Infof("Stopping %s", s.cfg.Name)

	}
}

func (s *ScheduleSimple) logTaskFailure(task Task, err error) {
	if !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name":      s.cfg.Name,
			"task_name": task.Name,
		}).Errorf("Task failed! %s", err)
	}
}

// |||| BATCH ||||

// ScheduleBatch composes a set of Schedule and runs each of them in a separate
// go -routine. You can construct hierarchies of tasks schedulers by composing
// multiple ScheduleBatch and ScheduleSimple.
type ScheduleBatch struct {
	schedulers []Schedule
	chanErr    chan error
}

// NewScheduleBatch creates a new ScheduleBatch from a set of Schedule.
func NewScheduleBatch(schedulers ...Schedule) Schedule {
	return &ScheduleBatch{schedulers: schedulers}
}

// Start implements Schedule.
// NOTE: this method is non-blocking, and starts a new goroutine for each child Schedule.
func (sb *ScheduleBatch) Start(ctx context.Context) {
	for _, s := range sb.schedulers {
		go s.Start(ctx)
		go sb.pipeErrors(s.Errors())
	}
}

// Stop implements Schedule.
func (sb *ScheduleBatch) Stop() {
	for _, s := range sb.schedulers {
		s.Stop()
	}
}

// Errors implements Schedule.
func (sb *ScheduleBatch) Errors() chan error {
	return sb.chanErr
}

func (sb *ScheduleBatch) pipeErrors(chanErr chan error) {
	err := <-chanErr
	sb.chanErr <- err
}

// |||| SCHEDULER OPTIONS ||||

// ScheduleConfig holds configuration information for a Schedule. It shouldn't be instantiated directly, and should
// instead be configured by passing ScheduleOpt to one of the Schedule constructors.
type ScheduleConfig struct {
	Name   string
	Accel  float64
	Silent bool
}

// ScheduleOpt is a modifier function that allows arbitrary configuration of a scheduler.
type ScheduleOpt func(opts *ScheduleConfig)

// ScheduleWithName assigns a name to the scheduler. This is primarily for logging purposes.
func ScheduleWithName(name string) ScheduleOpt {
	return func(opts *ScheduleConfig) { opts.Name = name }
}

// ScheduleWithAccel accelerates the scheduler by the specified multiple. A utility mostly for testing tasks with long
// Task.Interval values. It's probably not a good idea to use this in production, as it modifies the sense of time, which
// can be confusing.
func ScheduleWithAccel(accel float64) ScheduleOpt {
	return func(opts *ScheduleConfig) { opts.Accel = accel }
}

// ScheduleWithSilence disables all logging internal to the scheduler. NOTE: This won't disable logging from Tasks,
// unless the task accesses ScheduleConfig and implements the appropriate logic.
func ScheduleWithSilence() ScheduleOpt {
	return func(opts *ScheduleConfig) { opts.Silent = true }
}

// |||| MATH UTILITIES ||||

func durationGCD(durs ...time.Duration) time.Duration {
	if len(durs) == 0 {
		return time.Duration(0)
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
