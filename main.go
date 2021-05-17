package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

var blockChain []Block
const difficulty int = 4 // Changes how long it takes to mine one new block

type Block struct {
	Index int
	Hash string
	PrevHash string
	Data string
	TimeStamp string
	Nonce int
}

// Generating a digital fingerprint of a block
func calculateHash (block Block) string{
	h := sha256.New()
	unique := block.Data + block.PrevHash + block.TimeStamp + string(block.Nonce)
	h.Write([]byte(unique))
	
	return hex.EncodeToString(h.Sum(nil))
}

func getLastBlock () Block{
	return blockChain[len(blockChain) - 1]
}

func createNewBlock (index int, prevHash, data string) Block{
	newBlock := Block{
		Index: index,
		PrevHash: prevHash,
		Data: data,
		TimeStamp: time.Now().String(),
		Nonce: 0,
	}
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

// Proof of work algorithm: keep "mining" until we find a 
// nonce that makes our calculate hash satisfy the number
// of leading zeroes required by "difficulty".
func mineNewBlock (block *Block) {
	proofPrefix := strings.Repeat("0", difficulty)
	for calculateHash(*block)[:difficulty] != proofPrefix {
		block.Nonce ++
	}

	block.Hash = calculateHash(*block)
}

func addNewBlock (data string) {
	var lastBlockHash string
	var index int
	if len(blockChain) > 0 {
		lastBlock := getLastBlock()
		lastBlockHash = lastBlock.Hash
		index = lastBlock.Index + 1
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
			if (currBlock.Index != prevBlock.Index + 1) {
				return false
			}
			if (currBlock.PrevHash != prevBlock.Hash) {
				return false
			}
			if (calculateHash(currBlock) != currBlock.Hash) {
				return false
			}
			if (currBlock.Hash[:difficulty] != proofPrefix) {
				return false
			}
		} 
    }
	return true
}


// Blockchain as an API
func handleRequests() {
	r := mux.NewRouter()

	// Paths
    r.HandleFunc("/blockchain", getBlockchainHandler)
	r.HandleFunc("/blockchain/{index}", getBlockHandler)

	log.Fatal(http.ListenAndServe(":10000", r))
}

// View entire blockchain
func getBlockchainHandler(w http.ResponseWriter, r *http.Request){
    json.NewEncoder(w).Encode(blockChain)
    fmt.Println("Endpoint Hit: /blockchain")
}

// View single block on blockchain
func getBlockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    i, _ := strconv.Atoi(vars["index"])
	if (i < 0 || len(blockChain) <= i) {
		fmt.Fprintf(w, "Index out of range")
		return
	} 
	json.NewEncoder(w).Encode(blockChain[i])
}

func main () {
	addNewBlock("Genesis block")
	addNewBlock("2nd block")
	addNewBlock("3rd block")
	spew.Dump(blockChain)
	
	spew.Dump(isBlockChainValid())	// True
	blockChain[1].Data = "I've messed with the blockchain"
	spew.Dump(isBlockChainValid()) // False

	handleRequests()
}