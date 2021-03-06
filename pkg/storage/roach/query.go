package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/uptrace/bun"
	bunMigrate "github.com/uptrace/bun/migrate"
)

type base struct {
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
	exc      *model.Exchange
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
	beforeInsertSetUUID(p.Model())
	_, err := c.bunQ.Exec(ctx)
	return err
}

func (e *retrieve) exec(ctx context.Context, p *query.Pack) error {
	e.convertOpts(p)
	err := e.bunQ.Scan(ctx, e.scanArgs...)
	if e.exc != nil {
		e.exc.ToSource()
	}
	return err
}

func (u *update) exec(ctx context.Context, p *query.Pack) error {
	u.convertOpts(p)
	_, err := u.bunQ.Exec(ctx)
	return err
}

func (d *del) exec(ctx context.Context, p *query.Pack) error {
	d.convertOpts(p)
	_, err := d.bunQ.Exec(ctx)
	return err
}

func (m *migrate) exec(ctx context.Context, p *query.Pack) error {
	c := errutil.NewCatchContext(ctx)
	if m.verify(p) {
		_, err := m.db.NewSelect().Model((*models.ChannelConfig)(nil)).Count(ctx)
		return err
	}
	bindMigrations(m.bunQ)
	bunMig := bunMigrate.NewMigrator(m.db, m.bunQ)
	c.Exec(bunMig.Init)
	c.Exec(func(ctx context.Context) error {
		_, err := bunMig.Migrate(ctx)
		return err
	})
	return c.Error()
}

// |||| OPT CONVERTERS ||||

func (c *create) convertOpts(p *query.Pack) {
	query.OptConvertChain{c.model}.Exec(p)
}

func (u *update) convertOpts(p *query.Pack) {
	query.OptConvertChain{u.model, u.pk, u.fields, u.bulk}.Exec(p)
}

func (e *retrieve) convertOpts(p *query.Pack) {
	query.OptConvertChain{
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
	query.OptConvertChain{d.model, d.pk}.Exec(p)
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) interface{} {
	b.sql = sqlGen{db: b.db, m: p.Model()}
	return p.Model().Pointer()
}

func (c *create) model(p *query.Pack) {
	c.bunQ = c.bunQ.Model(c.base.model(p))
}

func (u *update) model(p *query.Pack) {
	u.bunQ = u.bunQ.Model(u.base.model(p))
}

func (e *retrieve) model(p *query.Pack) {
	if p.Model().IsChain() && p.Model().ChainValue().Len() > 0 {
		e.exc = model.NewExchange(e.base.model(p), p.Model().NewRaw())
		e.bunQ = e.bunQ.Model(e.exc.Dest().Pointer())
	} else {
		e.bunQ = e.bunQ.Model(e.base.model(p))
	}
}

func (d *del) model(p *query.Pack) {
	d.bunQ = d.bunQ.Model(d.base.model(p))
}

// |||| PK ||||

func (u *update) pk(p *query.Pack) {
	if pkc, ok := query.RetrievePKOpt(p); ok {
		u.bunQ = u.bunQ.Where(u.sql.pks(), bun.In(pkc.Raw()))
	}
}

func (d *del) pk(p *query.Pack) {
	if pkc, ok := query.RetrievePKOpt(p); ok {
		d.bunQ = d.bunQ.Where(d.sql.pks(), bun.In(pkc.Raw()))
	}
}

func (e *retrieve) pk(p *query.Pack) {
	if pkc, ok := query.RetrievePKOpt(p); ok {
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
	if wf, ok := query.RetrieveWhereFieldsOpt(p); ok {
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
	for _, opt := range query.RetrieveRelationOpts(p) {
		// CLARIFICATION: Still don't know exactly why it needs to be called this way, but it does for the
		// correct opt to be provided.
		func(opt query.RelationOpt) {
			e.bunQ = e.bunQ.Relation(opt.Name, func(sq *bun.SelectQuery) *bun.SelectQuery {
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
	if limit, ok := query.RetrieveLimitOpt(p); ok {
		e.bunQ = e.bunQ.Limit(limit)
	}
}

func (e *retrieve) order(p *query.Pack) {
	if o, ok := query.RetrieveOrderOpt(p); ok {
		e.bunQ = e.bunQ.Order(e.sql.order(o.Direction, o.Field))
	}
}

// |||| CUSTOM UPDATE OPTS ||||

func (u *update) bulk(p *query.Pack) {
	if blk := query.RetrieveBulkUpdateOpt(p); blk {
		u.bunQ = u.bunQ.Bulk()
	}
}

// |||| CUSTOM MIGRATE OPTS ||||

func (m *migrate) verify(p *query.Pack) bool {
	return query.RetrieveVerifyOpt(p)
}
