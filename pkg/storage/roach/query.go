package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/uptrace/bun"
)

// |||| CREATE ||||

type base struct {
	exchange *model.Exchange
	sqlGen   sqlGen
	db       *bun.DB
}

type create struct {
	base
	bunQ *bun.InsertQuery
}

func newCreate(db *bun.DB) *create {
	return &create{bunQ: db.NewInsert(), base: base{db: db}}
}

type retrieve struct {
	base
	bunQ     *bun.SelectQuery
	scanArgs []interface{}
}

func newRetrieve(db *bun.DB) *retrieve {
	return &retrieve{bunQ: db.NewSelect(), base: base{db: db}}
}

type update struct {
	base
	bunQ *bun.UpdateQuery
}

func newUpdate(db *bun.DB) *update {
	return &update{bunQ: db.NewUpdate(), base: base{db: db}}
}

type del struct {
	base
	bunQ *bun.DeleteQuery
}

func newDelete(db *bun.DB) *del {
	return &del{bunQ: db.NewDelete(), base: base{db: db}}
}

// |||| EXEC ||||

func (c *create) Exec(ctx context.Context, p *query.Pack) error {
	c.convertOpts(p)
	c.exchange.ToDest()
	beforeInsertSetUUID(c.exchange.Dest)
	_, err := c.bunQ.Exec(ctx)
	c.exchange.ToSource()
	return newErrorHandler().Exec(err)
}

func (e *retrieve) Exec(ctx context.Context, p *query.Pack) error {
	e.convertOpts(p)
	err := e.bunQ.Scan(ctx, e.scanArgs...)
	e.exchange.ToSource()
	return newErrorHandler().Exec(err)
}

func (u *update) Exec(ctx context.Context, p *query.Pack) error {
	u.convertOpts(p)
	u.exchange.ToDest()
	_, err := u.bunQ.Exec(ctx)
	u.exchange.ToSource()
	return newErrorHandler().Exec(err)
}

func (d *del) Exec(ctx context.Context, p *query.Pack) error {
	d.convertOpts(p)
	_, err := d.bunQ.Exec(ctx)
	return newErrorHandler().Exec(err)
}

// |||| OPT CONVERTERS ||||

type OptConverter func(p *query.Pack)

type OptConverters []OptConverter

func (ocs OptConverters) Exec(p *query.Pack) {
	for _, oc := range ocs {
		oc(p)
	}
}

func (c *create) convertOpts(p *query.Pack) {
	OptConverters{c.model}.Exec(p)
}

func (u *update) convertOpts(p *query.Pack) {
	OptConverters{u.model, u.pk, u.fields, u.bulk}.Exec(p)
}

func (e *retrieve) convertOpts(p *query.Pack) {
	OptConverters{e.model, e.pk, e.fields, e.whereFields, e.relations, e.whereFields, e.calculate}.Exec(p)
}

func (d *del) convertOpts(p *query.Pack) {
	OptConverters{d.model, d.pk}.Exec(p)
}

// |||| MODEL ||||

func (b *base) exchangeToDest() {
	b.exchange.ToDest()
}

func (b *base) exchangeToSource() {
	b.exchange.ToSource()
}

func (b *base) model(p *query.Pack) interface{} {
	ptr := p.Model().Pointer()
	b.exchange = model.NewExchange(ptr, catalog().New(ptr))
	b.sqlGen = sqlGen{db: b.db, m: b.exchange.Dest}
	return b.exchange.Dest.Pointer()
}

func (c *create) model(p *query.Pack) {
	c.bunQ = c.bunQ.Model(c.base.model(p))
}

func (u *update) model(p *query.Pack) {
	u.bunQ = u.bunQ.Model(u.base.model(p))
}

func (e *retrieve) model(p *query.Pack) {
	e.bunQ = e.bunQ.Model(e.base.model(p))
}

func (d *del) model(p *query.Pack) {
	d.bunQ = d.bunQ.Model(d.base.model(p))
}

// |||| PK ||||

func (u *update) pk(p *query.Pack) {
	if pkc, ok := query.PKOpt(p); ok {
		u.bunQ = u.bunQ.Where(u.sqlGen.pks(), bun.In(pkc.Raw()))
	}
}

func (d *del) pk(p *query.Pack) {
	if pkc, ok := query.PKOpt(p); ok {
		d.bunQ = d.bunQ.Where(d.sqlGen.pks(), bun.In(pkc.Raw()))
	}
}

func (e *retrieve) pk(p *query.Pack) {
	if pkc, ok := query.PKOpt(p); ok {
		e.bunQ = e.bunQ.Where(e.sqlGen.pks(), bun.In(pkc.Raw()))
	}
}

// |||| FIELDS ||||

func (e *retrieve) fields(p *query.Pack) {
	if f, ok := query.RetrieveFieldsOpt(p); ok {
		e.bunQ = e.bunQ.Column(e.sqlGen.fieldNames(f...)...)
	}
}

func (u *update) fields(p *query.Pack) {
	if f, ok := query.RetrieveFieldsOpt(p); ok {
		u.bunQ = u.bunQ.Column(u.sqlGen.fieldNames(f...)...)
	}
}

// |||| CUSTOM RETRIEVE OPTS ||||

func (e *retrieve) whereFields(p *query.Pack) {
	if wf, ok := query.WhereFieldsOpt(p); ok {
		for fldN, fldV := range wf {
			relN, _ := model.SplitLastFieldName(fldN)
			if relN != "" {
				e.bunQ = e.bunQ.Relation(relN)
			}
			fldExp, args := e.sqlGen.relFldExp(fldN, fldV)
			e.bunQ = e.bunQ.Where(fldExp, args...)
		}
	}
}

func (e *retrieve) relations(p *query.Pack) {
	for _, opt := range query.RelationOpts(p) {
		e.bunQ = e.bunQ.Relation(opt.Rel, func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Column(e.sqlGen.fieldNames(opt.Fields...)...)
		})
	}
}

func (e *retrieve) calculate(p *query.Pack) {
	if c, ok := query.RetrieveCalcOpt(p); ok {
		e.scanArgs = append(e.scanArgs, c.Into)
		e.bunQ = e.bunQ.ColumnExpr(e.sqlGen.calc(c.Op), bun.Ident(e.sqlGen.fieldName(c.FldName)))
	}
}

// |||| CUSTOM UPDATE OPTS ||||

func (u *update) bulk(p *query.Pack) {
	if blk := query.BulkUpdateOpt(p); blk {
		u.bunQ = u.bunQ.Bulk()
	}
}