package server

type EndpointBuilder struct {
	bases []string
}

func (eb EndpointBuilder) Base() string {
	return eb.Build()
}

func (eb EndpointBuilder) Build(args ...string) (e string) {
	for _, arg := range append(eb.bases, args...) {
		e += appendSlash(arg)
	}
	return e
}

func (eb EndpointBuilder) Child(args ...string) *EndpointBuilder {
	return NewEndpointBuilder(append(eb.bases, args...)...)
}

func NewEndpointBuilder(bases ...string) *EndpointBuilder {
	return &EndpointBuilder{bases}
}

func appendSlash(e string) (formattedE string) {
	if e[len(e)-1] != '/' {
		return e + "/"
	}
	return e
}
