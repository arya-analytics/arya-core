package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

type base struct {
	client *minio.Client
	exc    *exchange
}

type create struct {
	base
}

func newCreate(client *minio.Client) *create {
	return &create{base: base{client: client}}
}

type where struct {
	base
	pkc model.PKChain
}

type del struct {
	where
}

func newDelete(client *minio.Client) *del {
	return &del{where{base: base{client: client}}}
}

type retrieve struct {
	where
	dvc dataValueChain
}

func newRetrieve(client *minio.Client) *retrieve {
	return &retrieve{where: where{base: base{client: client}}}
}

type migrate struct {
	base
}

func newMigrate(client *minio.Client) *migrate {
	return &migrate{base: base{client: client}}
}

// |||| BASE ||||

func (b *base) bucket() string {
	return b.exc.bucket()
}

func (b *base) exchangeToSource() {
	b.exc.ToSource()
}

func (b *base) exchangeToDest() {
	b.exc.ToDest()
}

// |||| EXEC |||||

func (c *create) exec(ctx context.Context, p *query.Pack) error {
	c.convertOpts(p)
	c.exchangeToDest()
	for _, dv := range c.exc.dataVals() {
		if dv.Data == nil {
			return query.Error{
				Type:    query.ErrorTypeInvalidArgs,
				Message: fmt.Sprintf("Minio data to write is nil! Model %s with id %s", c.exc.Dest().Type(), dv.PK),
			}
		}
		dv.Data.Reset()
		_, err := c.client.PutObject(ctx, c.bucket(), dv.PK.String(), dv.Data, dv.Data.Size(), minio.PutObjectOptions{})
		dv.Data.Reset()
		if err != nil {
			return newErrorConvert().Exec(err)
		}
	}
	c.exchangeToSource()
	return nil
}

func (d *del) exec(ctx context.Context, p *query.Pack) error {
	d.convertOpts(p)
	c := errutil.NewCatchSimple(errutil.WithConvert(newErrorConvert()))
	for _, pk := range d.pkc {
		c.AddError(d.client.RemoveObject(ctx, d.bucket(), pk.String(), minio.RemoveObjectOptions{}))
	}
	return c.Error()
}

func (r *retrieve) exec(ctx context.Context, p *query.Pack) error {
	r.convertOpts(p)
	c := errutil.NewCatchSimple(errutil.WithConvert(newErrorConvert()))
	c.AddError(whereReqValidator().Exec(r.where).Error())
	var dvc dataValueChain
	for _, pk := range r.pkc {
		resObj, gErr := r.client.GetObject(ctx, r.bucket(), pk.String(), minio.GetObjectOptions{})
		c.AddError(gErr)
		stat, sErr := resObj.Stat()
		c.AddError(sErr)
		bulk := telem.NewChunkData(make([]byte, stat.Size))
		_, err := bulk.ReadFrom(resObj)
		c.AddError(err)
		c.Exec(resObj.Close)
		dvc = append(dvc, &dataValue{PK: pk, Data: bulk})
	}
	if err := c.Error(); err != nil {
		return err
	}
	r.exc.bindDataVals(dvc)
	r.exchangeToSource()
	return nil
}

func (m *migrate) exec(ctx context.Context, p *query.Pack) error {
	c := errutil.NewCatchSimple(errutil.WithConvert(newErrorConvert()))
	for _, mod := range catalog() {
		me := wrapExchange(model.NewExchange(mod, mod))
		exists, err := m.client.BucketExists(ctx, me.bucket())
		c.AddError(err)
		if !exists {
			if m.verify(p) {
				c.AddError(fmt.Errorf("bucket %s does not exist", err))
			}
			c.AddError(m.client.MakeBucket(ctx, me.bucket(), minio.MakeBucketOptions{}))
		}
	}
	return c.Error()
}

// |||| OPT CONVERTERS ||||

func (c *create) convertOpts(p *query.Pack) {
	internal.OptConverters{c.model}.Exec(p)
}

func (d *del) convertOpts(p *query.Pack) {
	internal.OptConverters{d.model, d.pk}.Exec(p)
}

func (r *retrieve) convertOpts(p *query.Pack) {
	internal.OptConverters{r.model, r.pk}.Exec(p)
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) {
	ptr := p.Model().Pointer()
	b.exc = wrapExchange(model.NewExchange(ptr, catalog().New(ptr)))
}

// |||| PK ||||

func (w *where) pk(p *query.Pack) {
	if pkc, ok := query.PKOpt(p); ok {
		w.pkc = pkc
	} else if p.Model().PKChain().AllNonZero() {
		// CLARIFICATION: If there wasn't a primary key specified, try to pull the primary key
		// from the model.
		w.pkc = p.Model().PKChain()
	} else {
		panic("where queries require a primary key! tried to pull from model, but was unable to")
	}
}

// |||| CUSTOM MIGRATE ||||

func (m *migrate) verify(p *query.Pack) bool {
	return query.VerifyOpt(p)
}

// |||| VALIDATORS ||||

// || WHERE ||

func whereReqValidator() *validate.Validate[where] {
	return validate.New([]func(w where) error{
		validatePKProvided,
	})
}

func validatePKProvided(w where) error {
	if w.pkc.AllZero() {
		panic("where queries require a primary key!")
	}
	return nil
}
