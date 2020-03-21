package main

import (
	"errors"
	"os"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "horcrux",
	Short: "A security question based secret sharing utility",
	Long: `Horcrux is a security question based secret sharing
utility. Yes, the name is inspired by Harry Potter.

> Documentation & Support: https://github.com/lukasmalkmus/horcrux
> Source & Copyright Information: https://github.com/lukasmalkmus/horcrux`,
}

func main() {
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handlePromptError(err error) error {
	if err == nil || err == promptui.ErrAbort {
		return nil
	} else if cause := errors.Unwrap(err); cause != nil {
		err = cause
	}
	return err
}

func validateString(input string) error {
	if input = strings.TrimSpace(input); input == "" {
		return errors.New("cannot be empty")
	}
	return nil
}
