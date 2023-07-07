package cmd

import (
	"fmt"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyFlags := destroyCmd.PersistentFlags()
	destroyFlags.StringToStringVar(&ecosConfig.PackageVariables.VariableMap, "set", nil, lang.CmdApplyFlagSet)
}

var destroyCmd = &cobra.Command{
	Use:     "destroy [flags] ARCHIVE",
	Aliases: []string{"d"},
	Args:    cobra.ExactArgs(1),
	Short:   lang.CmdDestroyShort,
	Long:    lang.CmdDestroyLong,
	Run: func(cmd *cobra.Command, args []string) {
		archiveName := args[0]
		fmt.Println(archiveName)
	},
}
