package cmd

import (
	"fmt"
	"os"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/mikevanhemert/ecos/src/pkg/archive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)

	applyFlags := applyCmd.PersistentFlags()
	applyFlags.StringToStringVar(&ecosConfig.PackageVariables, "set", nil, lang.CmdApplyFlagSet)
}

var applyCmd = &cobra.Command{
	Use:     "apply [flags] ARCHIVE",
	Aliases: []string{"a"},
	Args:    cobra.ExactArgs(1),
	Short:   lang.CmdApplyShort,
	Long:    lang.CmdApplyLong,
	Run: func(cmd *cobra.Command, args []string) {
		archiveName := args[0]

		archiver := archive.New(&ecosConfig)
		defer archiver.ClearTempPaths()

		if err := archiver.Apply(archiveName); err != nil {
			fmt.Printf("Failed to apply Terraform from archvie %s: %s\n", archiveName, err)
			os.Exit(1)
		}

		fmt.Print("\nComplete\n\n")
	},
}
