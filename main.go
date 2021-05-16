package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var blockChain []Block

type Block struct {
	hash string
	prevHash string
	data string
	timeStamp string
}

// Generating a digital fingerprint of a block
func calculateHash (block Block) string{
	h := sha256.New()
	unique := block.data + block.prevHash + block.timeStamp
	h.Write([]byte(unique))
	
	return hex.EncodeToString(h.Sum(nil))
}

func getLastBlock () Block{
	return blockChain[len(blockChain) - 1]
}

func addNewBlock (data string) {
	var oldBlockHash string
	if len(blockChain) > 0 {
		oldBlock := getLastBlock()
		oldBlockHash = oldBlock.hash
	} else {
		// Genesis block has no prev hash
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
	addNewBlock("Genesis block")
	addNewBlock("2nd block")
	spew.Dump(blockChain)
}