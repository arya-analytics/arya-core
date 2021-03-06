package minio

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"github.com/minio/minio-go/v7"
)

type base struct {
	client       *minio.Client
	wrappedModel *reflectMinio
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
	data []data
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

// |||| EXEC |||||

func (c *create) exec(ctx context.Context, p *query.Pack) error {
	c.convertOpts(p)
	for _, dv := range c.wrappedModel.dataValues() {
		if dv.Data == nil {
			return query.Error{
				Type:    query.ErrorTypeInvalidArgs,
				Message: fmt.Sprintf("Minio data to write is nil! Model %s with id %s", c.wrappedModel.Type(), dv.PK),
			}
		}
		dv.Data.Reset()
		_, err := c.client.PutObject(ctx, c.wrappedModel.bucket(), dv.PK.String(), dv.Data, dv.Data.Size(), minio.PutObjectOptions{})
		dv.Data.Reset()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *del) exec(ctx context.Context, p *query.Pack) error {
	d.convertOpts(p)
	for _, pk := range d.pkc {
		if err := d.client.RemoveObject(ctx, d.wrappedModel.bucket(), pk.String(), minio.RemoveObjectOptions{}); err != nil {
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
	var d []data
	for _, pk := range r.pkc {
		bulk, err := r.getObject(ctx, pk)
		if err != nil {
			return err
		}
		d = append(d, data{PK: pk, Data: bulk})
	}
	r.wrappedModel.bindDataVals(d)

	return nil
}

func migrationCatalog() model.Catalog {
	return model.Catalog{
		&models.ChannelChunkReplica{},
	}
}

func (m *migrate) exec(ctx context.Context, p *query.Pack) error {
	for _, mod := range migrationCatalog() {
		me := wrapReflect(model.NewReflect(mod))
		exists, err := m.client.BucketExists(ctx, me.bucket())
		if err != nil {
			return err
		}
		if !exists {
			if m.verify(p) {
				return fmt.Errorf("bucket %s does not exist", err)
			}
			if mErr := m.client.MakeBucket(ctx, me.bucket(), minio.MakeBucketOptions{}); mErr != nil {
				return mErr
			}
		}

	}
	return nil
}

// |||| OPT CONVERTERS ||||

func (c *create) convertOpts(p *query.Pack) {
	query.OptConvertChain{c.model}.Exec(p)
}

func (d *del) convertOpts(p *query.Pack) {
	query.OptConvertChain{d.model, d.pk}.Exec(p)
}

func (r *retrieve) convertOpts(p *query.Pack) {
	query.OptConvertChain{r.model, r.pk}.Exec(p)
}

// |||| MODEL ||||

func (b *base) model(p *query.Pack) {
	b.wrappedModel = wrapReflect(p.Model())
}

// |||| PK ||||

func (w *where) pk(p *query.Pack) {
	if pkc, ok := query.RetrievePKOpt(p); ok {
		w.pkc = pkc
	} else if p.Model().PKChain().AllNonZero() {
		// CLARIFICATION: If there wasn't a primary key specified, try to pull the primary key
		// from the model.
		w.pkc = p.Model().PKChain()
	} else {
		panic("where queries require a primary key! tried to pull from model, but was unable to")
	}
}

// |||| CUSTOM RETRIEVE ||||

func (r *retrieve) getObject(ctx context.Context, pk model.PK) (*telem.ChunkData, error) {
	var (
		c      = errutil.NewCatchSimple()
		resObj *minio.Object
		stat   minio.ObjectInfo
		bulk   *telem.ChunkData
	)
	c.Exec(func() (err error) {
		resObj, err = r.client.GetObject(ctx, r.wrappedModel.bucket(), pk.String(), minio.GetObjectOptions{})
		return err
	})
	c.Exec(func() (err error) {
		stat, err = resObj.Stat()
		return err
	})
	c.Exec(func() error {
		bulk = telem.NewChunkData(make([]byte, stat.Size))
		_, err := bulk.ReadFrom(resObj)
		return err
	})
	return bulk, c.Error()
}

// |||| CUSTOM MIGRATE ||||

func (m *migrate) verify(p *query.Pack) bool {
	return query.RetrieveVerifyOpt(p)
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
