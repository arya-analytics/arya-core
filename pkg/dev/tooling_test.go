package dev_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tooling", func() {
	Describe("Brew tooling", func() {
		Context("A tool is not installed", func() {
			Describe("Checking if a tool is installed", func() {
				It("Should return false", func() {
					Expect(tooling.Installed("myverystrangerandomtool")).To(BeFalse())
				})
			})
			Describe("Uninstalling a tool", func() {
				It("Should throw an error", func() {
					err := tooling.Uninstall("ra012hjad");
					Expect(err).ToNot(BeNil())
				})
			})
			Describe("Installing a tool", func() {
				It("Should install the tool", func() {
					if err := tooling.Uninstall(testTool); err != nil {
						Fail("Failed to uninstall test tool")
					}
					err := tooling.Install(testTool)
					Expect(err).To(BeNil())
					Expect(tooling.Installed(testTool)).To(BeTrue())
				})
			})
		})
		Context("A tool is installed", func() {
			Describe("Checking if a tool is installed", func() {
				It("Should return true if a tool is installed", func() {
					Expect(tooling.Installed(testTool)).To(BeTrue())
				})
			})
			Describe("Uninstalling a tool", func() {
				It("(SLOW) Should uninstall the tool", func() {
					err := tooling.Uninstall(testTool)
					Expect(err).To(BeNil())
					Expect(tooling.Installed(testTool)).To(BeFalse())
					if err := tooling.Install(testTool); err != nil {
						Fail("Failed to install test tool")
					}
				})
			})
			Describe("Installing a tool", func() {
				It("Shouldn't throw an error", func() {
					err := tooling.Install(testTool)
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
