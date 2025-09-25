package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/serptech/serp-go/api/client"
	"github.com/serptech/serp-go/api/const/conf"
	"github.com/serptech/serp-go/api/utility"
	"github.com/spf13/cobra"
)

var (
	utilityAsmPhotoPath       string
	utilityLivenessPhoto1Path string
	utilityLivenessPhoto2Path string

	utilityComparePhoto1      string
	utilityComparePhoto2      string
	utilityCompareConfValue   string
	utilityCompareLivenessOne bool
	utilityCompareLivenessTwo bool
)

var utilityCmd = &cobra.Command{
	Use:   "utility",
	Short: "Access utility endpoints (health, metrics, checks)",
}

var utilityHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Health check",
	Run: func(cmd *cobra.Command, args []string) {
		c := resolveUtilityClient(false)
		resp, err := c.Utility().Health()
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var utilityMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Platform metrics",
	Run: func(cmd *cobra.Command, args []string) {
		c := resolveUtilityClient(false)
		metrics, err := c.Utility().Metrics()
		ifErrorExit(err)
		fmt.Print(metrics)
	},
}

var utilityAsmCmd = &cobra.Command{
	Use:   "asm",
	Short: "Age/sex/mood prediction",
	Run: func(cmd *cobra.Command, args []string) {
		if utilityAsmPhotoPath == "" {
			printAndExit("photo is required")
		}
		c := resolveUtilityClient(false)
		req, err := utility.NewAsmRequest(utilityAsmPhotoPath)
		ifErrorExit(err)
		resp, err := c.Utility().Asm(req)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var utilityLivenessCmd = &cobra.Command{
	Use:   "liveness",
	Short: "Run liveness check",
	Run: func(cmd *cobra.Command, args []string) {
		if utilityLivenessPhoto1Path == "" {
			printAndExit("photo1 is required, supply --photo1 /path/to/image")
		}
		if utilityLivenessPhoto2Path == "" {
			printAndExit("photo2 is required, supply --photo2 /path/to/image")
		}
		c := resolveUtilityClient(false)
		req, err := utility.NewLivenessRequest(utilityLivenessPhoto1Path, utilityLivenessPhoto2Path)
		ifErrorExit(err)
		resp, err := c.Utility().Liveness(req)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

var utilityCompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two faces (access token only)",
	Run: func(cmd *cobra.Command, args []string) {
		if utilityComparePhoto1 == "" {
			printAndExit("photo1 is required, supply --photo1 /path/to/image")
		}
		if utilityComparePhoto2 == "" {
			printAndExit("photo2 is required, supply --photo2 /path/to/image")
		}
		c := resolveUtilityClient(true)

		var confPtr *conf.Conf
		if cmd.Flag("conf").Changed {
			parsedConf, err := resolveConf(utilityCompareConfValue)
			ifErrorExit(err)
			confPtr = &parsedConf
		}

		req, err := utility.NewCompareRequest(utilityComparePhoto1, utilityComparePhoto2, confPtr)
		ifErrorExit(err)
		req.LivenessPhoto1 = utilityCompareLivenessOne
		req.LivenessPhoto2 = utilityCompareLivenessTwo

		resp, err := c.Utility().Compare(req)
		ifErrorExit(err)
		writeOutput(resp)
	},
}

func init() {
	utilityAsmCmd.Flags().StringVar(&utilityAsmPhotoPath, "photo", "", "path to photo")
	utilityLivenessCmd.Flags().StringVar(&utilityLivenessPhoto1Path, "photo1", "", "path to first photo")
	utilityLivenessCmd.Flags().StringVar(&utilityLivenessPhoto2Path, "photo2", "", "path to second photo")

	utilityCompareCmd.Flags().StringVar(&utilityComparePhoto1, "photo1", "", "path to first photo")
	utilityCompareCmd.Flags().StringVar(&utilityComparePhoto2, "photo2", "", "path to second photo")
	utilityCompareCmd.Flags().StringVar(&utilityCompareConfValue, "conf", "", "optional confidence threshold (name or integer value)")
	utilityCompareCmd.Flags().BoolVar(&utilityCompareLivenessOne, "liveness-photo1", false, "mark first photo as liveness frame")
	utilityCompareCmd.Flags().BoolVar(&utilityCompareLivenessTwo, "liveness-photo2", false, "mark second photo as liveness frame")

	utilityCmd.AddCommand(utilityHealthCmd, utilityMetricsCmd, utilityAsmCmd, utilityLivenessCmd, utilityCompareCmd)
	rootCmd.AddCommand(utilityCmd)
}

func resolveUtilityClient(requireAccess bool) *client.Client {
	if requireAccess {
		access := strings.TrimSpace(os.Getenv("SERP_ACCESS_TOKEN"))
		if access == "" {
			printAndExit("SERP_ACCESS_TOKEN is required for this command")
		}
		return client.NewClientWithToken(access)
	}
	c, err := client.NewClient()
	ifErrorExit(err)
	return c
}
