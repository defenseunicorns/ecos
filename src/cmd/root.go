package cmd

import (
	"fmt"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/mikevanhemert/ecos/src/types"
	"github.com/spf13/cobra"
)

var (
	logLevel   string
	ecosConfig = types.EcosConfig{}
)

var rootCmd = &cobra.Command{
	Use:   "ecos COMMAND",
	Short: lang.RootCmdShort,
	Long:  lang.RootCmdLong,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from ecos")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

/*
	Commands, flags, all that jazz

	ecos collect
		terraform providers mirror [DIR]
		terraform get

	ecos apply [--set KEY=value]
		terraform init -plugin-dir=[DIR] -get=false -get-plugins=false
		terraform apply

	ecos update STATE_FILE [--set KEY=value]
		terraform init -plugin-dir=[DIR] -get=false -get-plugins=false
*/
