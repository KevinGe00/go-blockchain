package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
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
	unique := block.data + block.prevHash + block.timeStamp
	h.Write([]byte(unique))
	
	return hex.EncodeToString(h.Sum(nil))
}

func addNewBlock (data string) {
	var oldBlockHash string
	if len(blockChain) > 0 {
		oldBlock := blockChain[len(blockChain) - 1]
		oldBlockHash = oldBlock.hash
	} else {
		oldBlockHash = ""
	}
	
	newBlock := Block{
		prevHash: oldBlockHash,
		data: data,
		timeStamp: time.Now().String(),
	}
	newBlock.hash = calculateHash(newBlock)

	blockChain = append(blockChain, newBlock)
}

func main () {
	fmt.Println(blockChain)
	addNewBlock("Genesis block")
	fmt.Println(blockChain)

	addNewBlock("2nd block")
	fmt.Println(blockChain)
}