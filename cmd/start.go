// Copyright © 2022 Arya Analytics

package cmd

import (
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a node in a cluster",
	Long:  "Start an Arya Core node",
	Args:  cobra.NoArgs,
	RunE:  runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, _ []string) error {
	return nil
}
