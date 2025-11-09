package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type HMACAlgorithm string

const (
	AlgorithmSHA1   HMACAlgorithm = "sha1"
	AlgorithmSHA256 HMACAlgorithm = "sha256"
)

// ComputeHMAC computes HMAC digest for given data and secret
func ComputeHMAC(algorithm HMACAlgorithm, secret, data string) string {
	var h hash.Hash

	switch algorithm {
	case AlgorithmSHA1:
		h = hmac.New(sha1.New, []byte(secret))
	case AlgorithmSHA256:
		h = hmac.New(sha256.New, []byte(secret))
	default:
		h = hmac.New(sha1.New, []byte(secret))
	}

	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// ValidateHMAC validates if the provided digest matches the computed HMAC
func ValidateHMAC(algorithm HMACAlgorithm, secret, data, providedDigest string) bool {
	expectedDigest := ComputeHMAC(algorithm, secret, data)
	return hmac.Equal([]byte(expectedDigest), []byte(providedDigest))
}
