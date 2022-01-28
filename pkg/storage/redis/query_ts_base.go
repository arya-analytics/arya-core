package redis

type tsBaseQuery struct {
	baseQuery
}

func (tsb *tsBaseQuery) tsBaseModelWrapper() *tsModelWrapper {
	return &tsModelWrapper{rfl: tsb.modelAdapter.Dest()}
}
