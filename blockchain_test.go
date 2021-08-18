package main

import (
	"testing"
)

func TestAddNewBlock(t *testing.T){
	addNewBlock("Genesis block")
	addNewBlock("2nd block")
	addNewBlock("3rd block")
	
	if len(blockChain) != 3 {
		t.Errorf("Only added 3 blocks to block chain got %v blocks", len(blockChain))
	}
}

func TestIsBlockChainValid(t *testing.T){
	// Block chain has 3 exisiting blocks already from TestAddNewBlock

	if !isBlockChainValid() {
        t.Errorf("Block chain hasn't been tampered with, expected true upon validation but got false")
    }

	blockChain[1].Data = "I've messed with the blockchain"

    if isBlockChainValid() {
        t.Errorf("Block chain has been tampered with, expected false upon validation but got true")
    }
}