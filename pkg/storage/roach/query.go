package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type base struct {
	exc *model.Exchange
	sql sqlGen
	db  *bun.DB
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

type migrate struct {
	base
	bunQ *bunMigrate.Migrations
}

func newMigrate(db *bun.DB) *migrate {
	return &migrate{bunQ: bunMigrate.NewMigrations(), base: base{db: db}}
}

// |||| EXEC ||||

func (c *create) exec(ctx context.Context, p *query.Pack) error {
	c.convertOpts(p)
	c.exc.ToDest()
	beforeInsertSetUUID(c.exc.Dest())
	_, err := c.bunQ.Exec(ctx)
	c.exc.ToSource()
	return newErrorConvert().Exec(err)
}

func (e *retrieve) exec(ctx context.Context, p *query.Pack) error {
	e.convertOpts(p)
	err := e.bunQ.Scan(ctx, e.scanArgs...)
	e.exc.ToSource()
	return newErrorConvert().Exec(err)
}

func (u *update) exec(ctx context.Context, p *query.Pack) error {
	u.convertOpts(p)
	u.exc.ToDest()
	_, err := u.bunQ.Exec(ctx)
	u.exc.ToSource()
	return newErrorConvert().Exec(err)
}

func (d *del) exec(ctx context.Context, p *query.Pack) error {
	d.convertOpts(p)
	_, err := d.bunQ.Exec(ctx)
	return newErrorConvert().Exec(err)
}

func (m *migrate) exec(ctx context.Context, p *query.Pack) error {
	c := errutil.NewCatchContext(ctx)
	if m.verify(p) {
		_, err := m.db.NewSelect().Model((*ChannelConfig)(nil)).Count(ctx)
		return newErrorConvert().Exec(err)
	}
	bindMigrations(m.bunQ)
	bunMig := bunMigrate.NewMigrator(m.db, m.bunQ)
	c.Exec(bunMig.Init)
	c.Exec(func(ctx context.Context) error {
		_, err := bunMig.Migrate(ctx)
		return err
	})
	return newErrorConvert().Exec(c.Error())
}

// |||| OPT CONVERTERS ||||

func (c *create) convertOpts(p *query.Pack) {
	internal.OptConverters{c.model}.Exec(p)
}

func (u *update) convertOpts(p *query.Pack) {
	internal.OptConverters{u.model, u.pk, u.fields, u.bulk}.Exec(p)
}

func (e *retrieve) convertOpts(p *query.Pack) {
	internal.OptConverters{
		e.model,
		e.pk,
		e.fields,
		e.whereFields,
		e.relations,
		e.whereFields,
		e.calculate,
		e.limit,
		e.order,
	}.Exec(p)
}

func (d *del) convertOpts(p *query.Pack) {
	internal.OptConverters{d.model, d.pk}.Exec(p)
}

// |||| BASE ||||

func (b *base) exchangeToDest() {
	b.exc.ToDest()
}

func (b *base) exchangeToSource() {
	b.exc.ToSource()
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) interface{} {
	ptr := p.Model().Pointer()
	b.exc = model.NewExchange(ptr, catalog().New(ptr))
	b.sql = sqlGen{db: b.db, m: b.exc.Dest()}
	return b.exc.Dest().Pointer()
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
		u.bunQ = u.bunQ.Where(u.sql.pks(), bun.In(pkc.Raw()))
	}
}

func (d *del) pk(p *query.Pack) {
	if pkc, ok := query.PKOpt(p); ok {
		d.bunQ = d.bunQ.Where(d.sql.pks(), bun.In(pkc.Raw()))
	}
}

func (e *retrieve) pk(p *query.Pack) {
	if pkc, ok := query.PKOpt(p); ok {
		e.bunQ = e.bunQ.Where(e.sql.pks(), bun.In(pkc.Raw()))
	}
}

// |||| FIELDS ||||

func (e *retrieve) fields(p *query.Pack) {
	if f, ok := query.RetrieveFieldsOpt(p); ok {
		e.bunQ = e.bunQ.Column(e.sql.fieldNames(f...)...)
	}
}

func (u *update) fields(p *query.Pack) {
	if f, ok := query.RetrieveFieldsOpt(p); ok {
		u.bunQ = u.bunQ.Column(u.sql.fieldNames(f...)...)
	}
}

// |||| WHERE FIELDS

func (e *retrieve) whereFields(p *query.Pack) {
	if wf, ok := query.WhereFieldsOpt(p); ok {
		for fldN, fldV := range wf {
			relN, _ := model.SplitLastFieldName(fldN)
			if relN != "" {
				e.bunQ = e.bunQ.Relation(relN)
			}
			fldExp, args := e.sql.relFldExp(fldN, fldV)
			e.bunQ = e.bunQ.Where(fldExp, args...)
		}
	}
}

// |||| CUSTOM RETRIEVE OPTS ||||

func (e *retrieve) relations(p *query.Pack) {
	for _, opt := range query.RelationOpts(p) {
		// CLARIFICATION: Still don't know exactly why it needs to be called this way, but it does for the
		// correct opt to be provided.
		func(opt query.RelationOpt) {
			e.bunQ = e.bunQ.Relation(opt.Rel, func(sq *bun.SelectQuery) *bun.SelectQuery {
				return sq.Column(e.sql.fieldNames(opt.Fields...)...)
			})
		}(opt)
	}
}

func (e *retrieve) calculate(p *query.Pack) {
	if c, ok := query.RetrieveCalcOpt(p); ok {
		e.scanArgs = append(e.scanArgs, c.Into)
		e.bunQ = e.bunQ.ColumnExpr(e.sql.calc(c.Op), bun.Ident(e.sql.fieldName(c.Field)))
	}
}

func (e *retrieve) limit(p *query.Pack) {
	if limit, ok := query.LimitOpt(p); ok {
		e.bunQ = e.bunQ.Limit(limit)
	}
}

func (e *retrieve) order(p *query.Pack) {
	if o, ok := query.RetrieveOrderOpt(p); ok {
		e.bunQ = e.bunQ.Order(e.sql.order(o.Order, o.Field))
	}
}

// |||| CUSTOM UPDATE OPTS ||||

func (u *update) bulk(p *query.Pack) {
	if blk := query.BulkUpdateOpt(p); blk {
		u.bunQ = u.bunQ.Bulk()
	}
}

// |||| CUSTOM MIGRATE OPTS ||||

func (m *migrate) verify(p *query.Pack) bool {
	return query.VerifyOpt(p)
}
