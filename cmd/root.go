/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aryacore",
	Short: "Arya Core - the high performance time series engine for distributed hardware systems.",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	configureFlags()
	cobra.OnInitialize(initConfig)
}

type configureFlag func()

// ||| FLAG CONFIGURATION ||||

func configureFlags() {
	flags := []configureFlag{
		configureConfigFlag,
	}
	for _, cf := range flags {
		cf()
	}
}

// || CONFIG ||

const (
	envPrefix      = "arya"
	configFlag     = "config"
	configType     = "yaml"
	configRelPath  = ".arya"
	configFileName = "config"
)

func configureConfigFlag() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, configFlag, "", "config file (default is $HOME/.aryacore.yaml)")
}

func configName() string {
	return configRelPath + "/" + configFileName
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType(configType)
		viper.SetConfigName(configName())
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
