/*
Copyright © 2022 Arya Analytics

*/

package cmd

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development Utilities",
	Long: `Provides access to development utilities such as CLuster provisioning, 
			tooling installs, and configuration management`,
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(clusterCmd)
	configureClusterCommand()
	devCmd.AddCommand(toolingCmd)
	configureToolingCmd()
	devCmd.AddCommand(dockerCmd)
	configureDockerCmd()
	devCmd.AddCommand(loginCmd)
	devCmd.AddCommand(reloaderCmd)
	configureReloaderCmd()
}

// |||| CLUSTER |||

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Provision and manage development clusters",
}

var (
	devClusterNodes      int
	devClusterCores      int
	devClusterMemory     int
	devClusterStorage    int
	devClusterCIDROffset int
	devClusterReInit     bool
	devClusterName       string
)

func configureClusterCommand() {
	clusterCmd.AddCommand(clusterProvisionCmd)
	clusterCmd.AddCommand(clusterDeleteCmd)
	clusterCmd.Flags().IntVarP(&devClusterNodes, "nodes", "n", 3, "Number of nodes in cluster")
	clusterCmd.Flags().IntVarP(&devClusterCores, "cores", "c", dev.BaseAryaClusterCfg.Cores, "Number of cores per node")
	clusterCmd.Flags().IntVarP(&devClusterMemory, "memory", "m", dev.BaseAryaClusterCfg.Memory, "Amount of memory (gb) per node")
	clusterCmd.Flags().IntVarP(&devClusterStorage, "storage", "s", dev.BaseAryaClusterCfg.Storage, "Amount of storage (gb) per node")
	clusterCmd.Flags().IntVarP(&devClusterCIDROffset, "cidrOffset", "o", dev.BaseAryaClusterCfg.CidrOffset, "StructValue to offset Cluster node cidrs by (ex. "+
		"an offset of 10 would make the first nodes cidr 10.11.0.0/16)")
	clusterCmd.Flags().BoolVarP(&devClusterReInit, "reInit", "r", dev.BaseAryaClusterCfg.ReInit, "Whether to delete and replace an existing cluster")
	clusterCmd.Flags().StringVarP(&devClusterName, "name", "d", dev.BaseAryaClusterCfg.Name, "name of cluster")
}

var clusterProvisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision an Arya development cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := dev.ProvisionLocalDevCluster(devClusterNodes, devClusterName, devClusterCores, devClusterMemory, devClusterStorage, devClusterReInit, devClusterCIDROffset)
		return err
	},
}

var clusterDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Arya development cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.DeleteLocalDevCluster(devClusterName)
	},
}

// |||| TOOLING ||||

var toolingCmd = &cobra.Command{
	Use:   "tools",
	Short: "Install and manage development tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.InstallRequiredTools()
	},
}

func configureToolingCmd() {
	toolingCmd.AddCommand(toolingUninstallCmd)
	toolingCmd.AddCommand(toolingCheckCmd)
	toolingCmd.AddCommand(toolingInstallCmd)
}

var toolingInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install development tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.InstallRequiredTools()
	},
}

var toolingUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall development tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.UninstallRequiredTools()
	},
}

var toolingCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if required development tools are installed",
	RunE: func(cmd *cobra.Command, args []string) error {
		dev.RequiredToolsInstalled()
		return nil
	},
}

// |||| LOGIN ||||

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Github and register credentials in config",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.Login(args[0])
	},
}

// |||| DOCKER ||||
func configureDockerCmd() {
	dockerCmd.AddCommand(dockerStartCmd)
}

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Operate development docker containers",
}

var dockerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the development docker containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.StartDockerCompose()
	},
}

// |||| RELOADER ||||

var (
	buildCtx string
)

func configureReloaderCmd() {
	reloaderCmd.AddCommand(reloaderStartCmd)
	reloaderStartCmd.Flags().StringVarP(&buildCtx, "buildCtx", "b", dev.DefaultBuildCtxPath(), "build context")
}

var reloaderCmd = &cobra.Command{
	Use:   "reloader",
	Short: "operate the development hot-reloader",
}

var reloaderStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the hot-reloader",
	RunE: func(cmd *cobra.Command, args []string) error {
		dev.StartReloader(devClusterName, buildCtx)
		return nil
	},
}
