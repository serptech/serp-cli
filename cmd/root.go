package cmd

import (
	"fmt"
	"os"

	cliutils "github.com/serptech/serp-cli/utils"
	serpUtils "github.com/serptech/serp-go/utils"
	"github.com/spf13/cobra"
)

var (
	outputPath      string
	limit           int
	offset          int
	flagAccessToken string
	flagRootToken   string
	debug           bool
	baseURL         string
)

var rootCmd = &cobra.Command{
	Use:     "serptech",
	Short:   "SERP is a real-time facial recognition platform.",
	Version: cliutils.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flagAccessToken != "" {
			ifErrorExit(os.Setenv("SERP_ACCESS_TOKEN", flagAccessToken))
		}

		if flagRootToken != "" {
			ifErrorExit(os.Setenv("SERP_ROOT_TOKEN", flagRootToken))
		}

		if os.Getenv("SERP_ACCESS_TOKEN") == "" {
			if root := os.Getenv("SERP_ROOT_TOKEN"); root != "" {
				ifErrorExit(os.Setenv("SERP_ACCESS_TOKEN", root))
			}
		}
		if baseURL != "" {
			ifErrorExit(os.Setenv("SERP_BASE_URL", baseURL))
		}
		if debug {
			ifErrorExit(os.Setenv("SERP_DEBUG", fmt.Sprintf("%v", debug)))
			serpUtils.Warn().Msgf("%v", os.Environ())
		}
	},
	Long: `
	
         # ##########                                                              
      # ## ############ #                                                          
   # ## ## ############ ##                                                         
  ## ## ##           ## ## ##   ########### ############ ############ ############ 
  ## ## ##           ## ## ##   ##          ##           ##        ## ##        ## 
  ## ## ##                      ##          ##           ##        ## ##        ## 
  ## ## ##           ## ## ##   ##          ######################### ##        ## 
  ## ## ##           ## ## ##   ##          ##           ##           ##        ## 
    ################ ## ##      ######################## ##           ##        ## 
       ############# ## #                                                       ## 
         ########### ##                                                             

`,
}

func ifErrorExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printAndExit(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func Execute() {
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug cli and client")
	rootCmd.PersistentFlags().StringVar(&flagAccessToken, "token", "", "serptech.ru access token (SERP_ACCESS_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&flagRootToken, "root-token", "", "root API token (SERP_ROOT_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "serptech.ru API base URL override")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "path to file for writing output result")
	rootCmd.PersistentFlags().IntVar(&limit, "limit", 20, "the number of output items, maximum 1000 entries per request")
	rootCmd.PersistentFlags().IntVar(&offset, "offset", 0, "a sequential number of an output item, to return a sampling after this one")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
