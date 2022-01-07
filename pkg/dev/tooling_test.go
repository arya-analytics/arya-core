package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("Tooling", func() {
	BeforeEach(func() {
		dev.RequiredTools = dev.Tools{testTool}
	})
	Describe("Brew tools", func() {
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
	Describe("InstallRequiredTools()",func() {
		It("Should install the required tools correctly", func() {
			err := dev.InstallRequiredTools()
			Expect(err).To(BeNil())
		})
	})
	Describe("UninstallRequiredTools", func() {
		It("Should uninstall the required tools correctly", func() {
			err := dev.UninstallRequiredTools()
			Expect(err).To(BeNil())
		})
	})
	Describe("RequiredToolsInstalled",func() {
		It("Should return true if the required tools are installed", func() {
			if err := dev.InstallRequiredTools(); err != nil {
				log.Fatalln(err)
			}
			i := dev.RequiredToolsInstalled()
			Expect(i).To(BeTrue())
		})
		It("Should return false if the required tools aren't installed", func() {
			if err := dev.UninstallRequiredTools(); err != nil {
				log.Fatalln(err)
			}
			i := dev.RequiredToolsInstalled()
			Expect(i).To(BeFalse())
		})

	})
})
