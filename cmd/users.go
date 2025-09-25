package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/serptech/serp-go/api/client"
	"github.com/serptech/serp-go/api/common"
	serpusers "github.com/serptech/serp-go/api/users"
	"github.com/spf13/cobra"
)

const (
	userMe         = "me"
	userListTokens = "list-tokens"
	userStatistics = "statistics"
	userList       = "list"
	userGet        = "get"
	userUpdate     = "update"
	userPatch      = "patch"
)

var usersCmd = &cobra.Command{
	Use:       "users [action]",
	Short:     "Interact with user-related endpoints",
	ValidArgs: []string{userMe, userListTokens, userStatistics, userList, userGet, userUpdate, userPatch},
	Args:      cobra.MaximumNArgs(1),
	Long:      "Available actions: me, list-tokens, statistics, list, get, update, patch.",
	Example: `  serptech users me
  serptech users list-tokens
  serptech users statistics
  serptech users list --limit 50 --query admin
  serptech users get --id 12
  serptech users update --id 12 --username admin --is-active=false
  serptech users patch --id 12 --is-active=true`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		action := args[0]
		c := resolveUsersClient(action)

		switch action {
		case userMe:
			meOut, err := c.Users().Me()
			ifErrorExit(err)
			writeOutput(meOut)
		case userListTokens:
			tokens, err := c.Tokens().ListAccess(nil)
			ifErrorExit(err)
			writeOutput(tokens)
		case userStatistics:
			resp, err := c.Users().Statistics()
			ifErrorExit(err)
			writeOutput(resp)
		case userList:
			query := common.NewPaginationQuery(limit, offset)
			if cmd.Flag("query").Changed {
				trimmed := strings.TrimSpace(userQueryValue)
				if trimmed != "" {
					query["q"] = trimmed
				}
			}
			resp, err := c.Users().List(query)
			ifErrorExit(err)
			writeOutput(resp)
		case userGet:
			if userTargetID == 0 {
				printAndExit("user id is required")
			}
			resp, err := c.Users().Get(userTargetID)
			ifErrorExit(err)
			writeOutput(resp)
		case userUpdate:
			if userTargetID == 0 {
				printAndExit("user id is required")
			}
			if !cmd.Flag("username").Changed {
				printAndExit("username is required")
			}
			if !cmd.Flag("is-active").Changed {
				printAndExit("is-active is required")
			}
			req := serpusers.UpdateUserRequest{
				Username: strings.TrimSpace(userUsernameValue),
				IsActive: userIsActiveValue,
			}
			if err := req.Validate(); err != nil {
				printAndExit(err.Error())
			}
			resp, err := c.Users().Update(userTargetID, req)
			ifErrorExit(err)
			writeOutput(resp)
		case userPatch:
			if userTargetID == 0 {
				printAndExit("user id is required")
			}
			var req serpusers.PartialUpdateUserRequest
			if cmd.Flag("username").Changed {
				trimmed := strings.TrimSpace(userUsernameValue)
				req.Username = &trimmed
			}
			if cmd.Flag("is-active").Changed {
				active := userIsActiveValue
				req.IsActive = &active
			}
			if err := req.Validate(); err != nil {
				printAndExit(err.Error())
			}
			resp, err := c.Users().PartialUpdate(userTargetID, req)
			ifErrorExit(err)
			writeOutput(resp)
		default:
			printAndExit(fmt.Sprintf("unsupported command %q", action))
		}
		return nil
	},
}

var (
	userQueryValue    string
	userTargetID      int
	userUsernameValue string
	userIsActiveValue bool
)

func resolveUsersClient(action string) *client.Client {
	requiresRoot := action == userList || action == userGet || action == userUpdate || action == userPatch
	if requiresRoot {
		rootToken := strings.TrimSpace(os.Getenv("SERP_ROOT_TOKEN"))
		if rootToken == "" {
			printAndExit("SERP_ROOT_TOKEN environment variable is required for this action")
		}
		return client.NewClientWithToken(rootToken)
	}
	c, err := client.NewClient()
	ifErrorExit(err)
	return c
}

func init() {
	usersCmd.Flags().StringVar(&userQueryValue, "query", "", "filter users by username substring")
	usersCmd.Flags().IntVar(&userTargetID, "id", 0, "target user identifier")
	usersCmd.Flags().StringVar(&userUsernameValue, "username", "", "username value")
	usersCmd.Flags().BoolVar(&userIsActiveValue, "is-active", false, "toggle active status")
	rootCmd.AddCommand(usersCmd)
}
