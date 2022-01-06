package dev

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/kubectl"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"os"
	"strings"
)

type DeploymentConfig struct {
	name      string
	chartPath string
	cluster   *AryaCluster
	imageCfg  ImageCfg
}

type Deployment struct {
	cfg          DeploymentConfig
	actionConfig *action.Configuration
	settings     *cli.EnvSettings
}

const (
	helmDriverEnv          = "HELM_DRIVER"
	imageRepoKey           = "image.repository"
	imageTagKey            = "image.tag"
	aryaCoreDeploymentName = "aryacore-deployment"
)

// NewDeployment creates a new deployment based on the specified config
func NewDeployment(cfg DeploymentConfig) (*Deployment, error) {
	d := Deployment{
		cfg:          cfg,
		actionConfig: new(action.Configuration),
		settings:     cli.New(),
	}
	if err := d.initActionConfig(); err != nil {
		log.Fatalln(err)
	}
	return &d, nil
}

func (d Deployment) initActionConfig() error {
	if err := d.actionConfig.Init(d.settings.RESTClientGetter(), d.settings.Namespace(),
		os.Getenv(helmDriverEnv), log.Printf); err != nil {
		return err
	}
	return nil
}

// RedeployArya redeploys arya into the cluster
func (d Deployment) RedeployArya() error {
	name := d.cfg.name + "-" + aryaCoreDeploymentName
	var err error
	fmt.Println(name)
	d.iterNodes(func(node *K3sCluster) {
		err = kubectl.Exec("rollout", "restart", "deployment", name)
	})
	return err
}

// Install installs the deployment into the cluster
func (d Deployment) Install() error {
	client := action.NewInstall(d.actionConfig)
	client.ReleaseName = d.cfg.name
	repo := fmt.Sprintf("%s=%s", imageRepoKey, d.cfg.imageCfg.Repository)
	tag := fmt.Sprintf("%s=%s", imageTagKey, d.cfg.imageCfg.Tag)
	c, err := d.chart()
	if err != nil {
		log.Fatalln(err)
	}

	var nodeIPs []string
	for _, node := range d.cfg.cluster.Nodes() {
		info, err := node.VM.Info()
		if err != nil {
			log.Fatalln(err)
		}
		nodeIPs = append(nodeIPs, info.IPv4)
	}

	d.iterNodes(func(node *K3sCluster) {
		info, err := node.VM.Info()
		if err != nil {
			log.Fatalln(err)
		}
		nodeIP := info.IPv4
		nodeIPVal := fmt.Sprintf("%s=%s", "cockroachdb.nodeIP", nodeIP)
		var nodeIpList []string
		for _, v := range nodeIPs {
			if v != nodeIP {
				nodeIpList = append(nodeIpList, v)
			}
		}
		nodeIPsString := strings.Join(nodeIpList,"\\,")
		joinVal := fmt.Sprintf("%s=%s", "cockroachdb.join", nodeIPsString)
		imageVals := []string{repo, tag, nodeIPVal, joinVal }
		options := values.Options{Values: imageVals}
		v, err := options.MergeValues(getter.All(d.settings))
		if err != nil {
			log.Fatalln(err)
		}

		_, err = client.Run(c, v)
	})
	return err
}

func (d Deployment) chart() (*chart.Chart, error) {
	return loader.LoadDir(d.cfg.chartPath)
}

// Uninstall uninstalls the deployment from the cluster
func (d Deployment) Uninstall() error {
	d.iterNodes(func(node *K3sCluster) {
		client := action.NewUninstall(d.actionConfig)
		_, _ = client.Run(d.cfg.name)
	})
	return nil
}

func (d Deployment) iterNodes(exec func(node *K3sCluster)) {
	for _, node := range d.cfg.cluster.Nodes() {
		if err := kubectl.SwitchContext(node.VM.Name()); err != nil {
			log.Fatalln(err)
		}
		if err := d.initActionConfig(); err != nil {
			log.Fatalln(err)
		}
		exec(node)
	}
}


