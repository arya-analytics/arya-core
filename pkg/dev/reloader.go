package dev

import (
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

	if err := d.Uninstall(); err != nil {
		log.Fatalln(err)
	}

	if err := d.Install(); err != nil {
		log.Fatalln(err)
	}

	wCfg := WatcherConfig{
		Dirs:       []string{buildCtxPath},
		Recursive:  true,
		IgnoreDirs: ignoreDirs,
		Triggers:   []fsnotify.Op{fsnotify.Write},
		Action: func(event fsnotify.Event) {
			if err := img.Build(); err != nil {
				log.Fatalln(err)
			}
			if err := img.Push(); err != nil {
				log.Fatalln(err)
			}
			if err := d.RedeployArya(); err != nil {
				log.Fatalln(err)
			}

		},
	}

	w := Watcher{cfg: wCfg}
	w.Start()
	return nil
}
