package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var blockChain []Block
const difficulty int = 4 // Changes how long it takes to mine one new block

type Block struct {
	index int
	hash string
	prevHash string
	data string
	timeStamp string
	nonce int
}

// Generating a digital fingerprint of a block
func calculateHash (block Block) string{
	h := sha256.New()
	unique := block.data + block.prevHash + block.timeStamp + string(block.nonce)
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
		nonce: 0,
	}
	newBlock.hash = calculateHash(newBlock)

	return newBlock
}

// Proof of work algorithm: keep "mining" until we find a 
// nonce that makes our calculate hash satisfy the number
// of leading zeroes required by "difficulty".
func mineNewBlock (block *Block) {
	proofPrefix := strings.Repeat("0", difficulty)
	for calculateHash(*block)[:difficulty] != proofPrefix {
		block.nonce ++
	}

	block.hash = calculateHash(*block)
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

	// Only once our new block has proof-of-work through mining do we add it to the blockchain
	mineNewBlock(&newBlock)
	blockChain = append(blockChain, newBlock)
}

// Checks if anyone has tried tampering with the blockchain
func isBlockChainValid() bool {
	proofPrefix := strings.Repeat("0", difficulty)

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
			if (currBlock.hash[:difficulty] != proofPrefix) {
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