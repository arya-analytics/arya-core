package minio_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/util/pool"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	p      = pool.New[internal.Engine]()
	engine = minio.New(mock.DriverMinio{}, p)
	ctx    = context.Background()
)

var _ = BeforeSuite(func() {
	p.AddFactory(engine)
	Expect(engine.NewMigrate().Exec(ctx)).To(BeNil())
})

func TestMinio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
