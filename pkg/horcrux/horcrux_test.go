package horcrux_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lukasmalkmus/horcrux/pkg/horcrux"
)

func Example() {
	secret := []byte("my favorite password")
	questions := []horcrux.Question{
		{"What's your first pet's name?", "Spot"},
		{"What's your least favorite food?", "Broccoli"},
		{"What's your mother's maiden name?", "Hernandez"},
		{"What's your real name?", "Rumplestiltskin"},
	}

	// Split into four fragments, any two of which can be combined to recover
	// the secret.
	fragments, err := horcrux.Split(secret, questions, 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Answer two of the security questions.
	answers := make([]horcrux.Answer, 2)
	for i := range answers {
		answers[i] = horcrux.Answer{
			Fragment: fragments[i],
			Answer:   questions[i].Answer,
		}
	}

	// Recover the original secret.
	s, err := horcrux.Recover(answers)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", s)
	// Output:
	// my favorite password
}

var (
	secret    = []byte("my favorite password")
	questions = []horcrux.Question{
		{"What's your first pet's name?", "Spot"},
		{"What's your least favorite food?", "Broccoli"},
		{"What's your mother's maiden name?", "Hernandez"},
		{"What's your real name?", "Rumplestiltskin"},
	}
)

func TestRecoverTooFewAnswers(t *testing.T) {
	fragments, err := horcrux.Split(secret, questions, 2)
	ok(t, err)

	answers := make([]horcrux.Answer, 1)
	for i := range answers {
		answers[i] = horcrux.Answer{
			Fragment: fragments[i],
			Answer:   questions[i].Answer,
		}
	}

	s, err := horcrux.Recover(answers)
	assert(t, s == nil, "Expected nil, but was %v", s)

	expected := "need at least 2 answers but only have 1"
	actual := err.Error()
	equals(t, actual, expected)
}

func TestRecoverBadAnswers(t *testing.T) {
	fragments, err := horcrux.Split(secret, questions, 2)
	ok(t, err)

	answers := make([]horcrux.Answer, 2)
	for i := range answers {
		answers[i] = horcrux.Answer{
			Fragment: fragments[i],
			Answer:   questions[i].Answer + "woo",
		}
	}

	s, err := horcrux.Recover(answers)
	assert(t, s == nil, "Expected nil, but was %v", s)

	expected := "chacha20poly1305: message authentication failed"
	actual := err.Error()
	equals(t, actual, expected)
}

func benchmarkSplit(b *testing.B, size int) {
	secret := make([]byte, size)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := horcrux.Split(secret, questions, 2)
		ok(b, err)
	}
}

func BenchmarkSplit64KB(b *testing.B)  { benchmarkSplit(b, 1024*64) }
func BenchmarkSplit1MB(b *testing.B)   { benchmarkSplit(b, 1024*1024) }
func BenchmarkSplit128MB(b *testing.B) { benchmarkSplit(b, 1024*1024*128) }
func BenchmarkSplit1GB(b *testing.B)   { benchmarkSplit(b, 1024*1024*1024) }

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()
	if !condition {
		tb.Fatalf("\033[31m "+msg+"\033[39m\n\n", v...)
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31m unexpected error: %s\033[39m\n\n", err.Error())
	}
}

// equals fails the test if got is not equal to want.
func equals(tb testing.TB, got, want interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(got, want) {
		tb.Fatalf("\033[31m\n\n\tgot: %#v\n\n\twant: %#v\033[39m\n\n", got, want)
	}
}
