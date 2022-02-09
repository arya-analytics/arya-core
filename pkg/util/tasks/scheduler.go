package tasks

import (
	"context"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

// |||| SCHEDULER ||||

type Scheduler interface {
	Start(ctx context.Context)
	Stop()
	Errors() chan error
}

const defaultAccel = 1

type SchedulerBase struct {
	Tasks    []Task
	cfg      *SchedulerConfig
	chanErr  chan error
	chanStop chan bool
}

func NewBaseScheduler(tasks []Task, opts ...SchedulerOpt) Scheduler {
	s := &SchedulerBase{
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

func (s *SchedulerBase) Start(ctx context.Context) {
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

func (s *SchedulerBase) Stop() {
	s.chanStop <- true
}

func (s *SchedulerBase) Errors() chan error {
	return s.chanErr
}

func (s *SchedulerBase) tickInterval() time.Duration {
	var intervals []time.Duration
	for _, t := range s.Tasks {
		intervals = append(intervals, t.Interval)
	}
	i := s.accelerate(durationGCD(intervals...))
	return i
}

const taskExecThreshold = 50 * time.Millisecond

func (s *SchedulerBase) exec(ctx context.Context, t time.Duration) {
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

func (s *SchedulerBase) accelerate(t time.Duration) time.Duration {
	return time.Duration((t.Minutes() / s.cfg.Accel) * minToNanoSec)
}

func (s *SchedulerBase) bindOpts(opts ...SchedulerOpt) {
	for _, opt := range opts {
		opt(s.cfg)
	}
}

func (s *SchedulerBase) logStart() {
	if s.cfg.Name != "" && !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name":          s.cfg.Name,
			"task_count":    len(s.Tasks),
			"tick_interval": s.tickInterval(),
			"Accel":         s.cfg.Accel,
		}).Infof("Starting %s", s.cfg.Name)
	}
}

func (s *SchedulerBase) logStop() {
	if s.cfg.Name != "" && !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name": s.cfg.Name,
		}).Infof("Stopping %s", s.cfg.Name)

	}
}

func (s *SchedulerBase) logTaskFailure(task Task, err error) {
	if !s.cfg.Silent {
		log.WithFields(log.Fields{
			"Name":      s.cfg.Name,
			"task_name": task.Name,
		}).Errorf("Task failed! %s", err)
	}
}

// |||| SCHEDULER OPTIONS ||||

type SchedulerConfig struct {
	Name   string
	Accel  float64
	Silent bool
}

type SchedulerOpt func(opts *SchedulerConfig)

func ScheduleWithName(name string) SchedulerOpt {
	return func(opts *SchedulerConfig) { opts.Name = name }
}

func ScheduleWithAccel(accel float64) SchedulerOpt {
	return func(opts *SchedulerConfig) { opts.Accel = accel }
}

func ScheduleWithSilence() SchedulerOpt {
	return func(opts *SchedulerConfig) { opts.Silent = true }
}

func durationGCD(durs ...time.Duration) time.Duration {
	if len(durs) == 0 {
		panic("cannot get the duration gcd with no arguments")
	}
	if len(durs) < 2 {
		return durs[0]
	}
	gcd := GCD(durs[0].Nanoseconds(), durs[1].Nanoseconds())
	if len(durs) == 2 {
		return time.Duration(gcd)
	}
	for _, dur := range durs[2:] {
		gcd = GCD(dur.Nanoseconds(), gcd)
	}
	return time.Duration(gcd)
}

func GCD(a, b int64) int64 {
	bigA, bigB := big.NewInt(a), big.NewInt(b)
	gcd := big.NewInt(0)
	return gcd.GCD(nil, nil, bigA, bigB).Int64()
}

// |||| BATCH ||||

type SchedulerBatch struct {
	schedulers []Scheduler
	chanErr    chan error
	chanStop   chan bool
}

func NewBatchScheduler(schedulers ...Scheduler) Scheduler {
	return &SchedulerBatch{schedulers: schedulers}
}

func (sb *SchedulerBatch) Start(ctx context.Context) {
	for _, s := range sb.schedulers {
		go s.Start(ctx)
		go sb.pipeErrors(s.Errors())
	}
}

func (sb *SchedulerBatch) Stop() {
	for _, s := range sb.schedulers {
		s.Stop()
	}
}

func (sb *SchedulerBatch) Errors() chan error {
	return sb.chanErr
}

func (sb *SchedulerBatch) pipeErrors(chanErr chan error) {
	err := <-chanErr
	sb.chanErr <- err
}
