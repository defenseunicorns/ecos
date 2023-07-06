package cmd

import (
	"fmt"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/spf13/cobra"
)

var packageCmd = &cobra.Command{
	Use:     "package",
	Aliases: []string{"p"},
	Short:   lang.CmdPackageShort,
}

var packageCreateCmd = &cobra.Command{
	Use:     "create [ DIRECTORY ]",
	Aliases: []string{"c"},
	Args:    cobra.MaximumNArgs(1),
	Short:   lang.CmdPackageCreateShort,
	Long:    lang.CmdPackageCreateLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from create")
	},
}

var packageApplyCmd = &cobra.Command{
	Use:     "apply [ PACKAGE ]",
	Aliases: []string{"a"},
	Args:    cobra.MaximumNArgs(1),
	Short:   lang.CmdPackageApplyShort,
	Long:    lang.CmdPackageApplyLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from apply")
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.AddCommand(packageCreateCmd)
	packageCmd.AddCommand(packageApplyCmd)

	bindApplyFlags()
}

func bindApplyFlags() {
	applyFlags := packageCmd.PersistentFlags()

	applyFlags.StringToStringVar(&ecosConfig.PackageVariables.VariableMap, "set", nil, lang.CmdPackageApplyFlagSet)
}
