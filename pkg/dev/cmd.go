package dev

import "C"
import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
)

var Cmd = &cli.Command{
	Name: "dev",
	Usage: "Provides access to development services such as cluster provisioning, " +
		"tooling installs, and configuration management.",
	Subcommands: []*cli.Command{
		clusterCmd,
		toolingCmd,
		reloaderCmd,
		configCmd,
		loginCmd,
	},
}

var clusterCmd = &cli.Command{
	Name:  "cluster",
	Usage: "Provision and manage development nodes.",
	Subcommands: []*cli.Command{
		{
			Name:  "init",
			Usage: "Initialize an Arya development cluster.",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "vms",
					Aliases: []string{"n"},
					Value:   3,
					Usage:   "Number of vms in cluster",
				},
				&cli.IntFlag{
					Name:    "cores",
					Aliases: []string{"c"},
					Value:   BaseAryaClusterCfg.Cores,
					Usage:   "Number of cores per node",
				},
				&cli.IntFlag{
					Name:    "memory",
					Aliases: []string{"m"},
					Value:   BaseAryaClusterCfg.Memory,
					Usage:   "Amount of memory (gb) per node",
				},
				&cli.IntFlag{
					Name:    "storage",
					Aliases: []string{"s"},
					Value:   BaseAryaClusterCfg.Storage,
					Usage:   "Amount of storage (gb) per node",
				},
				&cli.IntFlag{
					Name: "cidrOffset",
					Aliases: []string{"co"},
					Value: BaseAryaClusterCfg.CidrOffset,
					Usage: "Value to offset cluster node cidrs by (ex. " +
						"an offset of 10 would make the first nodes cidr 10.11.0.0/16",
				},
				&cli.BoolFlag{
					Name:    "reInit",
					Aliases: []string{"r"},
					Value:   BaseAryaClusterCfg.ReInit,
					Usage:   "Whether to delete existing cluster infrastructure",
				},
				&cli.StringFlag{
					Name:    "clusterName",
					Aliases: []string{"cn"},
					Value:   BaseAryaClusterCfg.Name,
					Usage:   "Name of Arya cluster",
				},

			},
			Action: func(c *cli.Context) error {
				if err := InstallRequired(); err != nil {
					return err
				}
				numNodes := c.Int("vms")
				clusterName := c.String("clusterName")
				fmt.Printf("%s Initializing an Arya Cluster named %s with %v nodes \n",
					emoji("\\U0001F4AB"),clusterName, numNodes)
				aryaCfg := AryaClusterConfig{
					NumNodes: numNodes,
					Cores:    c.Int("cores"),
					Memory:   c.Int("memory"),
					Storage:  c.Int("storage"),
					ReInit:   c.Bool("reInit"),
					CidrOffset: c.Int("cidrOffset"),
				}
				aryaCluster := NewAryaCluster(aryaCfg)
				if err := aryaCluster.Provision(); err != nil {
					return err
				}
				aryaConfig := NewAryaConfig()
				fmt.Println("Merging kubeconfig")
				fmt.Println(aryaCluster.nodes)
				for _, c := range aryaCluster.nodes {
					if err := aryaConfig.MergeRemoteKubeConfig(*c); err != nil {
						return err
					}
				}
				fmt.Printf("Successfully initialized Arya Cluster %s", clusterName)
				return nil
			},
		},
		{
			Name: "stop",
			Usage: "Stop an Arya development cluster",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "name",
					Aliases: []string{"cn"},
					Usage: "Name of Arya cluster",
				},
			},
			Action: func(c *cli.Context) error {

				return nil
			},
		},
	},
}

func emoji(s string) string {
	r, _ := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(r)
}

// toolingCmd dispatches cli actions to the correct tooling function.
var toolingCmd = &cli.Command{
	Name:  "tooling",
	Usage: "Install and manage development tools.",
	Subcommands: []*cli.Command{
		{
			Name:  "install",
			Usage: "Install development tools.",
			Action: func(c *cli.Context) error {
				return InstallRequired()
			},
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall development tools.",
			Action: func(c *cli.Context) error {
				return UninstallRequired()
			},
		},
	},
}

// reloaderCmd dispatches cli actions to the correct tooling function.
var reloaderCmd = &cli.Command{
	Name:  "reloader",
	Usage: "Operate the development hot-reloader.",
	Action: func(c *cli.Context) error {
		return nil
	},
}

// configCmd dispatches cli actions to the correct tooling function
var configCmd = &cli.Command{
	Name:  "config",
	Usage: "Manage development configurations.",
}

var loginCmd = &cli.Command{
	Name:  "login",
	Usage: "Login to Github and register credentials in config",
}
