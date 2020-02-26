package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
)

//Encrypt encrypts message with their public key and creates a signature with your private key
func Encrypt(msg []byte, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) ([]byte, []byte) {
	label := []byte("")
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		pubKey,
		msg,
		label)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := msg
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(
		rand.Reader,
		privKey,
		newhash,
		hashed,
		&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return ciphertext, signature
}

//Decrypt decrypts message with your private key and checks the signature with their public key
func Decrypt(ciphertext []byte, signature []byte, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) string {
	var msg string
	label := []byte("")
	hash := sha256.New()

	plainText, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		privKey,
		ciphertext,
		label)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	newhash := crypto.SHA256
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	pssh := newhash.New()
	pssh.Write([]byte(plainText))
	hashed := pssh.Sum(nil)
	err = rsa.VerifyPSS(
		pubKey,
		newhash,
		hashed,
		signature,
		&opts)

	if err != nil {
		fmt.Println("Who are U? Verify Signature failed")
		os.Exit(1)
	} else {
		fmt.Println("Verify Signature successful")
	}

	return msg
}
