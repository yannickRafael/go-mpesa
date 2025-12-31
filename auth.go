package mpesa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// GenerateBearerToken generates a Bearer Token by encrypting the API Key with the Public Key.
// 1. Formats the public key into a PEM block.
// 2. Encrypts the API Key using RSA PKCS#1 v1.5.
// 3. Base64 encodes the result.
func GenerateBearerToken(apiKey, publicKeyStr string) (string, error) {
	if apiKey == "" || publicKeyStr == "" {
		return "", errors.New("missing API Key or Public Key")
	}

	// formatting certificate string as the Node.js lib does
	certificate := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", publicKeyStr)

	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("key is not an RSA public key")
	}

	// Encrypt API key using public key (PKCS1v15)
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt API key: %v", err)
	}

	// Return formatted string, Bearer token in base64 format
	return "Bearer " + base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
