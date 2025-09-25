package cmd

import (
	"fmt"
	"strings"

	"github.com/serptech/serp-go/api/client"
	"github.com/serptech/serp-go/api/common"
	"github.com/serptech/serp-go/api/profiles"
	"github.com/spf13/cobra"
)

const (
	profileCreate = "create"
	profileSearch = "search"
	profileDelete = "delete"
	profileReinit = "reinit"
)

var (
	primaryPhotoPath   string
	secondaryPhotoPath string
	profileOriginID    int
	profileID          string
	profileMinFacesize int
	profileMinConf     int
	profileAllowHa     bool
	profileAllowJunk   bool
)

var profilesCmd = &cobra.Command{
	Use:       "profiles [command]",
	Short:     "Manage recognition profiles",
	Long:      "Provides helpers for working with profile lifecycle using the SerpTech API.",
	Example:   "serptech profiles create --photo img/profile.png --origin-id 42",
	ValidArgs: []string{profileCreate, profileSearch, profileDelete, profileReinit},
	Args:      cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		c, err := client.NewClient()
		ifErrorExit(err)

		switch action {
		case profileCreate:
			handleProfileCreate(cmd, c)
		case profileSearch:
			handleProfileSearch(cmd, c)
		case profileDelete:
			handleProfileDelete(c)
		case profileReinit:
			handleProfileReinit(cmd, c)
		default:
			printAndExit(fmt.Sprintf("unsupported command %q", action))
		}
	},
}

func handleProfileCreate(cmd *cobra.Command, c *client.Client) {
	if primaryPhotoPath == "" {
		printAndExit("photo is required")
	}
	if profileOriginID == 0 {
		printAndExit("origin-id is required")
	}

	photo, err := common.NewPhotoFromFile(primaryPhotoPath)
	ifErrorExit(err)

	req := profiles.CreateRequest{
		Photo:    photo,
		OriginID: profileOriginID,
	}

	if cmd.Flag("create-min-facesize").Changed {
		req.CreateMinFacesize = intPtr(profileMinFacesize)
	}
	if cmd.Flag("create-ha").Changed {
		req.CreateHa = boolPtr(profileAllowHa)
	}
	if cmd.Flag("create-junk").Changed {
		req.CreateJunk = boolPtr(profileAllowJunk)
	}

	resp, err := c.Profiles().Create(req)
	ifErrorExit(err)
	writeOutput(resp)
}

func handleProfileSearch(cmd *cobra.Command, c *client.Client) {
	if primaryPhotoPath == "" {
		printAndExit("photo is required")
	}

	primary, err := common.NewPhotoFromFile(primaryPhotoPath)
	ifErrorExit(err)

	searchReq := profiles.SearchRequest{Photo: primary}
	if secondaryPhotoPath != "" {
		secondary, err := common.NewPhotoFromFile(secondaryPhotoPath)
		ifErrorExit(err)
		searchReq.SecondImageName = secondary.PhotoName
		searchReq.SecondImageData = secondary.PhotoData
	}

	resp, err := c.Profiles().Search(searchReq)
	ifErrorExit(err)
	writeOutput(resp)
}

func handleProfileDelete(c *client.Client) {
	if strings.TrimSpace(profileID) == "" {
		printAndExit("profile-id is required")
	}
	ifErrorExit(c.Profiles().Delete(profileID))
	fmt.Printf("profile %s successfully deleted\n", profileID)
}

func handleProfileReinit(cmd *cobra.Command, c *client.Client) {
	if strings.TrimSpace(profileID) == "" {
		printAndExit("profile-id is required")
	}
	if primaryPhotoPath == "" {
		printAndExit("photo is required")
	}

	photo, err := common.NewPhotoFromFile(primaryPhotoPath)
	ifErrorExit(err)

	req := profiles.ReinitRequest{Photo: photo}
	if cmd.Flag("create-min-facesize").Changed {
		req.CreateMinFacesize = intPtr(profileMinFacesize)
	}
	if cmd.Flag("min-conf").Changed {
		req.MinConf = intPtr(profileMinConf)
	}

	resp, err := c.Profiles().Reinit(profileID, req)
	ifErrorExit(err)
	if len(resp) == 0 {
		fmt.Println("reinit completed")
		return
	}
	writeOutput(resp)
}

func init() {
	profilesCmd.Flags().StringVarP(&primaryPhotoPath, "photo", "p", "", "path to primary photo for create/search/reinit")
	profilesCmd.Flags().StringVar(&secondaryPhotoPath, "second-photo", "", "optional path to secondary photo for search")
	profilesCmd.Flags().IntVar(&profileOriginID, "origin-id", 0, "origin identifier for profile create")
	profilesCmd.Flags().StringVar(&profileID, "profile-id", "", "profile identifier for delete or reinit")
	profilesCmd.Flags().IntVar(&profileMinFacesize, "create-min-facesize", 0, "minimum face size when creating or reinitializing")
	profilesCmd.Flags().IntVar(&profileMinConf, "min-conf", 0, "minimum match confidence for reinit")
	profilesCmd.Flags().BoolVar(&profileAllowHa, "create-ha", false, "allow creation when result confidence is HA")
	profilesCmd.Flags().BoolVar(&profileAllowJunk, "create-junk", false, "allow creation when result confidence is junk")

	rootCmd.AddCommand(profilesCmd)
}
