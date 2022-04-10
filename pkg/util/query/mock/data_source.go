package mock

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/filter"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"reflect"
	"strings"
	"sync"
)

type DataSourceMem struct {
	query.Assemble
	mu   sync.RWMutex
	Data *model.DataSource
	*query.HookRunner
}

func NewDataSourceMem() *DataSourceMem {
	ds := &DataSourceMem{Data: model.NewDataSource(), HookRunner: query.NewHookRunner()}
	ds.Assemble = query.NewAssemble(ds.Exec)
	return ds
}

func (s *DataSourceMem) Exec(ctx context.Context, p *query.Pack) error {
	c := query.NewCatch(ctx, p)
	c.Exec(s.Before)
	c.Exec(func(ctx context.Context, p *query.Pack) error {
		return query.Switch(ctx, p, query.Ops{
			&query.Create{}:       s.create,
			&streamq.TSCreate{}:   s.create,
			&query.Retrieve{}:     s.retrieve,
			&streamq.TSRetrieve{}: s.retrieve,
			&query.Update{}:       s.update,
			&query.Delete{}:       s.delete,
		})
	})
	c.Exec(s.After)
	return c.Error()
}

func (s *DataSourceMem) retrieve(ctx context.Context, p *query.Pack) error {
	d := s.Data.Retrieve(p.Model().Type())
	f, err := s.filter(p, d, filter.ErrOnNotFound())
	if err != nil {
		return err
	}
	s.retrieveRelations(p, d)
	if p.Model().IsStruct() {
		p.Model().Set(f.ChainValueByIndex(0))
	} else {
		p.Model().Set(f)
	}
	return nil
}

func (s *DataSourceMem) delete(ctx context.Context, p *query.Pack) error {
	d := s.Data.Retrieve(p.Model().Type())
	f, err := s.filter(p, d)
	if err != nil {
		return err
	}
	newData := d.Filter(func(rfl *model.Reflect, i int) (match bool) {
		f.ForEach(func(nRfl *model.Reflect, i int) {
			if nRfl.PK().Equals(rfl.PK()) {
				match = true
			}
		})
		return !match
	})
	s.Data.Write(newData)
	return nil
}

func (s *DataSourceMem) create(ctx context.Context, p *query.Pack) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data.Retrieve(p.Model().Type()).ChainAppendEach(p.Model())
	return nil
}

func (s *DataSourceMem) update(ctx context.Context, p *query.Pack) error {
	bulk := query.RetrieveBulkUpdateOpt(p)
	if bulk {
		return s.bulkUpdate(p)
	}
	return s.unaryUpdate(p)
	return nil
}

func (s *DataSourceMem) unaryUpdate(p *query.Pack) error {
	if !p.Model().IsStruct() {
		panic("model must be a struct when not using bulk update")
	}
	f, err := s.filter(
		p,
		s.Data.Retrieve(p.Model().Type()),
		filter.ErrOnNotFound(),
		filter.ErrOnMultipleResults(),
	)
	if err != nil {
		return err
	}
	fo, ok := query.RetrieveFieldsOpt(p)
	if !ok {
		panic("fields must be specified for updates")
	}
	u := f.ChainValueByIndex(0)
	for _, fn := range fo {
		fld := u.StructFieldByName(fn)
		if !fld.IsValid() {
			panic(fmt.Sprintf("field %s not found", fn))
		}
		fld.Set(p.Model().StructFieldByName(fn))
	}
	return nil
}

func (s *DataSourceMem) bulkUpdate(p *query.Pack) error {
	d := s.Data.Retrieve(p.Model().Type())
	fo, ok := query.RetrieveFieldsOpt(p)
	if !ok {
		panic("fields must be specified for bulk updates")
	}
	d.ForEach(func(nDRfl *model.Reflect, i int) {
		p.Model().ForEach(func(nSRfl *model.Reflect, i int) {
			if nDRfl.PK().Equals(nSRfl.PK()) {
				for _, f := range fo {
					nDRfl.StructFieldByName(f).Set(nSRfl.StructFieldByName(f))
				}
			}
		})
	})
	return nil
}

func (s *DataSourceMem) filter(p *query.Pack, d *model.Reflect, opts ...filter.Opt) (*model.Reflect, error) {
	s.retrieveWhereFieldRelations(p, d)
	return filter.Filter(p, d, opts...)
}

func (s *DataSourceMem) retrieveRelations(p *query.Pack, d *model.Reflect) {
	ro := query.RetrieveRelationOpts(p)
	for _, r := range ro {
		s.retrieveRelation(d, r)
	}
}

func (s *DataSourceMem) retrieveRelation(sRfl *model.Reflect, rel query.RelationOpt) {
	names := model.SplitFieldNames(rel.Name)
	name := names[0]
	fldT := sRfl.FieldTypeByName(name)
	st, ok := sRfl.StructTagChain().RetrieveByFieldName(name)
	if !ok {
		panic(fmt.Sprintf("field %s couldn't be found on model %s", name, sRfl.Type()))
	}
	if fldT.Kind() == reflect.Ptr {
		fldT = fldT.Elem()
	}
	joinStr, ok := st.Retrieve("model", "join")
	if !ok {
		panic(fmt.Sprintf("model %s must have a join tag specified in order to perform lookup", sRfl.Type().Name()))
	}
	str := strings.Split(joinStr, "=")
	if len(str) != 2 {
		panic(fmt.Sprintf("struct tag join improperlly formatted: %s", joinStr))
	}
	s.Data.Retrieve(fldT).ForEach(func(nDRfl *model.Reflect, i int) {
		dFld := nDRfl.StructFieldByName(str[1])
		sRfl.ForEach(func(nSRfl *model.Reflect, i int) {
			sFld := nSRfl.StructFieldByName(str[0])
			if sFld.Interface() == dFld.Interface() {
				if len(names) > 1 {
					s.retrieveRelation(nDRfl, query.RelationOpt{
						Name:   strings.Join(names[1:], "."),
						Fields: rel.Fields,
					})
				}
				nSRfl.StructFieldByName(names[0]).Set(nDRfl.PointerValue())
			}
		})
	})
}

func (s *DataSourceMem) retrieveWhereFieldRelations(p *query.Pack, sRfl *model.Reflect) {
	wFld, ok := query.RetrieveWhereFieldsOpt(p)
	if !ok {
		return
	}
	for k := range wFld {
		fn, ln := model.SplitLastFieldName(k)
		if fn != "" {
			s.retrieveRelation(sRfl, query.RelationOpt{Name: fn, Fields: query.FieldsOpt{ln}})
		}
	}
}
