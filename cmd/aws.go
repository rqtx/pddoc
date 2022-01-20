/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/rqtx/pddoc/utils"
	"github.com/rqtx/pddoc/workers/aws"

	"github.com/spf13/cobra"
)

var file string
var region string

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "AWS doc",
	Long:  `Create AWS documentation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		worker := aws.NewWorker(region)
		doc := utils.NewDocumet(worker.GetSetctions())
		return doc.WriteDocument(file)
	},
}

func init() {
	rootCmd.AddCommand(awsCmd)

	awsCmd.Flags().StringVarP(&file, "file", "f", "", "file is required")
	awsCmd.MarkFlagRequired("file")

	awsCmd.Flags().StringVarP(&region, "region", "r", "", "region is required")
	awsCmd.MarkFlagRequired("region")
}
