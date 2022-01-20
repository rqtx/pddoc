/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version string = "0.0.1"

var short bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `Print the version number of pdoc`,
	Run: func(cmd *cobra.Command, args []string) {
		if short {
			fmt.Printf("%s\n", version)
		} else {
			fmt.Printf("version %s\n", version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVarP(&short, "short", "s", false, "Output just version number")
}
