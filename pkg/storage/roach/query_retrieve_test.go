package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type BunSQLError interface {
	Code() int
}

var _ = FDescribe("QueryRetrieve", func() {
	BeforeEach(migrate)
	Describe("Edge cases + errors", func() {
		Context("Retrieving an item that doesn't exist", func() {
			It("Should return the correct error type", func() {
				somePKThatDoesntExist := 136987
				m := &storage.ChannelConfig{}
				err := dummyEngine.NewRetrieve(dummyAdapter).
					Model(m).
					WherePK(somePKThatDoesntExist).
					Exec(dummyCtx)

				log.Info(reflect.ValueOf(err).Elem().Type())
				//berr := err.(BunSQLError)
				//log.Info(berr.Code())
				//
				//t := reflect.ValueOf(err).Type()
				//for i := 0; i < t.NumMethod(); i++ {
				//	log.Info(t.Method(i))
				//}
				//log.Info(reflect.ValueOf(err).Type().Num)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
