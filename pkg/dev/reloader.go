package dev

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/git"
	"github.com/fsnotify/fsnotify"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var ignoreDirs = []string{
	".git",
	".idea",
}

func WatchAndReload() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)

	hash := git.CurrentCommitHash()
	username := git.Username()

	_, filename, _, _ := runtime.Caller(1)
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))

	bc := context.Background()

	nameTag := fmt.Sprintf("ghcr.io/arya-analytics/arya-core:%s-%s",
		hash[len(hash)-8:],
		strings.Split(username, "@")[0])

	fmt.Println(bc, nameTag)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("Modified file %s \n", event.Name)
					BuildDockerImage(bc, nameTag)
					PushDockerImage(bc, nameTag)
					RestartDeployment()
					//DeployHelmChart(nameTag)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)

			}
		}
	}()
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	fileSystem := os.DirFS(dir)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		ignoreDir := false
		for _, v := range ignoreDirs {
			if strings.Contains(path, v) {
				ignoreDir = true
				break
			}
		}
		if d.IsDir() && !ignoreDir {
			err = watcher.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	<-done
}

func BuildDockerImage(pCtx context.Context, nameTag string) context.CancelFunc {
	_, filename, _, _ := runtime.Caller(1)
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	ctx, cancel := context.WithCancel(pCtx)
	cmdString := fmt.Sprintf("docker build %s -t %s", dir, nameTag)
	fmt.Println(cmdString)
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdString)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return cancel
}

func PushDockerImage(pCtx context.Context, nameTag string) context.CancelFunc {
	cmdString := fmt.Sprintf("docker push %s", nameTag)
	ctx, cancel := context.WithCancel(pCtx)
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdString)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return cancel
}

func DeployHelmChart(imageNameTag string) {
	fmt.Println("Deploying helm chart")
	settings := cli.New()

	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(),
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		panic(err)
	}

	client := action.NewUpgrade(actionConfig)
	//client.ReleaseName = "arya-core"

	splitImageNameTag := strings.Split(imageNameTag, ":")

	vals := values.Options{Values: []string{fmt.Sprintf("image.repository=%s",
		splitImageNameTag[0]),
		fmt.Sprintf("image.tag=%s", splitImageNameTag[1])}}
	p := getter.All(settings)
	v, err := vals.MergeValues(p)

	fmt.Println(v)

	chart, err := loader.LoadDir("./kubernetes/aryacore")
	if err != nil {
		panic(err)
	}

	results, err := client.Run("arya-core", chart, v)

	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}

func RestartDeployment() {
	fmt.Println("Restarting deployment")
	cmd := exec.Command("bash", "-c",
		"kubectl rollout restart deployment arya-core-aryacore"+
			"-deployment")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
