package errutil

type Delta struct {
	inlet  []chan error
	outlet chan error
}

func NewDelta(outlet chan error, inlet ...chan error) *Delta {
	return &Delta{inlet: inlet, outlet: outlet}
}

// Exec pipes errors from inlet to outlet.
func (d *Delta) Exec() {
	for _, inlet := range d.inlet {
		go func(inlet chan error) {
			for err := range inlet {
				d.outlet <- err
			}
		}(inlet)
	}
}
