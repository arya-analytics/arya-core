package live

import (
	"github.com/arya-analytics/aryacore/telem"
	"github.com/google/uuid"
	"time"
)

type SenderConfig struct{}

// Relay manages the interactions between a set of telemetry 'senders' and 'receivers.'
type Relay struct {
	// Config
	chanConfigs map[int][]uuid.UUID
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
	chanConfigs := map[int][]uuid.UUID{}
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
func (h *Relay) Start() {
	ticker := time.NewTicker(rate)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case cfg := <-h.updateSenderConfig:
			h.handleConfigUpdate(cfg)
		case sender := <-h.addSender:
			h.senders[sender] = true
		case sender := <-h.removeSender:
			delete(h.senders, sender)
		case receiver := <-h.addReceiver:
			h.receivers[receiver] = true
		case receiver := <-h.removeReceiver:
			delete(h.receivers, receiver)
		case <-ticker.C:
			tlm := h.readFromReceivers()
			h.forwardToSenders(tlm)
		}
	}
}

// readFromReceivers reads, aggregates, and returns telemetry from receivers
func (h *Relay) readFromReceivers() (slc telem.Slice) {
	slc = telem.Slice{}
	for rc := range h.receivers {
		for key, val := range rc.receive() {
			slc[key] = val
		}
	}
	return slc
}

/// forwardToSenders forwards telemetry (p) to senders
func (h *Relay) forwardToSenders(slc telem.Slice) {
	for sr := range h.senders {
		sr.send(slc)
	}
}

func (h *Relay) handleConfigUpdate(cfg SenderConfig) {

}
