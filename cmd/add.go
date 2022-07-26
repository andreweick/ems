/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/missionfocus/ems/pkg/cloudflare"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add/Update photographs or videos",
	Long: `
Add photographs or videos to Eick.com management system..
$ ems add ~/edc/photomechanic/file1.jpg
$ ems add ~/edc/photomechanic/*.jpg # Add all files in ~/edc/photomechanic
`,
	Run: func(cmd *cobra.Command, args []string) {
		feedAdd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func feedAdd(filename string) {
	cloudflare.Add(filename)
}
