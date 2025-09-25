package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/serptech/serp-go/api/client"
	"github.com/serptech/serp-go/api/common"
	"github.com/serptech/serp-go/api/entries"
	"github.com/spf13/cobra"
)

var (
	entriesListOriginIDs string
	entriesListSpaceIDs  string
	entriesListPersonIDs string
	entriesListConf      string
	entriesListDateFrom  string
	entriesListDateTo    string

	entriesDeleteID int

	entriesStatsPersonIDs     string
	entriesStatsSourceID      int
	entriesStatsConfValue     string
	entriesStatsLivenessValue string
	entriesStatsEntryIDFrom   int
	entriesStatsDateFrom      string
	entriesStatsDateTo        string
)

var entriesCmd = &cobra.Command{
	Use:   "entries",
	Short: "Inspect recognition entries and statistics",
}

var entriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recognition entries",
	Run: func(cmd *cobra.Command, args []string) {
		c := resolveEntriesClient(false)

		query := common.NewPaginationQuery(limit, offset)
		if cmd.Flag("origin-ids").Changed {
			query["origin_ids"] = strings.TrimSpace(entriesListOriginIDs)
		}
		if cmd.Flag("spaces-ids").Changed {
			query["spaces_ids"] = strings.TrimSpace(entriesListSpaceIDs)
		}
		if cmd.Flag("person-ids").Changed {
			query["person_ids"] = strings.TrimSpace(entriesListPersonIDs)
		}
		if cmd.Flag("conf").Changed {
			query["conf"] = strings.TrimSpace(entriesListConf)
		}
		if cmd.Flag("date-from").Changed {
			trimmed := strings.TrimSpace(entriesListDateFrom)
			if trimmed != "" {
				parsed, err := parseDate(trimmed)
				ifErrorExit(err)
				query["date_from"] = parsed.Format(time.RFC3339)
			}
		}
		if cmd.Flag("date-to").Changed {
			trimmed := strings.TrimSpace(entriesListDateTo)
			if trimmed != "" {
				parsed, err := parseDate(trimmed)
				ifErrorExit(err)
				query["date_to"] = parsed.Format(time.RFC3339)
			}
		}

		resp, err := c.Entries().List(query)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var entriesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete entry by identifier",
	Run: func(cmd *cobra.Command, args []string) {
		if entriesDeleteID == 0 {
			printAndExit("entry id is required")
		}
		c := resolveEntriesClient(false)
		ifErrorExit(c.Entries().Delete(entriesDeleteID))
		fmt.Printf("entry %d successfully deleted\n", entriesDeleteID)
	},
}

var entriesStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Retrieve entry statistics",
}

var entriesStatsSourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Show statistics grouped by origins",
	Run: func(cmd *cobra.Command, args []string) {
		var req entries.StatsSourcesRequest
		c := resolveEntriesClient(false)

		if cmd.Flag("person-ids").Changed {
			req.PersonIDs = strings.TrimSpace(entriesStatsPersonIDs)
		}
		if cmd.Flag("conf").Changed {
			parsedConf, err := resolveConf(entriesStatsConfValue)
			ifErrorExit(err)
			req.Conf = parsedConf
		}
		if cmd.Flag("liveness").Changed {
			parsedLiveness, err := resolveLiveness(entriesStatsLivenessValue)
			ifErrorExit(err)
			req.Liveness = parsedLiveness
		}
		if cmd.Flag("source-id").Changed {
			req.Source = entriesStatsSourceID
		}
		if cmd.Flag("entry-id-from").Changed {
			req.EntryIdFrom = entriesStatsEntryIDFrom
		}
		if cmd.Flag("date-from").Changed {
			trimmed := strings.TrimSpace(entriesStatsDateFrom)
			if trimmed != "" {
				parsed, err := parseDate(trimmed)
				ifErrorExit(err)
				req.DateFrom = parsed
			}
		}
		if cmd.Flag("date-to").Changed {
			trimmed := strings.TrimSpace(entriesStatsDateTo)
			if trimmed != "" {
				parsed, err := parseDate(trimmed)
				ifErrorExit(err)
				req.DateTo = parsed
			}
		}

		resp, err := c.Entries().StatsSources(req)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

func init() {
	entriesListCmd.Flags().StringVar(&entriesListOriginIDs, "origin-ids", "", "comma-separated list of origin identifiers")
	entriesListCmd.Flags().StringVar(&entriesListSpaceIDs, "spaces-ids", "", "comma-separated list of space identifiers")
	entriesListCmd.Flags().StringVar(&entriesListPersonIDs, "person-ids", "", "comma-separated list of person identifiers")
	entriesListCmd.Flags().StringVar(&entriesListConf, "conf", "", "comma-separated list of confidence values")
	entriesListCmd.Flags().StringVar(&entriesListDateFrom, "date-from", "", "filter entries created after the given datetime (RFC3339)")
	entriesListCmd.Flags().StringVar(&entriesListDateTo, "date-to", "", "filter entries created before the given datetime (RFC3339)")

	entriesDeleteCmd.Flags().IntVar(&entriesDeleteID, "id", 0, "entry identifier")

	entriesStatsSourcesCmd.Flags().StringVar(&entriesStatsPersonIDs, "person-ids", "", "comma-separated list of person identifiers")
	entriesStatsSourcesCmd.Flags().StringVar(&entriesStatsConfValue, "conf", "", "filter by confidence (name or integer value)")
	entriesStatsSourcesCmd.Flags().StringVar(&entriesStatsLivenessValue, "liveness", "", "filter by liveness (passed|failed|undetermined)")
	entriesStatsSourcesCmd.Flags().IntVar(&entriesStatsSourceID, "source-id", 0, "filter by origin identifier")
	entriesStatsSourcesCmd.Flags().IntVar(&entriesStatsEntryIDFrom, "entry-id-from", 0, "filter entries starting from identifier")
	entriesStatsSourcesCmd.Flags().StringVar(&entriesStatsDateFrom, "date-from", "", "filter by start date (YYYY-MM-DD or RFC3339)")
	entriesStatsSourcesCmd.Flags().StringVar(&entriesStatsDateTo, "date-to", "", "filter by end date (YYYY-MM-DD or RFC3339)")

	entriesStatsCmd.AddCommand(entriesStatsSourcesCmd)
	entriesCmd.AddCommand(entriesListCmd, entriesDeleteCmd, entriesStatsCmd)
	rootCmd.AddCommand(entriesCmd)
}

func resolveEntriesClient(requireRoot bool) *client.Client {
	if requireRoot {
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
