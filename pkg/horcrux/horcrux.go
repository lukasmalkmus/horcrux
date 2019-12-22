package horcrux

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/hashicorp/vault/shamir"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"
)

const (
	saltLen = 32
	n       = 2 << 14
	r       = 8
	p       = 1
)

// Fragment is an encrypted fragment of the secret associated with a security
// question.
type Fragment struct {
	// K is the number of fragments required to recover the secret.
	K int

	// Question is the security question.
	Question string
	// Nonce is the random nonce used for encryption.
	Nonce []byte
	// Salt is the random salt used for scrypt.
	Salt []byte
	// Value is the encrypted share.
	Value []byte
}

// Answer is an encrypted fragment of the secret, plus the answer to the
// security question.
type Answer struct {
	// Fragment is the previously-encrypted fragment.
	Fragment
	// Answer is the answer to the security question.
	Answer string
}

// Split splits the given secret into encrypted fragments based on the given
// security questions. k is the number of fragments required to recover the
// secret.
func Split(secret []byte, questions map[string]string, k int) ([]Fragment, error) {
	shares, err := shamir.Split(secret, len(questions), k)
	if err != nil {
		return nil, err
	}

	var (
		i         = 0
		fragments = make([]Fragment, len(shares))
	)
	for question, answer := range questions {
		fragment := Fragment{
			K:        k,
			Question: question,
		}

		fragment.Salt = make([]byte, saltLen)
		if _, err := io.ReadFull(rand.Reader, fragment.Salt); err != nil {
			return nil, err
		}

		key, err := scrypt.Key([]byte(answer), fragment.Salt, n, r, p, chacha20poly1305.KeySize)
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
		if answer.K > len(answers) {
			return nil, fmt.Errorf("need at least %d answers but only have %d", answer.K, len(answers))
		}

		key, err := scrypt.Key([]byte(answer.Answer), answer.Salt, n, r, p, chacha20poly1305.KeySize)
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
