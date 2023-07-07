package cmd

import (
	"fmt"
	"os"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/mikevanhemert/ecos/src/pkg/archive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd := updateCmd.PersistentFlags()
	updateCmd.StringToStringVar(&ecosConfig.PackageVariables.VariableMap, "set", nil, lang.CmdUpdateFlagSet)
}

var updateCmd = &cobra.Command{
	Use:     "update [flags] ORIGINAL_ARCHIVE UPDATED_ARCHIVE",
	Aliases: []string{"u"},
	Args:    cobra.ExactArgs(2),
	Short:   lang.CmdApplyShort,
	Long:    lang.CmdApplyLong,
	Run: func(cmd *cobra.Command, args []string) {
		originalArchive := args[0]
		updatedArchive := args[1]

		update := archive.NewOrDieUpdate(originalArchive, updatedArchive)
		defer update.ClearTempPaths()

		if err := update.Update(); err != nil {
			fmt.Printf("Failed to update Terraform archive %s with values in %s: %s\n", originalArchive, updatedArchive, err)
			os.Exit(1)
		}

		fmt.Print("\nComplete\n\n")
	},
}
