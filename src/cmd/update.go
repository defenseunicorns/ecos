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
		oldArchive := args[0]
		archiveName := args[1]

		archiver := archive.New(&ecosConfig)
		defer archiver.ClearTempPaths()

		if err := archiver.Update(archiveName, oldArchive); err != nil {
			fmt.Printf("Failed to update Terraform archive %s with values in %s: %s\n", archiveName, oldArchive, err)
			os.Exit(1)
		}

		fmt.Print("\nComplete\n\n")
	},
}
