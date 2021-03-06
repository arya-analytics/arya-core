package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
)

type LocalStorage struct {
	delta *route.Delta[*models.ChannelSample, outletContext]
	qe    query.Execute
}

func NewLocalStorage(qe query.Execute) *LocalStorage {
	dr := telem.DataRate(100)
	di := &localDeltaInlet{
		dr:        dr,
		ctx:       context.Background(),
		qExec:     qe,
		errC:      make(chan error, 10),
		valStream: make(chan *models.ChannelSample, 1),
	}
	go di.start()
	d := route.NewDelta[*models.ChannelSample, outletContext](di)
	go d.Start()
	return &LocalStorage{delta: d, qe: qe}
}

func (ls *LocalStorage) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   newLocalStreamCreate(ls.qe).exec,
		&streamq.TSRetrieve{}: newLocalStreamRetrieve(ls.delta).exec,
	}, query.SwitchWithDefault(ls.qe))
}
