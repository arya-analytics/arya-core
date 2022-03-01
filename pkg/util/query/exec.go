package query

import "context"

type Execute interface {
	Exec(ctx context.Context, q *Pack) error
}
