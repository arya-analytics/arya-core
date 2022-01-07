package live

import (
	"github.com/arya-analytics/aryacore/pkg/telem"
	"github.com/google/uuid"
	"time"
)

// Relay manages the interactions between a set of telemetry 'senders' and 'receivers.'
type Relay struct {
	// Config
	chanConfigs map[int32]map[uuid.UUID]bool
	// Sender Management
	addSender          chan Sender
	removeSender       chan Sender
	updateSenderConfig chan SenderConfig
	senders            map[Sender]bool
	// Receiver Management
	receivers      map[Receiver]bool
	addReceiver    chan Receiver
	removeReceiver chan Receiver
	// Locator
	locator Locator
}

// NewRelay creates a new Relay. Returns a pointer to the created Relay.
func NewRelay(locator Locator) *Relay {
	chanConfigs := map[int32]map[uuid.UUID]bool{}
	addSender := make(chan Sender)
	removeSender := make(chan Sender)
	updateSenderConfig := make(chan SenderConfig)
	senders := map[Sender]bool{}
	receivers := map[Receiver]bool{}
	addReceiver := make(chan Receiver)
	removeReceiver := make(chan Receiver)
	return &Relay{
		chanConfigs,
		addSender,
		removeSender,
		updateSenderConfig,
		senders,
		receivers,
		addReceiver,
		removeReceiver,
		locator,
	}
}

// Start starts the Relay and begins listening to requests and sending responses to
// senders and receivers
func (r *Relay) Start() {
	ticker := time.NewTicker(rate)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case cfg := <-r.updateSenderConfig:
			r.handleConfigUpdate(cfg)
		case sender := <-r.addSender:
			r.senders[sender] = true
		case sender := <-r.removeSender:
			delete(r.senders, sender)
		case receiver := <-r.addReceiver:
			r.receivers[receiver] = true
		case receiver := <-r.removeReceiver:
			delete(r.receivers, receiver)
		case <-ticker.C:
			tlm := r.readFromReceivers()
			r.forwardToSenders(tlm)
		}
	}
}

// readFromReceivers reads, aggregates, and returns telemetry from receivers
func (r *Relay) readFromReceivers() (slc telem.Slice) {
	slc = telem.Slice{}
	for rc := range r.receivers {
		for key, val := range rc.receive() {
			slc[key] = val
		}
	}
	return slc
}

/// forwardToSenders forwards telemetry (p) to senders
func (r *Relay) forwardToSenders(slc telem.Slice) {
	for sr := range r.senders {
		sr.send(slc)
	}
}

func (r *Relay) handleConfigUpdate(cfg SenderConfig) {
	for _, senderIds := range r.chanConfigs {
		delete(senderIds, cfg.ID)

	}
	for _, chanCfg := range cfg.ChanCfgs {
		_, found := r.chanConfigs[chanCfg]
		if !found {
			r.chanConfigs[chanCfg] = map[uuid.UUID]bool{}
		}
		r.chanConfigs[chanCfg][cfg.ID] = true
	}
	for chanCfg, senderIds := range r.chanConfigs {
		if len(senderIds) == 0 {
			delete(r.chanConfigs, chanCfg)
		}
	}
	var chanCfgChain []int32
	for chanCfg := range r.chanConfigs {
		chanCfgChain = append(chanCfgChain, chanCfg)
	}
	receivers := r.locator.Locate(chanCfgChain)
	for _, r := range receivers {
		go r.Start(chanCfgChain)
	}
}
