package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

var blockChain []Block
const difficulty int = 2 // Changes how long it takes to mine one new block

type Block struct {
	Index int
	Hash string
	PrevHash string
	Data string
	TimeStamp string
	Nonce int
}

// Used for unmarshalling POST request body
type Data_ struct {
	Data string
}

// Generating a digital fingerprint of a block
func calculateHash (block Block) string{
	h := sha256.New()
	unique := block.Data + block.PrevHash + block.TimeStamp + strconv.Itoa(block.Nonce)
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

func addNewBlock (data string) Block{
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

	return newBlock
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
				// Work has not been done yet
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
	r.HandleFunc("/mine", mineBlockHandler).Methods("POST")
	r.HandleFunc("/{index}", getBlockHandler)
    r.HandleFunc("/", getBlockchainHandler)

	// Run the server
	port := ":10000"
	fmt.Println("\nListening on port " + port[1:])
	log.Fatal(http.ListenAndServe(port, r))
}

// View entire blockchain
func getBlockchainHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(blockChain)
}

// View single block on blockchain
func getBlockHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    i, _ := strconv.Atoi(vars["index"])
	if (i < 0 || len(blockChain) <= i) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("HTTP 400: Bad Request - Index out of range"))
		return
	} 

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blockChain[i])
}

func mineBlockHandler (w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var postData Data_
	json.Unmarshal(reqBody, &postData)

	newBlock := addNewBlock(postData.Data) // Mine for new block

	response, err := json.MarshalIndent(newBlock, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("HTTP 400: Bad Request"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
	fmt.Println("New block mined and added to block chain")
}

func main () {
	wallet1 := initializeNewWallet()
	wallet2 := initializeNewWallet()

	cipherText := encryptTransaction(wallet1.PublicKey)
	decryptTransaction(wallet1.PrivateKey, cipherText)
	
	transaction := Transaction{
		Sender: wallet1.PublicKey,
		Receiver: wallet2.PublicKey,
		Amount: 10,
	}

	hashSum := generateUniqueTransactionHashSum(transaction)
	sig := generateSignature(wallet1.PrivateKey, hashSum)
	transaction.Signature = sig

	fmt.Println(verifySignature(sig, transaction, hashSum))
	
	// Alter hashSum then try verifying again
	transaction.Amount = 10000
	hashSum = generateUniqueTransactionHashSum(transaction)
	fmt.Println(verifySignature(sig, transaction, hashSum))

	// Testing blockchain functions
	addNewBlock("Genesis block")
	addNewBlock("2nd block")
	addNewBlock("3rd block")
	spew.Dump(blockChain)
	
	spew.Dump(isBlockChainValid())	// True
	blockChain[1].Data = "I've messed with the blockchain"
	spew.Dump(isBlockChainValid()) // False

	handleRequests()
}