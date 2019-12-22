// Package horcrux provides security question style secret recovery while
// preserving end-to-end cryptographic security.
//
// Given N pairs of security questions and answers, the secret is split using
// Shamir's Secret Sharing algorithm into N shares, one for each question. A
// 256-bit key is derived from the answer to each question using scrypt, and the
// share is then encrypted with that key using ChaCha20Poly1305.
//
// To recover the secret given K of N answers, the secret keys are re-derived
// and the shares are decrypted and combined.
//
// This package has not been audited by cryptography or security professionals.
package horcrux
