// Package horcrux is a security question based secret sharing utility.
//
// A secret is split into multiple fragments and every fragment is associated
// with a security question. A key derived from the answer to that question is
// used to encrypt the fragment using ChaCha20Poly1305. Only a given number of
// fragments is needed to fully restore the original secret.
//
// This package has not been audited by cryptography or security professionals.
package horcrux
