package query

type Delete struct {
	where
}

func NewDelete() *Delete {
	d := &Delete{}
	d.baseInit(d)
	return d
}

func (d *Delete) Model(m interface{}) *Delete {
	d.baseModel(m)
	return d
}

func (d *Delete) WherePK(pk interface{}) *Delete {
	newPKOpt(d.Pack(), pk)
	return d
}

func (d *Delete) WherePKs(pks interface{}) *Delete {
	newPKsOpt(d.Pack(), pks)
	return d
}

func (d *Delete) BindExec(e Execute) *Delete {
	d.baseBindExec(e)
	return d
}
