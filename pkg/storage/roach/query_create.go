package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/uptrace/bun"
	"reflect"
)

type createQuery struct {
	baseQuery
	q *bun.InsertQuery
}

func newCreate(db *bun.DB) *createQuery {
	r := &createQuery{q: db.NewInsert()}
	r.baseInit()
	return r
}

func (c *createQuery) Model(m interface{}) storage.MDCreateQuery {
	rm := c.baseModel(m)
	c.baseAdaptToDest()
	c.catcher.Exec(func() error {
		beforeInsertSetUUID(rm)
		c.q = c.q.Model(rm.Pointer())
		return nil
	})
	return c
}

func (c *createQuery) Exec(ctx context.Context) error {
	c.catcher.Exec(func() error {
		_, err := c.q.Exec(ctx)
		return err
	})
	return c.baseErr()
}

func (c *createQuery) validateReq(rm interface{}) {
	c.catcher.Exec(func() error { return createReqValidator.Exec(rm) })
}

// |||| VALIDATORS ||||
var createReqValidator = validate.New([]validate.Func{
	validatePK,
})

func validatePK(v interface{}) (err error) {
	rfl := v.(*model.Reflect)
	if rfl.IsChain() {
		for i := 0; i < rfl.ChainValue().Len(); i++ {
			err = validatePK(rfl.ChainValueByIndex(i))
		}
	} else {
		f := rfl.Value().FieldByName(model.KeyPK)
		switch f.Kind() {
		case reflect.Int:
			if f.Interface() == 0 {
				err = storage.Error{Type: storage.ErrTypeNoPK}
			}
		}
	}
	return err
}
