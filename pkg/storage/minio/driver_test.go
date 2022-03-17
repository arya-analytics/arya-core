package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("Driver", func() {
	BeforeEach(func() {
		viper.SetConfigFile("./testdata/config.json")
		Expect(viper.ReadInConfig()).To(BeNil())
	})
	Describe("Config", func() {
		It("Shoild load the viper config correctly", func() {
			cfg := minio.Config{}.Viper()
			Expect(cfg.Endpoint).To(Equal("badep:9000"))
		})
	})
	Describe("Connection Errors", func() {
		It("Should return the correct query error", func() {
			pool := storage.NewPool()
			driver := &minio.DriverMinio{Config: minio.Config{}.Viper()}
			engine := minio.New(driver, pool)
			err := engine.NewRetrieve().Model(&models.ChannelChunkReplica{}).WherePK(uuid.New()).Exec(ctx)
			Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeConnection))
		})
	})
	Describe("DemandCap", func() {
		It("Should return the correct demand cap", func() {
			driver := &minio.DriverMinio{Config: minio.Config{}.Viper()}
			Expect(driver.DemandCap()).To(Equal(500))
		})
	})
})
