package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"strconv"
)

type Wallet struct {
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type Transaction struct {
	Sender, Receiver *rsa.PublicKey
	Signature []byte
	Amount int
	Inputs []TransactionInput 	// Proof that send has funds to send
	Outputs []TransactionOutput //
}

type TransactionInput struct {
	TransactionOutputId string
	UTXO TransactionOutput // Unspent Transaction Output
}

// These are referenced when new transactions are made as proof the sender has coins
type TransactionOutput struct {
	Id string
	Receiver *rsa.PublicKey; 
	Amout int;
	parentTransactionId string;
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
		[]byte("transaction data"),
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

// Generate the data to be signed by a private key
func generateUniqueTransactionHashSum(transaction Transaction) []byte{
	h := sha256.New()
	unique := strconv.Itoa(transaction.Sender.E) + strconv.Itoa(transaction.Receiver.E) + strconv.Itoa(transaction.Amount)
	_, err := h.Write([]byte(unique))
	if err != nil {
		panic(err)
	}

	hashSum := h.Sum(nil)
	return hashSum
}

// Sign unique data(msgHashSum) using a private key
// Returns the signature and data that was signed
func signTransaction(privateKey *rsa.PrivateKey, transaction Transaction) ([]byte, []byte){
	msgHashSum := generateUniqueTransactionHashSum(transaction)

	signature, err := rsa.SignPSS(
		rand.Reader, 
		privateKey, // private key of sender
		crypto.SHA256, 
		msgHashSum, 
		nil,
	)

	if err != nil {
		panic(err)
	}

	return signature, msgHashSum
}

// Verify that the signature is indeed signed of the input data(msgHashSum) by the owner of the sender in the transaction
func verifySignature(signature []byte, transaction Transaction, msgHashSum []byte) bool{
	err := rsa.VerifyPSS(
		transaction.Sender, 
		crypto.SHA256, 
		msgHashSum, 
		signature,
		nil,
	)
	if err != nil {
		fmt.Println("Could not verify signature: ", err)
		return false
	}

	// Signature is valid
	return true
}