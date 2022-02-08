package tasks

import (
	"context"
	"time"
)

type Action func(ctx context.Context) error

type Task struct {
	Interval time.Duration
	Action   Action
	Name     string
}
