package horcrux

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"

	"github.com/lukasmalkmus/horcrux/pkg/shamir"
)

const (
	saltLen = 32
	scryptN = 2 << 14
	scryptR = 8
	scryptP = 1
)

// Question is a securtiy question composed of the actual question and the
// answer.
type Question struct {
	// Question is the actual question.
	Question string
	// Answer is the answer to that question.
	Answer string
}

// Fragment is an encrypted fragment of the secret associated with a security
// question.
type Fragment struct {
	// Threshold is the number of fragments required to recover the secret.
	// Must be between 2 and 255.
	Threshold int
	// Question is the security question.
	Question string
	// Nonce is the random nonce used for encryption.
	Nonce []byte
	// Salt is the random salt used for scrypt.
	Salt []byte
	// Value is the encrypted share.
	Value []byte
}

// Answer is an encrypted fragment of the secret and the answer to the security
// question.
type Answer struct {
	// Fragment is the previously encrypted fragment.
	Fragment
	// Answer is the answer to the security question.
	Answer string
}

// Split the given secret into as much encrypted fragments as security questions
// given. The threshold specifies how many fragments are needed to recover the
// secret and ranges from 2 to 255.
func Split(secret []byte, questions []Question, threshold int) ([]Fragment, error) {
	shares, err := shamir.Split(secret, len(questions), threshold)
	if err != nil {
		return nil, err
	}

	var (
		i         = 0
		fragments = make([]Fragment, len(shares))
	)
	for _, question := range questions {
		fragment := Fragment{
			Threshold: threshold,
			Question:  question.Question,
		}

		fragment.Salt = make([]byte, saltLen)
		if _, err := io.ReadFull(rand.Reader, fragment.Salt); err != nil {
			return nil, err
		}

		key, err := scrypt.Key([]byte(question.Answer), fragment.Salt,
			scryptN, scryptR, scryptP, chacha20poly1305.KeySize)
		if err != nil {
			return nil, err
		}

		aead, err := chacha20poly1305.New(key)
		if err != nil {
			return nil, err
		}

		fragment.Nonce = make([]byte, aead.NonceSize())
		if _, err = io.ReadFull(rand.Reader, fragment.Nonce); err != nil {
			return nil, err
		}

		fragment.Value = aead.Seal(nil, fragment.Nonce, shares[i], nil)

		fragments[i] = fragment

		i++
	}

	return fragments, nil
}

// Recover combines the given answers and returns the original secret.
func Recover(answers []Answer) ([]byte, error) {
	shares := make([][]byte, 0)

	for _, answer := range answers {
		if answer.Threshold > len(answers) {
			return nil, fmt.Errorf("need at least %d answers but only have %d",
				answer.Threshold, len(answers))
		}

		key, err := scrypt.Key([]byte(answer.Answer), answer.Salt,
			scryptN, scryptR, scryptP, chacha20poly1305.KeySize)
		if err != nil {
			return nil, err
		}

		aead, err := chacha20poly1305.New(key)
		if err != nil {
			return nil, err
		}

		share, err := aead.Open(nil, answer.Nonce, answer.Value, nil)
		if err != nil {
			return nil, err
		}

		shares = append(shares, share)
	}

	return shamir.Combine(shares)
}
