package redis

type tsBaseQuery struct {
	baseQuery
}

func (tsb *tsBaseQuery) tsBaseModelWrapper() *TSModelWrapper {
	return &TSModelWrapper{rfl: tsb.modelAdapter.Dest()}
}
