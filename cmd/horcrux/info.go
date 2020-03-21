package main

import (
	"text/template"

	"github.com/spf13/cobra"
)

const templateStr = `ID: {{ .ID }}
Owner: {{ .Owner }}
Question: {{ .Question }}
`

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:   "info [file]",
	Short: "Get info of a fragment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fragment, err := getFragementFromDisk(args[0])
		if err != nil {
			return err
		}

		tmpl, err := template.New("fragment").Parse(templateStr)
		if err != nil {
			return err
		}
		return tmpl.Execute(cmd.OutOrStdout(), fragment)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
