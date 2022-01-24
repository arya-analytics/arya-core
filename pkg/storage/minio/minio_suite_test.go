package minio_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	mockEngine = minio.New(minio.Config{
		Driver:    minio.DriverMinIO,
		Endpoint:  "play.min.io",
		AccessKey: "Q3AM3UQ867SPQQA43P2F",
		SecretKey: "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
	})
	mockAdapter = mockEngine.NewAdapter()
	mockCtx     = context.Background()
)

func TestMinio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Minio Suite")
}
