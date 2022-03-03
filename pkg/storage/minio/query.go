//nolint:typecheck
package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
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
			return storage.Error{
				Type:    storage.ErrorTypeInvalidArgs,
				Message: fmt.Sprintf("Minio data to write is nil! Model %s with id %s", c.exc.Dest().Type(), dv.PK),
			}
		}
		dv.Data.Reset()
		_, err := c.client.PutObject(ctx, c.bucket(), dv.PK.String(), dv.Data, dv.Data.Size(), minio.PutObjectOptions{})
		dv.Data.Reset()
		if err != nil {
			return newErrorHandler().Exec(err)
		}
	}
	c.exchangeToSource()
	return nil
}

func (d *del) exec(ctx context.Context, p *query.Pack) error {
	d.convertOpts(p)
	for _, pk := range d.pkc {
		if err := d.client.RemoveObject(ctx, d.bucket(), pk.String(), minio.RemoveObjectOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (r *retrieve) exec(ctx context.Context, p *query.Pack) error {
	r.convertOpts(p)
	if err := whereReqValidator().Exec(r.where).Error(); err != nil {
		return err
	}
	var dvc dataValueChain
	for _, pk := range r.pkc {
		resObj, gErr := r.client.GetObject(ctx, r.bucket(), pk.String(), minio.GetObjectOptions{})
		if gErr != nil {
			return newErrorHandler().Exec(gErr)
		}
		stat, sErr := resObj.Stat()
		if sErr != nil {
			return newErrorHandler().Exec(sErr)
		}
		bulk := telem.NewChunkData(make([]byte, stat.Size))
		if _, err := bulk.ReadFrom(resObj); err != nil {
			return newErrorHandler().Exec(err)
		}
		if err := resObj.Close(); err != nil {
			return newErrorHandler().Exec(err)
		}
		dvc = append(dvc, &dataValue{PK: pk, Data: bulk})
	}
	r.exc.bindDataVals(dvc)
	r.exchangeToSource()
	return nil
}

func (m *migrate) Exec(ctx context.Context) error {
	for _, mod := range catalog() {
		me := newWrappedExchange(model.NewExchange(mod, mod))
		bucketExists, err := m.client.BucketExists(ctx, me.bucket())
		if err != nil {
			return newErrorHandler().Exec(err)
		}
		if bucketExists {
			break
		}
		if mErr := m.client.MakeBucket(ctx, me.bucket(), minio.MakeBucketOptions{}); mErr != nil {
			return newErrorHandler().Exec(mErr)
		}
	}
	return nil
}

// |||| OPT CONVERTERS ||||

func (c *create) convertOpts(p *query.Pack) {
	storage.OptConverters{c.model}.Exec(p)
}

func (d *del) convertOpts(p *query.Pack) {
	storage.OptConverters{d.model, d.pk}.Exec(p)
}

func (r *retrieve) convertOpts(p *query.Pack) {
	storage.OptConverters{r.model, r.pk}.Exec(p)
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) {
	ptr := p.Model().Pointer()
	b.exc = newWrappedExchange(model.NewExchange(ptr, catalog().New(ptr)))
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

func (m *migrate) Verify(ctx context.Context) error {
	for _, mod := range catalog() {
		me := newWrappedExchange(model.NewExchange(mod, mod))
		exists, err := m.client.BucketExists(ctx, me.bucket())
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("bucket %s does not exist", err)
		}
	}
	return nil
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
