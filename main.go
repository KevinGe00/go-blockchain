package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var blockChain []Block

type Block struct {
	hash string
	prevHash string
	data string
	timeStamp string
}

func calculateHash (block Block) string{
	// Generating a digital fingerprint of input block
	h := sha256.New()
	unique := block.data + block.hash + block.prevHash + block.timeStamp
	h.Write([]byte(unique))
	
	return hex.EncodeToString(h.Sum(nil))
}

func main () {
	b := Block{
		hash: "111",
		prevHash: "",
		data: "genesis",
		timeStamp: "11pm",
	}
	fmt.Println(calculateHash(b))
	
}