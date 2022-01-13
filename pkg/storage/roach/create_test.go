package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
)

var c = &storage.ChannelConfig{
	ID:   432,
	Name: "Cool Name",
}

var _ = Describe("Create", func() {
	Describe("Create a new Channel Config", func() {
		It("Should create it without error", func() {
			ctx := context.Background()
			a := dummyEngine.NewAdapter()
			if err := dummyEngine.Migrate(ctx, a); err != nil {
				log.Fatalln(err)
			}
			err := dummyEngine.NewCreate(a).Model(c).Exec(ctx)
			Expect(err).To(BeNil())
		})
	})

})
