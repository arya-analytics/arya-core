package minio_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	engine  = minio.New(mock.DriverMinio{})
	adapter = engine.NewAdapter()
	ctx     = context.Background()
)

func TestMinio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
