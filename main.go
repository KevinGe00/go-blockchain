package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var blockChain []Block

type Block struct {
	index int
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

func createNewBlock (index int, prevHash, data string) Block{
	newBlock := Block{
		index: index,
		prevHash: prevHash,
		data: data,
		timeStamp: time.Now().String(),
	}
	newBlock.hash = calculateHash(newBlock)

	return newBlock
}

func addNewBlock (data string) {
	var lastBlockHash string
	var index int
	if len(blockChain) > 0 {
		lastBlock := getLastBlock()
		lastBlockHash = lastBlock.hash
		index = lastBlock.index + 1
	} else {
		// Genesis block has no prev hash
		lastBlockHash = ""
		index = 0
	}

	newBlock := createNewBlock(index, lastBlockHash, data)
	blockChain = append(blockChain, newBlock)
}

// Checks if anyone has tried tampering with the blockchain
func isBlockChainValid() bool {
	for i, currBlock := range blockChain {
    	if (i > 0 ) { // Skip first block
			prevBlock := blockChain[i - 1]
			if (currBlock.index != prevBlock.index + 1) {
				return false
			}
			if (currBlock.prevHash != prevBlock.hash) {
				return false
			}
			if (calculateHash(currBlock) != currBlock.hash) {
				return false
			}
		} 
    }
	return true
}

func main () {
	addNewBlock("Genesis block")
	addNewBlock("2nd block")
	addNewBlock("3rd block")
	spew.Dump(blockChain)
	
	spew.Dump(isBlockChainValid())

	blockChain[1].data = "I've messed with the blockchain"
	blockChain[1].hash = calculateHash(blockChain[1])

	spew.Dump(isBlockChainValid())
}