package rng

type Service struct {
	obs Observe
	p   Persist
}

func NewService(obs Observe, p Persist) *Service {
	return &Service{obs: obs, p: p}
}

func (s *Service) NewAllocate() *Allocate {
	return &Allocate{obs: s.obs, pst: s.p}
}
