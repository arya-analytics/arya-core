package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Driver", func() {
	Describe("Connection Errors", func() {
		It("Should return the correct query error", func() {
			pool := storage.NewPool()
			driver := &minio.DriverMinio{minio.Config{
				Endpoint:  "badhost:1234",
				AccessKey: "badAccessKey",
				SecretKey: "badSecretKey",
				UseTLS:    false,
			}}
			engine := minio.New(driver, pool)
			err := engine.NewRetrieve().Model(&models.ChannelChunkReplica{}).WherePK(uuid.New()).Exec(ctx)
			Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeConnection))
		})
	})
})
