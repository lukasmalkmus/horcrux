package main

import (
	"encoding/gob"
	"errors"
	"os"
	"strings"

	"github.com/golang/snappy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/lukasmalkmus/horcrux/horcrux"
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

func writeFragmentToDisk(fileName string, fragment *horcrux.Fragment) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	w := snappy.NewBufferedWriter(f)
	defer w.Close()

	if err = gob.NewEncoder(w).Encode(fragment); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}
	return f.Sync()
}

func getFragementFromDisk(fileName string) (*horcrux.Fragment, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := snappy.NewReader(f)

	var fragment horcrux.Fragment
	if err := gob.NewDecoder(r).Decode(&fragment); err != nil {
		return nil, err
	}
	return &fragment, nil
}
