package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
)

type receiver interface {
	Read() []*models.ChannelSample
}

type senderConfig struct {
	pkc model.PKChain
}

type sender interface {
	Config() senderConfig
}

type Relay struct {
	receivers map[receiver]bool
	sender    map[sender]bool
	dr        telem.DataRate
}

type Locator struct {
	ctx     context.Context
	configs map[model.PK]*models.ChannelConfig
}

func (l *Locator) buildReceivers(cfg []senderConfig) ([]receiver, error) {
	pkc := aggregatePKC(cfg)
	cfgs, err := l.retrieveConfigs(pkc)
	if err != nil {
		return nil, err
	}
}

func (l *Locator) retrieveConfigs(pkc model.PKChain) (configs []*models.ChannelConfig, error) {
	var ccToRetrieve model.PKChain
	for _, pk := range pkc {
		cfg, ok := l.configs[pk]
		if !ok {
			ccToRetrieve = append(ccToRetrieve, pk)
		} else {
			configs = append(configs, cfg)
		}
	}
	retrievedConfigs, err := l.retrieveConfigsFromStorage(ccToRetrieve)
	if err != nil {
		return configs, err
	}
	configs = append(configs, retrievedConfigs...)
	for _, cfg := retrievedConfigs {
		l.configs[model.NewReflect(cfg).PK()] = cfg
	}
	return configs, nil
}

func (l *Locator) retrieveConfigsFromStorage(pkc model.PKChain) (configs []*models.ChannelConfig, err error) {
	return configs, query.NewRetrieve().Model(&configs).WherePKs(pkc).Relation(cfgRelNode, nodeFields()...).Exec(l.ctx)
}

func aggregatePKC(cfgs []senderConfig) model.PKChain {
	pkc := model.NewPKChain([]uuid.UUID{})
	for _, cfg := range cfgs {
		pkc = append(pkc, cfg.pkc...)
	}
	return pkc.Unique()
}

type RemoteReceiver struct {
	pkc  model.PKChain
	data *model.Reflect
}

func (r *RemoteReceiver) Read() *model.Reflect {
	rfl := model.NewReflect(&[]*models.ChannelSample{})
	for {
		nRfl, ok := r.data.ChanRecv()
		if !ok {
			break
		}
		rfl.ChainAppend(nRfl)
	}
	return rfl
}
