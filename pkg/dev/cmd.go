package dev

import "C"
import (
	"github.com/urfave/cli/v2"
)

// || GENERAL CLI ||

var Cmd = &cli.Command{
	Name: "dev",
	Usage: "Provides access to development services such as Cluster provisioning, " +
		"tools installs, and configuration management.",
	Subcommands: []*cli.Command{
		clusterCmd,
		toolingCmd,
		reloaderCmd,
		loginCmd,
	},
}

// || LOCAL DEV CLUSTER CLI ||

var clusterCmd = &cli.Command{
	Name:  "cluster",
	Usage: "Provision and manage development clusters.",
	Subcommands: []*cli.Command{
		{
			Name:  "provision",
			Usage: "Provision an Arya development Cluster.",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "nodes",
					Aliases: []string{"n"},
					Value:   3,
					Usage:   "Number of vms in Cluster",
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
					Name:    "cidrOffset",
					Aliases: []string{"co"},
					Value:   BaseAryaClusterCfg.CidrOffset,
					Usage: "Value to offset Cluster node cidrs by (ex. " +
						"an offset of 10 would make the first nodes cidr 10.11.0.0/16)",
				},
				&cli.BoolFlag{
					Name:    "reInit",
					Aliases: []string{"r"},
					Value:   BaseAryaClusterCfg.ReInit,
					Usage:   "Whether to delete existing Cluster infrastructure",
				},
				&cli.StringFlag{
					Name:    "Name",
					Aliases: []string{"cn"},
					Value:   BaseAryaClusterCfg.Name,
					Usage:   "Name of Arya Cluster",
				},
			},
			Action: func(c *cli.Context) error {
				_, err := ProvisionLocalDevCluster(
					c.Int("nodes"),
					c.String("Name"),
					c.Int("cores"),
					c.Int("memory"),
					c.Int("storage"),
					c.Bool("reInit"),
					c.Int("cidrOffset"),
				)
				return err
			},
		},
		{
			Name:  "delete",
			Usage: "Delete an Arya development Cluster",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "Name",
					Aliases: []string{"cn"},
					Value:   BaseAryaClusterCfg.Name,
					Usage:   "Name of Arya Cluster",
				},
			},
			Action: func(c *cli.Context) error {
				return DeleteLocalDevCluster(c.String("Name"))
			},
		},
	},
}

// || TOOLING CLI ||

var toolingCmd = &cli.Command{
	Name:  "tools",
	Usage: "Install and manage development tools.",
	Subcommands: []*cli.Command{
		{
			Name:  "install",
			Usage: "Install development tools.",
			Action: func(c *cli.Context) error {
				return InstallRequiredTools()
			},
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall development tools.",
			Action: func(c *cli.Context) error {
				return UninstallRequiredTools()
			},
		},
		{
			Name:  "check",
			Usage: "Check if required development tools are installed",
			Action: func(c *cli.Context) error {
				RequiredToolsInstalled()
				return nil
			},
		},
	},
}

// || RELOADER CLI ||

var reloaderCmd = &cli.Command{
	Name:  "reloader",
	Usage: "Operate the development hot-reloader.",
	Subcommands: []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the hot-reloader",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "clusterName",
					Aliases: []string{"cn"},
					Value:   BaseAryaClusterCfg.Name,
					Usage:   "Name of arya Cluster to deploy into",
				},
				&cli.StringFlag{
					Name:    "buildCtx",
					Aliases: []string{"bc"},
					Value:   DefaultBuildCtxPath(),
				},
			},
			Action: func(c *cli.Context) error {
				StartReloader(c.String("clusterName"), c.String("buildCtx"))
				return nil
			},
		},
	},
}

// || LOGIN CLI ||

var loginCmd = &cli.Command{
	Name:  "login",
	Usage: "Login to Github and register credentials in config",
	Action: func(c *cli.Context) error {
		return Login(c.Args().Get(0))
	},
}
