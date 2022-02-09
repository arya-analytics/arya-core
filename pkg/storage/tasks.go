package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/tasks"
	log "github.com/sirupsen/logrus"
)

func startTaskRunner(ctx context.Context, s *storage, opts ...tasks.SchedulerOpt) {
	if s.taskRunner != nil {
		log.Warn("Storage task runner already started. Can't start again.")
	}
	log.Info("Starting storage task runner")
	schedulers := []*tasks.Scheduler{
		s.cfg.EngineMD.NewTasks(s.adapter(s.cfg.EngineMD), opts...),
	}
	s.taskRunner = &taskRunner{errors: make(chan error), schedulers: schedulers,
		_stop: make(chan bool)}
	s.taskRunner.start(ctx)
}

func stopTaskRunner(s *storage) {
	if s.taskRunner == nil {
		log.Warn("Storage task runner hasn't been started! Can't stop.")
		return
	}
	s.taskRunner.stop()
}

type taskRunner struct {
	schedulers []*tasks.Scheduler
	errors     chan error
	_stop      chan bool
}

func (tr *taskRunner) start(ctx context.Context) {
	for _, s := range tr.schedulers {
		go s.Start(ctx)
		go tr.pipeErrors(s.Errors)
	}
	go tr.listenForErrors()
}

func (tr *taskRunner) stop() {
	log.Info("Stopping storage task runner")
	for _, s := range tr.schedulers {
		s.Stop()
	}
	tr._stop <- true
}

func (tr *taskRunner) pipeErrors(errChan chan error) {
	err := <-errChan
	if err != nil {
		tr.errors <- err
	}
}

func (tr *taskRunner) listenForErrors() {
	for {
		select {
		case <-tr._stop:
			close(tr._stop)
			close(tr.errors)
			return
		case err := <-tr.errors:
			log.Errorf("Storage task failed. Error %v", err)
		}
	}
}
