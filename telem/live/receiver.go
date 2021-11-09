package live

type Receiver struct {
	send chan Telemetry
}
