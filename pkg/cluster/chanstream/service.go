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
	"github.com/google/uuid"
	"reflect"
)

// |||| SERVICE |||

type Service struct {
	localStorage *LocalStorage
	remote       *RemoteRPC
	ccMemo       *query.Memo
}

func NewService(local query.Execute, remote *RemoteRPC) *Service {
	relay := &localRelay{
		ctx: context.Background(),
		qe:  local,
		dr:  telem.DataRate(100),
		add: make(chan *query.Pack),
	}

	go relay.start()

	l := &LocalStorage{relay: relay, qe: local}

	return &Service{
		remote:       remote,
		localStorage: l,
		ccMemo:       query.NewMemo(model.NewReflect(&[]*models.ChannelConfig{})),
	}
}

func (s *Service) CanHandle(p *query.Pack) bool {
	if !p.Model().IsChan() {
		panic("chanstream service can't handle non-channel models yet!")
	}
	return catalog().Contains(p.Model())
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&tsquery.Create{}:   s.create,
		&tsquery.Retrieve{}: s.retrieve,
	})
}

// Abbreviation reminder:
// cfg - Channel Config
// cs - channel Sample
const (
	cfgRelNode        = "Node"
	cfgNodeIsHost     = "Node.IsHost"
	csFieldCfgID      = "ChannelConfigID"
	csFieldCfg        = "ChannelConfig"
	csFieldNodeIsHost = "ChannelConfig.Node.IsHost"
	csFieldNode       = "ChannelConfig.Node"
)

func nodeFields() []string {
	return []string{"ID", "Address", "IsHost", "RPCPort"}
}

func (s *Service) create(ctx context.Context, p *query.Pack) error {
	goe, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream create queries must be run using goexec")
	}
	var (
		rs    = make(chan *models.ChannelSample)
		ls    = make(chan *models.ChannelSample)
		rsRfl = model.NewReflect(&rs)
		lsRfl = model.NewReflect(&ls)
	)

	goeOne := tsquery.NewCreate().Model(rsRfl).BindExec(s.remote.exec).GoExec(ctx)
	goeTwo := tsquery.NewCreate().Model(lsRfl).BindExec(s.localStorage.exec).GoExec(ctx)
	go errutil.NewDelta(goe.Errors, goeTwo.Errors, goeOne.Errors).Exec()

	for {
		sample, sampleOk := p.Model().ChanRecv()
		if !sampleOk {
			break
		}
		cfgPK := sample.StructFieldByName(csFieldCfgID).Interface().(uuid.UUID)
		cc, err := s.retrieveConfig(ctx, model.NewPKChain([]uuid.UUID{cfgPK}))
		if err != nil {
			goe.Errors <- err
			continue
		}
		sample.StructFieldByName(csFieldCfg).Set(reflect.ValueOf(cc[0]))
		nodeIsHostForEachSwitch(sample, csFieldNodeIsHost, lsRfl.ChanSend, rsRfl.ChanSend)
	}
	return nil
}

func (s *Service) retrieve(ctx context.Context, p *query.Pack) error {
	goe, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream queries must be run using goexec")
	}

	pkc, ok := query.PKOpt(p)
	if !ok {
		panic("chanstream queries require a pk")
	}

	cc, err := s.retrieveConfig(ctx, pkc)
	if err != nil {
		return err
	}

	var (
		le = tsquery.GoExecOpt{Errors: make(chan error)}
		re = tsquery.GoExecOpt{Errors: make(chan error)}
	)

	nodeIsHostSwitch(
		model.NewReflect(&cc),
		cfgNodeIsHost,
		func(m *model.Reflect) { le = retrieveSamplesQuery(p, m).BindExec(s.localStorage.exec).GoExec(ctx) },
		func(m *model.Reflect) { re = retrieveSamplesQuery(p, m).BindExec(s.remote.exec).GoExec(ctx) },
	)

	go errutil.NewDelta(goe.Errors, le.Errors, re.Errors).Exec()

	return nil
}

func retrieveSamplesQuery(p *query.Pack, m *model.Reflect) *tsquery.Retrieve {
	return tsquery.NewRetrieve().Model(p.Model()).WherePKs(m.PKChain())
}

func (s *Service) retrieveConfig(ctx context.Context, pkc model.PKChain) (cc []*models.ChannelConfig, err error) {
	return cc, query.NewRetrieve().
		Model(&cc).
		WherePKs(pkc).
		WithMemo(s.ccMemo).
		Relation(cfgRelNode, nodeFields()...).
		BindExec(s.localStorage.exec).
		Exec(ctx)
}

// |||| QUERY |||

// |||| ROUTING ||||

func nodeIsHostForEachSwitch(mRfl *model.Reflect, fld string, localF, remoteF func(m *model.Reflect)) {
	nodeIsHostSwitch(
		mRfl,
		fld,
		func(m *model.Reflect) { m.ForEach(func(rfl *model.Reflect, i int) { localF(rfl) }) },
		func(m *model.Reflect) { m.ForEach(func(rfl *model.Reflect, i int) { remoteF(rfl) }) },
	)
}

func nodeIsHostSwitch(mRfl *model.Reflect, fld string, localF, remoteF func(m *model.Reflect)) {
	route.ModelSwitchBoolean(
		mRfl,
		fld,
		func(_ bool, m *model.Reflect) { localF(m) },
		func(_ bool, m *model.Reflect) { remoteF(m) },
	)
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelSample{}}
}
