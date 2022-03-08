/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aryacore",
	Short: "Arya Core - the high performance time series engine for distributed hardware systems.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	configureRootFlags()
	cobra.OnInitialize(initConfig)
}

type configureFlag func()

// ||| FLAG CONFIGURATION ||||

func configureRootFlags() {
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
	configType     = "json"
	configRelPath  = ".arya"
	configFileName = "core-config"
)

func configureConfigFlag() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, configFlag, "", "config file (default is $HOME/.arya/core-config.json)")
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
