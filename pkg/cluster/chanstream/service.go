package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
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
		&tsquery.Create{}:   s.tsCreate,
		&tsquery.Retrieve{}: s.tsRetrieve,
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
		goe   = goExecOpt(p)
		rs    = make(chan *models.ChannelSample)
		ls    = make(chan *models.ChannelSample)
		rsRfl = model.NewReflect(&rs)
		lsRfl = model.NewReflect(&ls)
	)

	le := tsquery.NewCreate().Model(lsRfl).BindExec(s.local.exec).GoExec(ctx)
	re := tsquery.NewCreate().Model(rsRfl).BindExec(s.remote.exec).GoExec(ctx)
	go errutil.NewDelta(goe.Errors, le.Errors, re.Errors).Exec()

	for {
		sample, sampleOk := p.Model().ChanRecv()
		if !sampleOk {
			break
		}
		pkc := sample.FieldsByName(csFieldCfgID).ToPKChain()
		cc, err := s.retrieveConfigs(ctx, pkc)
		if err != nil {
			goe.Errors <- err
			continue
		}
		sample.StructFieldByName(csFieldCfg).Set(reflect.ValueOf(cc[0]))
		sampleNodeIsHostSwitchEach(sample, lsRfl.ChanSend, rsRfl.ChanSend)
	}
	return nil
}

func (s *Service) tsRetrieve(ctx context.Context, p *query.Pack) error {
	var (
		goe = goExecOpt(p)
		pkc = pkOpt(p)
		le  = tsquery.GoExecOpt{Errors: make(chan error)}
		re  = tsquery.GoExecOpt{Errors: make(chan error)}
	)

	cc, err := s.retrieveConfigs(ctx, pkc)
	if err != nil {
		return err
	}

	route.ModelSwitchBoolean(
		model.NewReflect(&cc),
		CfgFieldNodeIsHost,
		func(m *model.Reflect) {
			le = retrieveSamplesQuery(p, m).BindExec(s.local.exec).GoExec(ctx)
		},
		func(m *model.Reflect) { re = retrieveSamplesQuery(p, m).BindExec(s.remote.exec).GoExec(ctx) },
	)

	go errutil.NewDelta(goe.Errors, le.Errors, re.Errors).Exec()

	return nil
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

func goExecOpt(p *query.Pack) tsquery.GoExecOpt {
	goe, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream queries must be run using goexec")
	}
	return goe
}

func pkOpt(p *query.Pack) model.PKChain {
	pkc, ok := query.PKOpt(p)
	if !ok {
		panic("chanstream queries require a pk")
	}
	return pkc

}

// || SAMPLES ||

func retrieveSamplesQuery(p *query.Pack, m *model.Reflect) *tsquery.Retrieve {
	q := tsquery.NewRetrieve().Model(p.Model()).WherePKs(m.PKChain())
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
