package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

const (
	dummyRepo = "mystrangerepo"
	dummyTag  = "mystrangetag"
)

var cfg = dev.ImageCfg{
	Repository:   dummyRepo,
	Tag:          dummyTag,
	BuildCtxPath: "../../",
}

var _ = Describe("Docker", func() {
	Describe("DockerImage", func() {
		Describe("Test NameTag", func() {
			It("Should generate the correct image nameTag", func() {
				di := dev.NewDockerImage(cfg)
				Expect(di.NameTag()).To(Equal(dummyRepo + ":" + dummyTag))
			})
		})
		Describe("Test Build", func() {
			It("Should build the docker image correctly", func() {
				di := dev.NewDockerImage(cfg)
				err := di.Build()
				Expect(err).To(BeNil())
			})
		})
		Describe("Test Push", func() {
			It("Should push the docker image correctly", func() {
				cfg := dev.ImageCfg{
					Repository: dev.DefaultImageRepo,
					Tag: dev.GitImageTag(),
					BuildCtxPath: "../../",
				}
				di := dev.NewDockerImage(cfg)
				if err := di.Build(); err != nil {
					log.Fatalln(err)
				}
				Expect(di.Push()).To(BeNil())
			})
		})
	})
})
