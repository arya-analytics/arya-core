package storage_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Tasks", func() {
	It("Should start and stop the task runner correctly", func() {
		store.StartTaskRunner(ctx)
		store.StopTaskRunner()
	})

})
