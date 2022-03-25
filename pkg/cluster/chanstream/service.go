package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
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
	configs      map[model.PK]*models.ChannelConfig
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
		configs:      map[model.PK]*models.ChannelConfig{},
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

	goeOne := tsquery.NewCreate().Model(rsRfl).BindExec(s.remote.Create).GoExec(ctx)
	goeTwo := tsquery.NewCreate().Model(lsRfl).BindExec(s.localStorage.exec).GoExec(ctx)
	go func() {
		for {
			select {
			case err := <-goeOne.Errors:
				goe.Errors <- err
			case err := <-goeTwo.Errors:
				goe.Errors <- err
			}
		}
	}()
	for {
		sample, sampleOk := p.Model().ChanRecv()
		if !sampleOk {
			break
		}
		cfgPK := model.NewPK(sample.StructFieldByName(csFieldCfgID).Interface().(uuid.UUID))
		cfg, err := s.retrieveConfig(ctx, cfgPK)
		if err != nil {
			goe.Errors <- err
			continue
		}
		sample.StructFieldByName(csFieldCfg).Set(reflect.ValueOf(cfg))
		configNodeIsHostForEachSwitch(sample, lsRfl.ChanSend, rsRfl.ChanSend)
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

	cfg, err := s.retrieveConfig(ctx, pkc[0])
	if err != nil {
		goe.Errors <- err
	}

	q := tsquery.NewRetrieve().Model(p.Model()).WherePK(cfg.ID)
	if cfg.Node.IsHost {
		q = q.BindExec(s.localStorage.exec)
	} else {
		panic("Hello")
		//q = q.BindExec(s.remote)
	}

	qGoe := q.GoExec(ctx)

	go func() {
		for {
			select {
			case err := <-qGoe.Errors:
				goe.Errors <- err
				//	goe.Errors <- err
				//case err := <-goeTwo.Errors:
				//	goe.Errors <- err
			}
		}
	}()
	return nil
}

func (s *Service) retrieveConfig(ctx context.Context, pk model.PK) (*models.ChannelConfig, error) {
	cfg, ok := s.configs[pk]
	if !ok {
		cfg = &models.ChannelConfig{}
		if err := query.NewRetrieve().
			Model(cfg).
			WherePK(pk).
			BindExec(s.localStorage.exec).
			Relation(cfgRelNode, nodeFields()...).
			Exec(ctx); err != nil {
			return cfg, err
		}
		s.configs[pk] = cfg
	}
	return cfg, nil
}

// |||| ROUTING ||||

func configNodeIsHostForEachSwitch(mRfl *model.Reflect, localF, remoteF func(m *model.Reflect)) {
	route.ModelSwitchBoolean(
		mRfl,
		csFieldNodeIsHost,
		func(_ bool, m *model.Reflect) { m.ForEach(func(rfl *model.Reflect, i int) { localF(rfl) }) },
		func(_ bool, m *model.Reflect) { m.ForEach(func(rfl *model.Reflect, i int) { remoteF(rfl) }) },
	)
}

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{&models.ChannelSample{}}
}
