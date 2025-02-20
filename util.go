package goverify

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// loadPrivateKey loads an parses a PEM encoded private key file.
func LoadPublicKeyFromFile(path string) (Verifier, error) {
	dat, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return parsePublicKey(dat)
}

func LoadPublicKeyFromString(key string) (Verifier, error) {

	return parsePublicKey([]byte(key))
}

// parsePublicKey parses a PEM encoded private key.
func parsePublicKey(pemBytes []byte) (Verifier, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
        case "RSA PUBLIC KEY":
		rsa, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}

	return newVerifierFromKey(rawkey)
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func LoadPrivateKeyFromFile(path string) (Signer, error) {
	dat, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return parsePrivateKey(dat)
}

func LoadPrivateKeyFromString(key string) (Signer, error) {
	return parsePrivateKey([]byte(key))
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, err := pem.Decode(pemBytes)
	fmt.Println("block=", block, "err=", err)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}
	fmt.Println("block.Type=", block.Type)

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	case "PRIVATE KEY":
		privkey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = privkey
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}
