package main

import (
	"encoding/gob"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/lukasmalkmus/horcrux/pkg/horcrux"
)

var outputFile string

// restoreCmd represents the restore command.
var restoreCmd = &cobra.Command{
	Use:   "restore [files]",
	Short: "Restore a file from the given horcruxes",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := make([]horcrux.Answer, len(args))
		for k, arg := range args {
			fragment, err := getFragementFromDisk(arg)
			if err != nil {
				return err
			}
			answers[k] = horcrux.Answer{Fragment: fragment}
		}

		for k, answer := range answers {
			prompt := promptui.Prompt{
				Label:    answer.Question,
				Validate: validateString,
			}
			res, err := prompt.Run()
			if err != nil {
				return handlePromptError(err)
			}
			answers[k].Answer = res
		}

		res, err := horcrux.Recover(answers)
		if err != nil {
			return err
		}

		if outputFile == "" {
			cmd.Printf("%s\n", res)
			return nil
		}

		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = f.Write(res); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	restoreCmd.Flags().StringVarP(&outputFile, "output", "o", "", "file to write restored content to")
	rootCmd.AddCommand(restoreCmd)
}

func getFragementFromDisk(fileName string) (horcrux.Fragment, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return horcrux.Fragment{}, err
	}
	defer f.Close()

	var fragment horcrux.Fragment
	if err := gob.NewDecoder(f).Decode(&fragment); err != nil {
		return horcrux.Fragment{}, err
	}
	return fragment, nil
}
