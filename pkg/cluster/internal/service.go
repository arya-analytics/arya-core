package internal

import (
	"context"
)

type ServiceOperation func(ctx context.Context, qr *QueryRequest) error
