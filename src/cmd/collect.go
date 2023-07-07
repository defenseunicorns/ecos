package cmd

import (
	"fmt"
	"os"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/mikevanhemert/ecos/src/pkg/collect"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(collectCmd)

	// collectFlags := collectCmd.PersistentFlags()
	// TODO
}

var collectCmd = &cobra.Command{
	Use:     "collect [flags]",
	Aliases: []string{"c"},
	Short:   lang.CmdCollectShort,
	Long:    lang.CmdCollectLong,
	Run: func(cmd *cobra.Command, args []string) {
		collect := collect.NewOrDie(&ecosConfig)
		//defer collect.ClearTempPaths()

		if err := collect.Collect(); err != nil {
			fmt.Printf("Failed to collect Terraform resources: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\nComplete\n\n")
	},
}
