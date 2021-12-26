package dev

import (
	"github.com/arya-analytics/aryacore/pkg/util/emoji"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

var ignoreDirs = []string{
	".git",
	".idea",
}

func WatchAndDeploy(cluster *AryaCluster, repository, tag, chartPath, buildCtxPath string) error {
	imgCfg := ImageCfg{
		Repository:   repository,
		Tag:          tag,
		BuildCtxPath: buildCtxPath,
	}

	img := NewDockerImage(imgCfg)
	dCfg := DeploymentConfig{
		name:      "aryacore",
		chartPath: chartPath,
		cluster:   cluster,
		imageCfg:  imgCfg,
	}

	d, err := NewDeployment(dCfg)
	if err != nil {
		return err
	}

	log.Infof("%s Compiling and deploying Arya", emoji.Rainbow)

	if err := d.Uninstall(); err != nil {
		log.Fatalln(err)
	}

	if err := d.Install(); err != nil {
		log.Fatalln(err)
	}

	log.Infof("%s Successfully deployed", emoji.Check)

	wCfg := WatcherConfig{
		Dirs:       []string{buildCtxPath},
		Recursive:  true,
		IgnoreDirs: ignoreDirs,
		Triggers:   []fsnotify.Op{fsnotify.Write},
		Action: func(event fsnotify.Event) {
			log.Infof("%s Building Image", emoji.Tools)
			if err := img.Build(); err != nil {
				log.Fatalln(err)
			}

			log.Infof("%s Pushing Image", emoji.Flame)
			if err := img.Push(); err != nil {
				log.Fatalln(err)
			}
			log.Infof("%s Re-deploying Image", emoji.Bison)
			if err := d.RedeployArya(); err != nil {
				log.Fatalln(err)
			}

		},
	}

	w := Watcher{cfg: wCfg}
	log.Infof("%s Listening for changes", emoji.Sparks)
	w.Start()
	return nil
}
