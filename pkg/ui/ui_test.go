package ui_test

import (
	"github.com/arya-analytics/aryacore/pkg/ui"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UI", func() {
	Describe("Serving the application", func() {
		It("Should serve the application correctly", func() {
			app := fiber.New()
			server := ui.NewServer()
			Expect(func() {
				server.BindTo(app)
			}).ToNot(Panic())
		})
	})
})
