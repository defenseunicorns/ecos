package cmd

import (
	"fmt"
	"os"

	"github.com/mikevanhemert/ecos/src/config/lang"
	"github.com/mikevanhemert/ecos/src/pkg/message"
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

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", lang.RootCmdFlagLogLevel)

	match := map[string]message.LogLevel{
		"warn":  message.WarnLevel,
		"info":  message.InfoLevel,
		"debug": message.DebugLevel,
		"trace": message.TraceLevel,
	}

	// No log level set, so use the default
	if logLevel != "" {
		if lvl, ok := match[logLevel]; ok {
			message.SetLogLevel(lvl)
			message.Debug("Log level set to " + logLevel)
		} else {
			message.Warn(lang.RootCmdErrInvalidLogLevel)
		}
	}

	// Disable progress bars for CI envs
	if os.Getenv("CI") == "true" {
		message.Debug("CI environment detected, disabling progress bars")
		message.NoProgress = true
	}
}
