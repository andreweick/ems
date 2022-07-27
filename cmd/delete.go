/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/missionfocus/ems/pkg/cloudflare"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete photos or videos",
	Long: `Delete content from Eick.com Managment System.  Delete from both the blob and database store.

$ ems delete ~/edc/photomechanic/file1.jpg  # Delete file1.jpg from cloudflare
$ ems delete ~/edc/photomechanic/           # Delete all files from cloudflare
`,
	Run: func(cmd *cobra.Command, args []string) {
		feedDelete(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func feedDelete(filename string) {
	cloudflare.Delete(filename)
}
