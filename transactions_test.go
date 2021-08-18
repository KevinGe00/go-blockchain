package main

import (
	"testing"
)

func TestTransactionEncryptionAndDecryption(t *testing.T){
	wallet1 := initializeNewWallet()

	// encrypt transaction using transaction data and wallet1's public key
	encryptedTransaction := encryptTransaction(wallet1.PublicKey)

	// decrypt encrypted transaction using wallet1's private key
	decryptedMsg := decryptTransaction(wallet1.PrivateKey, encryptedTransaction)

	if decryptedMsg != "transaction data" {
		t.Errorf("Expected \"transaction data\" but got: %v", decryptedMsg)
	}
}

func TestDSFS(t *testing.T){
	wallet1 := initializeNewWallet()
	wallet2 := initializeNewWallet()

	transaction := Transaction{
		Sender: wallet1.PublicKey,
		Receiver: wallet2.PublicKey,
		Amount: 10,
	}
	
	sig, hashSum := signTransaction(wallet1.PrivateKey, transaction)
	transaction.Signature = sig

	if !verifySignature(sig, transaction, hashSum) {
		t.Errorf("Verification that signature is indeed signed of the input data(msgHashSum) by the owner of the sender in the transaction failed when it shouldn't have")
	}
	
	// Alter hashSum then try verifying again
	transaction.Amount = 10000
	hashSum = generateUniqueTransactionHashSum(transaction)

	if verifySignature(sig, transaction, hashSum) {
		t.Errorf("Verification that signature is indeed signed of the input data(msgHashSum) by the owner of the sender in the transaction passed when it shouldn't have")
	}

}