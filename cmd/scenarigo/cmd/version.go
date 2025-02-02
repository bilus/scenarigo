package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/bilus/scenarigo/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("print %s version", appName),
	Long:  fmt.Sprintf("Prints %s version.", appName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s version %s\n", appName, version.String())
	},
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "%s version %s\n", appName, version.String())
}
