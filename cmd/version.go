/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version string = "0.0.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `Print the version number of pdoc`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
