package cmd

import (
	"github.com/serptech/serp-go/api/client"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display API version information",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient()
		ifErrorExit(err)

		resp, err := c.Users().Version()
		ifErrorExit(err)

		writeOutput(resp)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
