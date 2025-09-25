package cmd

import (
	"fmt"
	"strings"

	"github.com/serptech/serp-go/api/client"
	"github.com/serptech/serp-go/api/common"
	"github.com/serptech/serp-go/api/origins"
	"github.com/spf13/cobra"
)

const (
	originsList   = "list"
	originsDelete = "delete"
	originsGet    = "get"
	originsUpdate = "update"
	originsCreate = "create"
)

var (
	originSearch      string
	originID          int
	originName        string
	originMinFacesize int
	originEntryDays   int
	originCreateMin   int
	originCreateHa    bool
	originCreateJunk  bool
	originIsActive    bool
)

var originsCmd = &cobra.Command{
	Use:       "origins [action]",
	Short:     "Manage origin configuration",
	Long:      "Provides helpers for listing and maintaining origin configuration in SerpTech.",
	ValidArgs: []string{originsList, originsDelete, originsGet, originsUpdate, originsCreate},
	Args:      cobra.MaximumNArgs(1),
	Example: `  serptech origins list --limit 20
  serptech origins get --id 3
  serptech origins update --id 3 --name "Lobby" --is-active=false
  serptech origins delete --id 7
  serptech origins create --name "Warehouse"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		action := args[0]
		c, err := client.NewClient()
		ifErrorExit(err)

		switch action {
		case originsList:
			query := common.NewSearchPaginationQuery(originSearch, limit, offset)
			resp, err := c.Origins().List(query)
			ifErrorExit(err)
			writeOutput(resp)
		case originsDelete:
			if originID == 0 {
				printAndExit("origin id is required")
			}
			ifErrorExit(c.Origins().Delete(originID))
			fmt.Printf("origin %d successfully deleted\n", originID)
		case originsGet:
			if originID == 0 {
				printAndExit("origin id is required")
			}
			resp, err := c.Origins().Get(originID)
			ifErrorExit(err)
			writeOutput(resp)
		case originsUpdate:
			if originID == 0 {
				printAndExit("origin id is required")
			}
			req := origins.UpdateRequest{ID: originID}
			if cmd.Flag("name").Changed {
				req.Name = stringPtr(strings.TrimSpace(originName))
			}
			if cmd.Flag("is-active").Changed {
				req.IsActive = boolPtr(originIsActive)
			}
			if cmd.Flag("min-facesize").Changed {
				req.MinFacesize = intPtr(originMinFacesize)
			}
			if cmd.Flag("entry-storage-days").Changed {
				req.EntryStorageDays = intPtr(originEntryDays)
			}
			if cmd.Flag("create-min-facesize").Changed {
				req.CreateMinFacesize = intPtr(originCreateMin)
			}
			if cmd.Flag("create-ha").Changed {
				req.CreateHa = boolPtr(originCreateHa)
			}
			if cmd.Flag("create-junk").Changed {
				req.CreateJunk = boolPtr(originCreateJunk)
			}

			if err := req.Validate(); err != nil {
				printAndExit(err.Error())
			}

			resp, err := c.Origins().Update(req)
			ifErrorExit(err)
			writeOutput(resp)
		case originsCreate:
			if strings.TrimSpace(originName) == "" {
				printAndExit("origin name is required")
			}
			req := origins.DefaultSourceWithName(strings.TrimSpace(originName))
			if cmd.Flag("is-active").Changed {
				req.IsActive = boolPtr(originIsActive)
			}
			if cmd.Flag("min-facesize").Changed {
				req.MinFacesize = intPtr(originMinFacesize)
			}
			if cmd.Flag("entry-storage-days").Changed {
				printAndExit("entry-storage-days is not supported during origin creation; use update instead")
			}
			if cmd.Flag("create-min-facesize").Changed {
				req.CreateMinFacesize = intPtr(originCreateMin)
			}
			if cmd.Flag("create-ha").Changed {
				req.CreateHa = boolPtr(originCreateHa)
			}
			if cmd.Flag("create-junk").Changed {
				req.CreateJunk = boolPtr(originCreateJunk)
			}

			resp, err := c.Origins().Create(req)
			ifErrorExit(err)
			writeOutput(resp)
		default:
			printAndExit(fmt.Sprintf("unsupported command %q", action))
		}
		return nil
	},
}

func init() {
	originIsActive = true
	originsCmd.Flags().StringVarP(&originSearch, "search", "s", "", "filtering by partially specified name")
	originsCmd.Flags().IntVar(&originID, "id", 0, "origin identifier")
	originsCmd.Flags().StringVar(&originName, "name", "", "origin name")
	originsCmd.Flags().BoolVar(&originIsActive, "is-active", true, "whether origin is active")
	originsCmd.Flags().IntVar(&originMinFacesize, "min-facesize", 0, "minimum facesize for uploads")
	originsCmd.Flags().IntVar(&originEntryDays, "entry-storage-days", 0, "number of days to keep entries")
	originsCmd.Flags().IntVar(&originCreateMin, "create-min-facesize", 0, "minimum facesize when creating profiles")
	originsCmd.Flags().BoolVar(&originCreateHa, "create-ha", false, "allow profile creation when confidence is HA")
	originsCmd.Flags().BoolVar(&originCreateJunk, "create-junk", false, "allow profile creation when confidence is junk")

	rootCmd.AddCommand(originsCmd)
}
