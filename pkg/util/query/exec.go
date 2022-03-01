package query

import "context"

type Execute func(ctx context.Context, p *Pack) error
