package cmd

import (
	"fmt"
	"os"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/mikevanhemert/ecos/src/pkg/archive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyFlags := destroyCmd.PersistentFlags()
	destroyFlags.StringToStringVar(&ecosConfig.PackageVariables, "set", nil, lang.CmdApplyFlagSet)
}

var destroyCmd = &cobra.Command{
	Use:     "destroy [flags] ARCHIVE",
	Aliases: []string{"d"},
	Args:    cobra.ExactArgs(1),
	Short:   lang.CmdDestroyShort,
	Long:    lang.CmdDestroyLong,
	Run: func(cmd *cobra.Command, args []string) {
		archiveName := args[0]

		archiver := archive.New(&ecosConfig)
		defer archiver.ClearTempPaths()

		if err := archiver.Destroy(archiveName); err != nil {
			fmt.Printf("Failed to destroy Terraform in archive %s: %s\n", archiveName, err)
			os.Exit(1)
		}

		fmt.Print("\nComplete\n\n")
	},
}
