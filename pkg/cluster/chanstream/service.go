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
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"reflect"
)

// |||| SERVICE |||

type Service struct {
	local  *LocalStorage
	remote *RemoteRPC
	ccMemo *query.Memo
}

func NewService(local query.Execute, remote *RemoteRPC) *Service {
	r := &localRelay{
		ctx: context.Background(),
		qe:  local,
		dr:  telem.DataRate(100),
		add: make(chan *query.Pack),
	}
	go r.start()

	l := &LocalStorage{relay: r, qe: local}
	memo := query.NewMemo(model.NewReflect(&[]*models.ChannelConfig{}))
	return &Service{remote: remote, local: l, ccMemo: memo}
}

// CanHandle implements cluster.Service.
func (s *Service) CanHandle(p *query.Pack) bool {
	if !p.Model().IsChan() {
		panic("chanstream service can't handle non-channel models yet!")
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
	csFieldCfgID       = "ChannelConfigID"
	csFieldCfg         = "ChannelConfig"
	csFieldNodeIsHost  = "ChannelConfig.Node.IsHost"
	csFieldNode        = "ChannelConfig.Node"
)

func nodeFields() []string {
	return []string{"ID", "Address", "IsHost", "RPCPort"}
}

func (s *Service) tsCreate(ctx context.Context, p *query.Pack) error {
	var (
		st    = stream(p)
		rs    = make(chan *models.ChannelSample)
		ls    = make(chan *models.ChannelSample)
		rsRfl = model.NewReflect(&rs)
		lsRfl = model.NewReflect(&ls)
	)

	_, err := streamq.NewTSCreate().Model(lsRfl).BindStream(st).BindExec(s.local.exec).Stream(ctx)
	if err != nil {
		return err
	}
	_, err = streamq.NewTSCreate().Model(rsRfl).BindStream(st).BindExec(s.remote.exec).Stream(ctx)
	if err != nil {
		return err
	}

	st.Segment(func() {
		for {
			sample, sampleOk := p.Model().ChanRecv()
			if !sampleOk {
				break
			}
			pkc := sample.FieldsByName(csFieldCfgID).ToPKChain()
			cc, err := s.retrieveConfigs(ctx, pkc)
			if err != nil {
				st.Errors <- err
				continue
			}
			sample.StructFieldByName(csFieldCfg).Set(reflect.ValueOf(cc[0]))
			sampleNodeIsHostSwitchEach(sample, lsRfl.ChanSend, rsRfl.ChanSend)
		}
	})

	return nil
}

func (s *Service) tsRetrieve(ctx context.Context, p *query.Pack) error {
	var (
		goe = stream(p)
		pkc = pkOpt(p)
		le  = &streamq.Stream{Errors: make(chan error)}
		re  = &streamq.Stream{Errors: make(chan error)}
	)

	cc, err := s.retrieveConfigs(ctx, pkc)
	if err != nil {
		return err
	}

	c := errutil.NewCatchContext(ctx)
	route.ModelSwitchBoolean(
		model.NewReflect(&cc),
		CfgFieldNodeIsHost,
		func(m *model.Reflect) {
			c.Exec(func(ctx context.Context) (err error) {
				le, err = retrieveSamplesQuery(p, m).BindExec(s.local.exec).Stream(ctx)
				return err
			})
		},
		func(m *model.Reflect) {
			c.Exec(func(ctx context.Context) (err error) {
				re, err = retrieveSamplesQuery(p, m).BindExec(s.remote.exec).Stream(ctx)
				return err
			})
		},
	)

	errutil.NewDelta(goe.Errors, le.Errors, re.Errors).Exec()

	return c.Error()
}

func (s *Service) retrieveConfigs(ctx context.Context, pkc model.PKChain) (cc []*models.ChannelConfig, err error) {
	return cc, query.NewRetrieve().
		Model(&cc).
		WherePKs(pkc).
		WithMemo(s.ccMemo).
		Relation(cfgRelNode, nodeFields()...).
		BindExec(s.local.exec).
		Exec(ctx)
}

// |||| QUERY |||

// || OPT ||

func stream(p *query.Pack) *streamq.Stream {
	stream, ok := streamq.StreamOpt(p)
	if !ok {
		panic("chanstream queries must be run using goexec")
	}
	return stream
}

func pkOpt(p *query.Pack) model.PKChain {
	pkc, ok := query.PKOpt(p)
	if !ok {
		panic("chanstream queries require a pk")
	}
	return pkc

}

// || SAMPLES ||

func retrieveSamplesQuery(p *query.Pack, m *model.Reflect) *streamq.TSRetrieve {
	q := streamq.NewTSRetrieve().Model(p.Model()).WherePKs(m.PKChain())
	newNodeOpt(
		q.Pack(),
		m.FieldsByName(cfgRelNode).ToReflect().RawValue().Interface().([]*models.Node),
	)
	return q
}

// |||| ROUTING ||||

func sampleNodeIsHostSwitchEach(mRfl *model.Reflect, localF, remoteF func(m *model.Reflect)) {
	route.ModelSwitchBoolean(mRfl,
		csFieldNodeIsHost,
		func(m *model.Reflect) { m.ForEach(func(rfl *model.Reflect, i int) { localF(rfl) }) },
		func(m *model.Reflect) { m.ForEach(func(rfl *model.Reflect, i int) { remoteF(rfl) }) },
	)
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
	n, ok := p.RetrieveOpt(nodeOptKey)
	if !ok {
		panic("node opt not specified. this is a bug")
	}
	return n.([]*models.Node)
}
