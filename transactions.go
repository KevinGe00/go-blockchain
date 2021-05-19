package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type Wallet struct {
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func initializeNewWallet() Wallet{
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	newWallet := Wallet{
		PublicKey: &publicKey,
		PrivateKey: privateKey,
	}

	return newWallet
}

func encryptTransaction(publicKey *rsa.PublicKey) []byte{
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte("secret message"),
	nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("encrypted bytes: ", encryptedBytes)
	return encryptedBytes
}

func decryptTransaction(privateKey *rsa.PrivateKey, cipherText []byte) string{
	decryptedBytes, err := privateKey.Decrypt(
		nil, 
		cipherText, 
		&rsa.OAEPOptions{Hash: crypto.SHA256}, // Same hashing we used to encrypt the message
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("decrypted message: ", string(decryptedBytes))
	return string(decryptedBytes)
}