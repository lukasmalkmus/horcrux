package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lukasmalkmus/horcrux/pkg/horcrux"
)

// Question extends horcrux.Question.
type Question struct {
	horcrux.Question

	Owner string
}

// Questions are multiple questions.
type Questions []Question

// HorcruxQuestions returns a slice of *horcrux.Question.
func (q Questions) HorcruxQuestions() []horcrux.Question {
	questions := make([]horcrux.Question, len(q))
	for k, question := range q {
		questions[k] = question.Question
	}
	return questions
}

var questions Questions

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create [file]",
	Short: "Create new horcruxes from a given file",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var lenHorcruxes int
		prompt := promptui.Prompt{
			Label: "Amount of horcruxes to split the file into",
			Validate: func(input string) error {
				if v, err := strconv.Atoi(input); err != nil {
					return errors.New("invalid number")
				} else if v < 2 || v > 255 {
					return fmt.Errorf("must be a number between 2 and 255")
				}
				return nil
			},
		}
		res, err := prompt.Run()
		if err != nil {
			return handlePromptError(err)
		}
		lenHorcruxes, _ = strconv.Atoi(res)

		var threshold int
		prompt = promptui.Prompt{
			Label: "Amount of horcruxes needed to restore the file",
			Validate: func(input string) error {
				if v, err := strconv.Atoi(input); err != nil {
					return errors.New("invalid number")
				} else if v < 2 || v > 255 {
					return errors.New("must be a number between 2 and 255")
				} else if v > lenHorcruxes {
					return fmt.Errorf("must be smaller or equal to the amount of horcruxes to create (%d)", lenHorcruxes)
				}
				return nil
			},
			Default: strconv.Itoa(lenHorcruxes),
		}
		if res, err = prompt.Run(); err != nil {
			return handlePromptError(err)
		}
		threshold, _ = strconv.Atoi(res)
		viper.Set("threshold", threshold)

		questions = make([]Question, lenHorcruxes)
		for k := range questions {
			cmd.Printf("Creating horcrux %d of %d. Please enter some information to protect your horcrux:\n", k+1, lenHorcruxes)

			prompt = promptui.Prompt{
				Label:   "Your name",
				Default: defaultUsernameIfNotTaken(),
				Validate: func(input string) error {
					if err := validateString(input); err != nil {
						return err
					} else if hasOwner(questions, input) {
						return fmt.Errorf("%q is already taken", input)
					}
					return nil
				},
			}
			if questions[k].Owner, err = prompt.Run(); err != nil {
				return handlePromptError(err)
			}

			prompt = promptui.Prompt{
				Label:    "Your security question",
				Validate: validateString,
			}
			if questions[k].Question.Question, err = prompt.Run(); err != nil {
				return handlePromptError(err)
			}

			prompt = promptui.Prompt{
				Label:    "Answer to the security question",
				Validate: validateString,
			}
			if questions[k].Answer, err = prompt.Run(); err != nil {
				return handlePromptError(err)
			}
		}

		prompt = promptui.Prompt{
			Label:     fmt.Sprintf("Create %d horcruxes with %d needed to restore the file", lenHorcruxes, threshold),
			IsConfirm: true,
		}
		if _, err = prompt.Run(); err != nil {
			return handlePromptError(err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		secret, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		fragments, err := horcrux.Split(secret, questions.HorcruxQuestions(), viper.GetInt("threshold"))
		if err != nil {
			return err
		}

		for k, fragment := range fragments {
			owner := questions[k].Owner
			owner = strings.ToLower(owner)
			owner = strings.ReplaceAll(owner, " ", "_")
			filename := fmt.Sprintf("%s.horcrux", owner)
			if err := writeFragmentToDisk(filename, fragment); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func writeFragmentToDisk(fileName string, fragment horcrux.Fragment) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(fragment); err != nil {
		return err
	}
	return nil
}

func hasOwner(questions []Question, owner string) bool {
	for _, question := range questions {
		if question.Owner == owner {
			return true
		}
	}
	return false
}

func defaultUsernameIfNotTaken() string {
	var username string
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}
	return username
}