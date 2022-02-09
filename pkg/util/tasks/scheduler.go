package tasks

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
)

// |||| SCHEDULER ||||

const defaultAccel = 1

type Scheduler struct {
	Tasks        []Task
	TickInterval time.Duration
	Errors       chan error
	stop         chan bool
	opts         *schedulerOpts
}

func NewScheduler(tasks []Task, tickInterval time.Duration, opts ...SchedulerOpt) *Scheduler {
	s := &Scheduler{
		Tasks:        tasks,
		TickInterval: tickInterval,
		Errors:       make(chan error),
		stop:         make(chan bool),
		opts: &schedulerOpts{
			accel: defaultAccel,
		},
	}
	s.bindOpts(opts...)
	return s
}

func (s *Scheduler) Stop() {
	s.stop <- true
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logStart()
	t, t0 := time.NewTicker(s.accelerate(s.TickInterval)), time.Now()
	defer t.Stop()
	for {
		select {
		case ct := <-t.C:
			if err := s.exec(ctx, ct.Sub(t0)); err != nil {
				s.Errors <- err
			}
		case <-s.stop:
			s.logStop()
			close(s.Errors)
			close(s.stop)
			return
		case <-ctx.Done():
			s.logStop()
			return
		}
	}
}

const taskExecThreshold = 50 * time.Millisecond

func (s *Scheduler) exec(ctx context.Context, t time.Duration) error {
	for _, task := range s.Tasks {
		if t%s.accelerate(task.Interval) < taskExecThreshold {
			if err := task.Action(ctx); err != nil {
				s.logTaskFailure(task, err)
				return err
			}
		}
	}
	return nil
}

const minToNanoSec = 1000000000 * 60

func (s *Scheduler) accelerate(t time.Duration) time.Duration {
	return time.Duration((t.Minutes() / s.opts.accel) * minToNanoSec)
}

func (s *Scheduler) bindOpts(opts ...SchedulerOpt) {
	for _, opt := range opts {
		opt(s.opts)
	}
}

func (s *Scheduler) logStart() {
	if s.opts.name != "" && !s.opts.silent {
		log.WithFields(log.Fields{
			"name":       s.opts.name,
			"task_count": len(s.Tasks),
			"interval":   s.TickInterval,
			"accel":      s.opts.accel,
		}).Infof("Starting %s", s.opts.name)
	}
}

func (s *Scheduler) logStop() {
	if s.opts.name != "" && !s.opts.silent {
		log.WithFields(log.Fields{
			"name": s.opts.name,
		}).Infof("Stopping %s", s.opts.name)

	}
}

func (s *Scheduler) logTaskFailure(task Task, err error) {
	if !s.opts.silent {
		log.WithFields(log.Fields{
			"name":      s.opts.name,
			"task_name": task.Name,
		}).Errorf("Task failed! %s", err)
	}
}

// |||| SCHEDULER OPTIONS ||||

type schedulerOpts struct {
	name   string
	accel  float64
	silent bool
}

type SchedulerOpt func(opts *schedulerOpts)

func ScheduleWithName(name string) SchedulerOpt {
	return func(opts *schedulerOpts) { opts.name = name }
}

func ScheduleWithAccel(accel float64) SchedulerOpt {
	return func(opts *schedulerOpts) { opts.accel = accel }
}

func ScheduleWithSilence() SchedulerOpt {
	return func(opts *schedulerOpts) { opts.silent = true }
}
