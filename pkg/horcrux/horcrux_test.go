package horcrux_test

import (
	"fmt"
	"testing"

	"github.com/lukasmalkmus/horcrux/pkg/horcrux"
)

func Example() {
	secret := []byte("my favorite password")
	questions := map[string]string{
		"What's your first pet's name?":     "Spot",
		"What's your least favorite food?":  "Broccoli",
		"What's your mother's maiden name?": "Hernandez",
		"What's your real name?":            "Rumplestiltskin",
	}

	// Split into four fragments, any two of which can be combined to recover
	// the secret.
	frags, err := horcrux.Split(secret, questions, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Answer two of the security questions.
	answers := make([]horcrux.Answer, 2)
	for i := range answers {
		answers[i] = horcrux.Answer{
			Fragment: frags[i],
			Answer:   questions[frags[i].Question],
		}
	}

	// Recover the original secret.
	s, err := horcrux.Recover(answers)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(s))
	// Output:
	// my favorite password
}

var (
	secret    = []byte("my favorite password")
	questions = map[string]string{
		"What's your first pet's name?":     "Spot",
		"What's your least favorite food?":  "Broccoli",
		"What's your mother's maiden name?": "Hernandez",
		"What's your real name?":            "Rumplestiltskin",
	}
)

func TestRecoverTooFewAnswers(t *testing.T) {
	frags, err := horcrux.Split(secret, questions, 2)
	if err != nil {
		t.Fatal(err)
	}

	answers := make([]horcrux.Answer, 1)
	for i := range answers {
		answers[i] = horcrux.Answer{
			Fragment: frags[i],
			Answer:   questions[frags[i].Question],
		}
	}

	s, err := horcrux.Recover(answers)
	if s != nil {
		t.Fatalf("Expected nil, but was %v", s)
	}

	expected := "need at least 2 answers but only have 1"
	actual := err.Error()
	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}

func TestRecoverBadAnswers(t *testing.T) {
	frags, err := horcrux.Split(secret, questions, 2)
	if err != nil {
		t.Fatal(err)
	}

	answers := make([]horcrux.Answer, 2)
	for i := range answers {
		answers[i] = horcrux.Answer{
			Fragment: frags[i],
			Answer:   questions[frags[i].Question] + "woo",
		}
	}

	s, err := horcrux.Recover(answers)
	if s != nil {
		t.Fatalf("Expected nil, but was %v", s)
	}

	expected := "chacha20poly1305: message authentication failed"
	actual := err.Error()
	if actual != expected {
		t.Fatalf("Expected %v but was %v", expected, actual)
	}
}
