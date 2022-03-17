package minio_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	pool   = internal.NewPool()
	engine = minio.New(mock.DriverMinio{}, pool)
	ctx    = context.Background()
)

var _ = BeforeSuite(func() {
	Expect(engine.NewMigrate().Exec(ctx)).To(BeNil())
})

func TestMinio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
