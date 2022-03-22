package chanstream

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

// |||| SERVICE |||

type ServiceRemote interface {
}

type Service struct {
	localExec query.Execute
	remote    ServiceRemote
}

func NewService(local query.Execute, remote ServiceRemote) *Service {
	return &Service{remote: remote, localExec: local}
}

func (s *Service) CanHandle(p *query.Pack) bool {
	if !p.Model().IsChan() {
		panic("chanstream service can't handle non-channel models yet!")
	}
	return catalog().Contains(p.Model())
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{})
}

// Abbreviation reminder:
// cfg - Channel Config
// cs - channel Sample
const (
	cfgRelNode        = "Node"
	csFieldCfgID      = "ChannelConfigID"
	csFieldCfg        = "ChannelConfig"
	csFieldNodeIsHost = "ChannelConfig.Node.IsHost"
	csFieldNode       = "ChannelConfig.Node"
)

func nodeFields() []string {
	return []string{"ID", "Address", "IsHost", "RPCPort"}
}

func retrieveChannelConfigsQuery(p *query.Pack, cfgs interface{}) *query.Retrieve {
	q := query.NewRetrieve().Model(cfgs).Relation(cfgRelNode, nodeFields()...)
	pkc, ok := query.PKOpt(p)
	if !ok {
		q.WherePKs(pkc)
	}
	return q
}

func (s *Service) create(ctx context.Context, p *query.Pack) {
	goExecOpt, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream create queries must be run using goexec")
	}
	var (
		rs       = make(chan *models.ChannelSample)
		ls       = make(chan *models.ChannelSample)
		rsRfl    = model.NewReflect(&rs)
		lsRfl    = model.NewReflect(&ls)
		cfgChain = model.NewReflect(&[]*models.ChannelConfig{})
		c        = errutil.NewCatchContext(ctx, errutil.WithHooks(errutil.NewPipeHook(goExecOpt.Errors)))
	)
	// CLARIFICATION: Retrieves information about the  channel configs and nodes the model belongs to.
	c.Exec(retrieveChannelConfigsQuery(p, cfgChain).BindExec(s.localExec).Exec)
	for {
		sample, sampleOk := p.Model().ChanRecv()
		if !sampleOk {
			break
		}
		cfgPK := model.NewPK(sample.StructFieldByName(csFieldCfgID))
		cfg, cfgOk := cfgChain.ValueByPK(cfgPK)
		if !cfgOk {
			goExecOpt.Errors <- query.NewSimpleError(query.ErrorTypeInvalidArgs, errors.New("invalid config"))
			continue
		}
		sample.StructFieldByName(csFieldCfg).Set(cfg.PointerValue())
		configNodeIsHostSwitch(
			sample,
			func(m *model.Reflect) { lsRfl.ChanSend(sample) },
			func(m *model.Reflect) { rsRfl.ChanSend(sample) },
		)
	}
}

func (s *Service) retrieve() {

}

// |||| ROUTING ||||

func configNodeIsHostSwitch(mRfl *model.Reflect, localF, remoteF func(m *model.Reflect)) {
	route.ModelSwitchBoolean(
		mRfl,
		csFieldNodeIsHost,
		func(_ bool, m *model.Reflect) { localF(m) },
		func(_ bool, m *model.Reflect) { remoteF(m) },
	)
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelSample{}}
}
