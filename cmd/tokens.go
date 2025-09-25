package cmd

import (
	"fmt"
	"strings"

	"github.com/serptech/serp-go/api/client"
	"github.com/serptech/serp-go/api/common"
	"github.com/serptech/serp-go/api/tokens"
	"github.com/spf13/cobra"
)

var (
	tokensAccessPermanent   bool
	tokensAccessKey         string
	tokensAccessFilterSpace int

	tokensStreamPermanent   bool
	tokensStreamKey         string
	tokensStreamFilterSpace int
)

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Manage API tokens",
}

var tokensAccessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage access tokens",
}

var tokensAccessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List access tokens",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient()
		ifErrorExit(err)

		query := common.NewPaginationQuery(limit, offset)
		if cmd.Flag("space-id").Changed {
			query["space_id"] = tokensAccessFilterSpace
		}

		resp, err := c.Tokens().ListAccess(query)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var tokensAccessCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create access token",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient()
		ifErrorExit(err)

		req := tokens.CreateTokenRequest{Permanent: tokensAccessPermanent}

		resp, err := c.Tokens().CreateAccess(req)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var tokensAccessDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete access token",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(tokensAccessKey) == "" {
			printAndExit("token key is required")
		}
		c, err := client.NewClient()
		ifErrorExit(err)
		ifErrorExit(c.Tokens().DeleteAccess(tokensAccessKey))
		fmt.Printf("access token %s successfully deleted\n", tokensAccessKey)
	},
}

var tokensStreamsCmd = &cobra.Command{
	Use:   "streams",
	Short: "Manage stream tokens",
}

var tokensStreamsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stream tokens",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient()
		ifErrorExit(err)

		query := common.NewPaginationQuery(limit, offset)
		if cmd.Flag("space-id").Changed {
			query["space_id"] = tokensStreamFilterSpace
		}

		resp, err := c.Tokens().ListStreams(query)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var tokensStreamsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create stream token",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.NewClient()
		ifErrorExit(err)

		req := tokens.CreateTokenRequest{Permanent: tokensStreamPermanent}

		resp, err := c.Tokens().CreateStream(req)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var tokensStreamsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete stream token",
	Run: func(cmd *cobra.Command, args []string) {
		if strings.TrimSpace(tokensStreamKey) == "" {
			printAndExit("token key is required")
		}
		c, err := client.NewClient()
		ifErrorExit(err)
		ifErrorExit(c.Tokens().DeleteStream(tokensStreamKey))
		fmt.Printf("stream token %s successfully deleted\n", tokensStreamKey)
	},
}

func init() {
	tokensAccessListCmd.Flags().IntVar(&tokensAccessFilterSpace, "space-id", 0, "filter by space identifier")

	tokensAccessCreateCmd.Flags().BoolVar(&tokensAccessPermanent, "permanent", false, "create permanent token")

	tokensAccessDeleteCmd.Flags().StringVar(&tokensAccessKey, "key", "", "token key")

	tokensStreamsListCmd.Flags().IntVar(&tokensStreamFilterSpace, "space-id", 0, "filter by space identifier")

	tokensStreamsCreateCmd.Flags().BoolVar(&tokensStreamPermanent, "permanent", false, "create permanent token")

	tokensStreamsDeleteCmd.Flags().StringVar(&tokensStreamKey, "key", "", "token key")

	tokensAccessCmd.AddCommand(tokensAccessListCmd, tokensAccessCreateCmd, tokensAccessDeleteCmd)
	tokensStreamsCmd.AddCommand(tokensStreamsListCmd, tokensStreamsCreateCmd, tokensStreamsDeleteCmd)
	tokensCmd.AddCommand(tokensAccessCmd, tokensStreamsCmd)
	rootCmd.AddCommand(tokensCmd)
}
