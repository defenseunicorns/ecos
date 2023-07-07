package cmd

import (
	"fmt"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)

	applyFlags := applyCmd.PersistentFlags()
	applyFlags.StringToStringVar(&ecosConfig.PackageVariables.VariableMap, "set", nil, lang.CmdApplyFlagSet)
}

var applyCmd = &cobra.Command{
	Use:     "apply [flags] STATE",
	Aliases: []string{"a"},
	Args:    cobra.ExactArgs(1),
	Short:   lang.CmdApplyShort,
	Long:    lang.CmdApplyLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from apply")
	},
}
