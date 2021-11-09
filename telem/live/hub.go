package live

import (
	"github.com/arya-analytics/aryacore/ds"
	"github.com/google/uuid"
	"time"
)

type TelemetryHub struct {
	cp *ds.ConnPooler
	senders map[*Sender] bool
	receivers map[*Receiver] bool
	senderConfig map[uuid.UUID] SenderConfig
	addSender chan *Sender
	removeSender chan *Sender
	updateConfig chan *SenderConfig
}

func (th *TelemetryHub) start() {
	ticker := time.NewTicker(receiverUpdateRate)
	for {
		select {
			case sender := <-th.addSender:
				th.senders[sender] = true
			case sender := <-th.removeSender:
				delete(th.senders, sender)
			case config := <-th.updateConfig:
				th.onUpdateConfig(config)
			case <- ticker.C:
				tlm := th.pullTelemetryFromReceivers()
				th.pushTelemetryToSenders(tlm)
		}
	}
}

func (th *TelemetryHub) onUpdateConfig(config *SenderConfig) {

}

func (th *TelemetryHub) pullTelemetryFromReceivers() Telemetry {
	var aggregatedTlm = Telemetry{}
	for receiver := range th.receivers {
		select {
			case tlm := <- receiver.send:
				for k, v := range tlm {
					aggregatedTlm[k] = v
				}
			default:
		}
	}
	return aggregatedTlm
}

func (th *TelemetryHub) pushTelemetryToSenders(tlm Telemetry) {
	for sender := range th.senders {
		cfg := th.senderConfig[sender.id]
		senderTlm := Telemetry{}
		for i := range cfg.channelConfigs {
			channelConfig := cfg.channelConfigs[i]
			if val, ok := tlm[channelConfig.ID]; ok {
				senderTlm[channelConfig.ID] = val
			}
		}
		sender.send <- senderTlm
	}
}