package dev

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/emoji"
	"github.com/arya-analytics/aryacore/pkg/util/git"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

var ignoreDirs = []string{
	".git",
	".idea",
}

const (
	imageRepo = "ghcr.io/arya-analytics/arya-core"
	chartRelPath = "kubernetes/aryacore"
)

// DefaultBuildCtxPath returns the default build context for the arya image.
func DefaultBuildCtxPath() string {
	ctx, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln(err)
	}
	return ctx
}

// StartReloader starts the development reloader.
func StartReloader(clusterName string, buildCtxPath string) error {
	log.Infof("%s Starting Reloader", emoji.Bolt)
	tag := GitImageTag()
	cfg := AryaClusterConfig{Name: clusterName}
	cluster := NewAryaCluster(cfg)
	cluster.Bind()
	chartPath := filepath.Join(buildCtxPath, chartRelPath)
	err := WatchAndDeployToLocalCluster(cluster, imageRepo, tag, chartPath,
		buildCtxPath)
	return err
}

// GitImageTag returns an image tag built off of the current commit hash and git
// username.
func GitImageTag() string {
	ch := git.CurrentCommitHash()
	u := git.Username()
	shortHash := ch[len(ch)-8:]
	shortUser := strings.Split(u, "@")[0]
	return fmt.Sprintf("%s-%s", shortHash, shortUser)
}


// WatchAndDeployToLocalCluster starts watching for file changes and continuously
// deploys those changes to a local development cluster.
func WatchAndDeployToLocalCluster(cluster *AryaCluster, repository, tag, chartPath, buildCtxPath string) error {
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
