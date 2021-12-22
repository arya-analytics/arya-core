package dev_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)




var _ = Describe("Tooling", func() {
	Describe("Brew tooling", func() {
		Describe("Checking if a tool is installed", func() {
			It("Should return false if a tool is not installed", func() {
				Expect(tooling.Installed("myverystrangerandomtool")).To(BeFalse())
			})
			It("Should return true if a tool is installed", func() {
				Expect(tooling.Installed(testTool)).To(BeTrue())
			})
		})

	})
})
