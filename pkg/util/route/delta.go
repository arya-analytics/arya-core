package route

// |||| CONTEXT ||||

type DeltaContext[V, C any] struct {
	Outlets map[DeltaOutlet[V, C]]bool
}

// |||| INLET ||||

type DeltaInlet[V, C any] interface {
	Next() <-chan V
	Errors() <-chan error
	Update(DeltaContext[V, C])
}

// |||| OUTLET ||||

type DeltaOutlet[V, C any] interface {
	Context() C
	Send() chan<- V
	SendError() chan<- error
}

// |||| DELTA ||||

func NewDelta[V, C any](inlet DeltaInlet[V, C]) *Delta[V, C] {
	return &Delta[V, C]{
		inlet:        inlet,
		outlets:      make(map[DeltaOutlet[V, C]]bool),
		addOutlet:    make(chan DeltaOutlet[V, C]),
		removeOutlet: make(chan DeltaOutlet[V, C]),
	}
}

type Delta[V, C any] struct {
	inlet        DeltaInlet[V, C]
	outlets      map[DeltaOutlet[V, C]]bool
	addOutlet    chan DeltaOutlet[V, C]
	removeOutlet chan DeltaOutlet[V, C]
}

func (d *Delta[V, C]) Start() {
	for {
		select {
		case o := <-d.addOutlet:
			d.processAddOutlet(o)
		case o := <-d.removeOutlet:
			d.processRemoveOutlet(o)
		case e := <-d.inlet.Errors():
			d.relayError(e)
		case v := <-d.inlet.Next():
			d.relay(v)
		}
	}
}

// || OUTLETS ||

func (d *Delta[V, C]) AddOutlet(o DeltaOutlet[V, C]) {
	d.addOutlet <- o
}

func (d *Delta[V, C]) RemoveOutlet(o DeltaOutlet[V, C]) {
	d.removeOutlet <- o
}

func (d *Delta[V, C]) processAddOutlet(o DeltaOutlet[V, C]) {
	d.outlets[o] = true
	d.updateInlet()
}

func (d *Delta[V, C]) processRemoveOutlet(o DeltaOutlet[V, C]) {
	delete(d.outlets, o)
	d.updateInlet()
}

// || INLET ||

func (d *Delta[V, C]) updateInlet() {
	d.inlet.Update(DeltaContext[V, C]{Outlets: d.outlets})
}

// || RELAY ||

func (d *Delta[V, C]) relay(v V) {
	for outlet := range d.outlets {
		select {
		case outlet.Send() <- v:
		default:
		}
	}
}

func (d *Delta[V, C]) relayError(e error) {
	for outlet := range d.outlets {
		for {
			select {
			case outlet.SendError() <- e:
			default:
			}
		}
	}
}
