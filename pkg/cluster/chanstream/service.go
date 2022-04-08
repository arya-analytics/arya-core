package chanstream

import "C"
import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

// |||| SERVICE |||

type Service struct {
	local  *LocalStorage
	remote *RemoteRPC
	ccMemo *query.Memo
}

func NewService(qExec query.Execute, remote *RemoteRPC) *Service {
	memo := query.NewMemo(&[]*models.ChannelConfig{})
	return &Service{remote: remote, local: NewLocalStorage(qExec), ccMemo: memo}
}

// CanHandle implements cluster.Service.
func (s *Service) CanHandle(p *query.Pack) bool {
	if !p.Model().IsChan() {
		panic("telemstream service can't handle non-channel models yet!")
	}
	return catalog().Contains(p.Model())
}

// Exec implements cluster.Service.
func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   s.tsCreate,
		&streamq.TSRetrieve{}: s.tsRetrieve,
	})
}

const (
	cfgRelNode         = "Node"
	CfgFieldNodeIsHost = "Node.IsHost"
	csFieldNode        = "ChannelConfig.Node"
)

func nodeFields() []string {
	return []string{"ID", "Address", "IsHost", "RPCPort"}
}

func (s *Service) tsCreate(ctx context.Context, p *query.Pack) error {
	c := *query.ConcreteModel[*chan *models.ChannelSample](p)
	remoteStream, localStream, st, cancel, err := s.openTSCreateQueries(ctx, p)
	if err != nil {
		return err
	}

	st.Segment(func() {
		defer func() {
			cancel()
			close(remoteStream)
			close(localStream)
		}()
		for sa := range c {
			if route.CtxDone(ctx) {
				break
			}
			cc, err := s.retrieveConfigsQuery(ctx, sa.ChannelConfigID)
			if err != nil {
				st.Errors <- err
				continue
			}
			sa.ChannelConfig = cc[0]
			if sa.ChannelConfig.Node.IsHost {
				remoteStream <- sa
			} else {
				localStream <- sa
			}
		}
	}, streamq.WithSegmentName("telemstream.Service.tsCreate"))
	return nil
}

func (s *Service) openTSCreateQueries(ctx context.Context, p *query.Pack) (rs, ls chan *models.ChannelSample, st *streamq.Stream, cancel context.CancelFunc, err error) {
	st, _ = streamq.RetrieveStreamOpt(p, query.RequireOpt())
	rs = make(chan *models.ChannelSample)
	ls = make(chan *models.ChannelSample)
	bCtx, cancel := context.WithCancel(ctx)
	_, err = streamq.NewTSCreate().Model(&rs).BindStream(st).BindExec(s.local.exec).Stream(bCtx)
	if err != nil {
		cancel()
		return nil, nil, nil, nil, err
	}
	_, err = streamq.NewTSCreate().Model(&ls).BindStream(st).BindExec(s.remote.exec).Stream(bCtx)
	if err != nil {
		cancel()
		return nil, nil, nil, nil, err
	}
	return rs, ls, st, cancel, nil
}

func (s *Service) tsRetrieve(ctx context.Context, p *query.Pack) error {
	st, _ := streamq.RetrieveStreamOpt(p, query.RequireOpt())
	pkc, _ := query.PKOpt(p, query.RequireOpt())
	cc, err := s.retrieveConfigsQuery(ctx, pkc)
	if err != nil {
		return err
	}

	c := errutil.NewCatchContext(ctx)
	route.ModelSwitchBoolean(
		&cc,
		CfgFieldNodeIsHost,
		func(m *model.Reflect) {
			c.Exec(func(ctx context.Context) (err error) {
				_, err = retrieveSamplesQuery(p, m).BindStream(st).BindExec(s.local.exec).Stream(ctx)
				return err
			})
		},
		func(m *model.Reflect) {
			c.Exec(func(ctx context.Context) (err error) {
				_, err = retrieveSamplesQuery(p, m).BindStream(st).BindExec(s.remote.exec).Stream(ctx)
				return err
			})
		},
	)

	return c.Error()
}

func (s *Service) retrieveConfigsQuery(ctx context.Context, pks interface{}) (cc []*models.ChannelConfig, err error) {
	return cc, query.NewRetrieve().
		Model(&cc).
		WherePKs(pks).
		WithMemo(s.ccMemo).
		Relation(cfgRelNode, nodeFields()...).
		BindExec(s.local.exec).Exec(ctx)
}

// |||| QUERY |||

// || SAMPLES ||

func retrieveSamplesQuery(p *query.Pack, m *model.Reflect) *streamq.TSRetrieve {
	q := streamq.NewTSRetrieve().Model(p.Model()).WherePKs(m.PKChain())
	newNodeOpt(
		q.Pack(),
		m.FieldsByName(cfgRelNode).ToReflect().RawValue().Interface().([]*models.Node),
	)
	return q
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelSample{}}
}

const nodeOptKey query.OptKey = "chanStreamNode"

func newNodeOpt(p *query.Pack, nodes []*models.Node) {
	p.SetOpt(nodeOptKey, nodes)
}

func nodeOpt(p *query.Pack) []*models.Node {
	n, _ := p.RetrieveOpt(nodeOptKey, query.RequireOpt())
	return n.([]*models.Node)
}
