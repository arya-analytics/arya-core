package tasks

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	log "github.com/sirupsen/logrus"
	"time"
)

type SchedulerOpt func(opts *schedulerOpts)

func SchedulerWithName(name string) SchedulerOpt {
	return func(opts *schedulerOpts) {
		opts.name = name
	}
}

type schedulerOpts struct {
	name string
}

type Scheduler struct {
	Tasks        []Task
	TickInterval time.Duration
	Errors       chan error
	catcher      *errutil.Catcher
	opts         *schedulerOpts
}

func NewScheduler(tasks []Task, tickInterval time.Duration, opts ...SchedulerOpt) *Scheduler {
	errs := make(chan error)
	catcher := &errutil.Catcher{}
	s := &Scheduler{
		Tasks:        tasks,
		TickInterval: tickInterval,
		Errors:       errs,
		catcher:      catcher,
		opts:         &schedulerOpts{},
	}
	s.bindOpts(opts...)
	return s
}

func (s *Scheduler) logStart() {
	log.WithFields(log.Fields{
		"name":       s.opts.name,
		"task_count": len(s.Tasks),
		"interval":   s.TickInterval,
	}).Infof("Starting %s", s.opts.name)
}

func (s *Scheduler) logTaskFailure(task Task, err error) {
	log.WithFields(log.Fields{
		"name":      s.opts.name,
		"task_name": task.Name,
	}).Errorf("Task failed! %s", err)

}

func (s *Scheduler) Start(ctx context.Context) {
	if s.opts.name != "" {
		s.logStart()
	}

	t0 := time.Now()
	t := time.NewTicker(s.TickInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			break
		case ct := <-t.C:
			if err := s.exec(ctx, ct.Sub(t0)); err != nil {
				s.Errors <- err
			}
		}
	}
}

const taskExecThreshold = 20 * time.Millisecond

func (s *Scheduler) exec(ctx context.Context, t time.Duration) error {
	for _, task := range s.Tasks {
		if t%task.Interval < taskExecThreshold {
			s.catcher.Exec(func() error {
				if err := task.Action(ctx); err != nil {
					s.logTaskFailure(task, err)
					return err
				}
				return nil
			})
		}
	}
	return s.catcher.Error()
}

func (s *Scheduler) bindOpts(opts ...SchedulerOpt) {
	for _, opt := range opts {
		opt(s.opts)
	}
}
