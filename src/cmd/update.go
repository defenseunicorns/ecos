package cmd

import (
	"fmt"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd := updateCmd.PersistentFlags()
	updateCmd.StringToStringVar(&ecosConfig.PackageVariables.VariableMap, "set", nil, lang.CmdUpdateFlagSet)
}

var updateCmd = &cobra.Command{
	Use:     "update [flags] STATE",
	Aliases: []string{"u"},
	Args:    cobra.ExactArgs(1),
	Short:   lang.CmdApplyShort,
	Long:    lang.CmdApplyLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from update")
	},
}
