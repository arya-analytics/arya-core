package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/google/uuid"
	"reflect"
)

// |||| SERVICE |||

type Service struct {
	localExec query.Execute
	remote    *ServiceRemoteRPC
	configs   map[uuid.UUID]*models.ChannelConfig
}

func NewService(local query.Execute, remote *ServiceRemoteRPC) *Service {
	return &Service{
		remote:    remote,
		localExec: local,
		configs:   map[uuid.UUID]*models.ChannelConfig{},
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
		&tsquery.Retrieve{}: s.retrieve,
		&tsquery.Create{}:   s.create,
	})
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

func (s *Service) create(ctx context.Context, p *query.Pack) error {
	goExecOpt, ok := tsquery.RetrieveGoExecOpt(p)
	if !ok {
		panic("chanstream create queries must be run using goexec")
	}
	var (
		rs    = make(chan *models.ChannelSample)
		ls    = make(chan *models.ChannelSample)
		rsRfl = model.NewReflect(&rs)
		lsRfl = model.NewReflect(&ls)
	)

	tsquery.NewCreate().Model(rsRfl).BindExec(s.remote.Create).GoExec(ctx, goExecOpt.Errors)
	tsquery.NewCreate().Model(lsRfl).BindExec(s.localExec).GoExec(ctx, goExecOpt.Errors)

	for {
		sample, sampleOk := p.Model().ChanRecv()
		if !sampleOk {
			break
		}
		cfgPK := sample.StructFieldByName(csFieldCfgID).Interface().(uuid.UUID)
		cfg, err := s.retrieveConfig(ctx, cfgPK)
		if err != nil {
			goExecOpt.Errors <- err
			continue
		}
		sample.StructFieldByName(csFieldCfg).Set(reflect.ValueOf(cfg))
		configNodeIsHostSwitch(sample, lsRfl.ChanSend, rsRfl.ChanSend)
	}
	return nil
}

func (s *Service) retrieveConfig(ctx context.Context, pk uuid.UUID) (*models.ChannelConfig, error) {
	cfg, ok := s.configs[pk]
	if !ok {
		cfg = &models.ChannelConfig{}
		if err := query.NewRetrieve().
			Model(cfg).
			WherePK(pk).
			BindExec(s.localExec).
			Relation(cfgRelNode, nodeFields()...).
			Exec(ctx); err != nil {
			return cfg, err
		}
		s.configs[pk] = cfg
	}
	return cfg, nil
}

func (s *Service) retrieve(ctx context.Context, p *query.Pack) error {
	return nil
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
