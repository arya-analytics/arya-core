package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
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
		It("Should load the viper config correctly", func() {
			cfg := redis.Config{}.Viper()
			Expect(cfg.Host).To(Equal("redistsfake"))
		})
	})
	Describe("Connection errors", func() {
		It("Should return the correct query error", func() {
			pool := internal.NewPool()
			cfg := redis.Config{}.Viper()
			driver := &redis.DriverRedis{Config: cfg}
			engine := redis.New(driver, pool)
			err := engine.NewTSRetrieve().Model(&models.ChannelSample{}).WherePK(uuid.New()).Exec(ctx)
			Expect(err).ToNot(BeNil())
			Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeConnection))
		})
	})
	Describe("DemandCap", func() {
		It("Should return the correct demand cap", func() {
			cfg := redis.Config{}.Viper()
			driver := &redis.DriverRedis{Config: cfg}
			Expect(driver.DemandCap()).To(Equal(500))
		})
	})

})
