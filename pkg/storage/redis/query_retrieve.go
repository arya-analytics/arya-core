package redis

//type retrieveQuery struct {
//	baseQuery
//	PKs    []model.PK
//	fromTS time.Time
//	toTS   time.Time
//}
//
//func newRetrieve(client *redistimeseries.Client) *retrieveQuery {
//	c := &retrieveQuery{}
//	c.baseInit(client)
//	return c
//}
//
//func (r *retrieveQuery) Model(m interface{}) storage.CacheRetrieveQuery {
//	r.baseModel(m)
//	r.baseAdaptToDest()
//	return r
//}
//
//func (r *retrieveQuery) WherePK(pk interface{}) storage.CacheRetrieveQuery {
//	r.PKs = append(r.PKs, model.NewPK(pk))
//	return r
//}
//
//func (r *retrieveQuery) WherePKs(pks interface{}) storage.CacheRetrieveQuery {
//	rfl := reflect.ValueOf(pks)
//	for i := 0; i < rfl.Len(); i++ {
//		r.WherePK(rfl.Index(i).Interface())
//	}
//	return r
//}
//
//func (r *retrieveQuery) WhereTimeRange(fromTS time.Time,
//	toTS time.Time) storage.CacheRetrieveQuery {
//	if fromTS.IsZero() {
//	} else {
//		r.fromTS = fromTS
//	}
//	if toTS.IsZero() {
//		r.toTS = time.Unix(redistimeseries.TimeRangeMaximum, 0)
//	} else {
//		r.toTS = toTS
//	}
//	return r
//}
//
//func (r *retrieveQuery) Exec(ctx context.Context) error {
//	dRfl := r.modelAdapter.Dest()
//	switch dRfl.Type() {
//	case reflect.TypeOf(&ChannelSample{}):
//		for _, pk := range r.PKs {
//			if !r.toTS.IsZero() {
//				r.catcher.Exec(func() error {
//					dv, err := r.baseClient().
//						RangeWithOptions(pk.
//							String(),
//							r.fromTS.Unix(), r.toTS.Unix(),
//							redistimeseries.RangeOptions{})
//					return err
//				})
//			}
//
//		}
//	}
//	return r.baseErr()
//}
//
//type retrieveQueryBuilder struct {
//}
