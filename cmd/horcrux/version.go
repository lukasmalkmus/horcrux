package main

import (
	"github.com/spf13/cobra"

	"github.com/axiomhq/pkg/version"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and build details",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version.Print("horcrux"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
