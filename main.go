package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Block struct {
	Index     int
	Timestamp string
	Hash      string
	PrevHash  string
}

type Blockchain []Block

// Creates new block based on last block in the chain
func genBlock(bc Blockchain) Block {
	t := time.Now()
	lastBlock := bc[len(bc)-1]
	var newBlock Block
	newBlock.Index = lastBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = lastBlock.Hash
	newBlock.Hash = genHash(newBlock)
	return newBlock
}

// Adds new block to the chain if it is valid
func (bc Blockchain) addBlock(newBlock Block) (Blockchain, error) {
	lastBlock := bc[len(bc)-1]
	if blockValid(lastBlock, newBlock) {
		bc = append(bc, newBlock)
		return bc, nil
	}
	return bc, errors.New("New block is invalid.")
}

// Initializes a new blockchain with a "genesis block" unless it's non-empty
func (bc Blockchain) init() Blockchain {
	// Return the chain if it already has blocks
	if len(bc) > 0 {
		return bc
	}
	t := time.Now()

	var genesisBlock Block
	genesisBlock.Index = 0
	genesisBlock.Timestamp = t.String()
	genesisBlock.PrevHash = ""
	genesisBlock.Hash = genHash(genesisBlock)

	bc = append(bc, genesisBlock)
	return bc
}

// Creates a SHA256 hash from a block's Index, Timestamp, and Previous Hash
func genHash(block Block) string {
	record := string(block.Index) + string(block.Timestamp) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Checks the validity of the newest block being added
// compared to the last block in the chain
func blockValid(lastBlock, newBlock Block) bool {
	if newBlock.Index != lastBlock.Index+1 {
		return false
	}

	if newBlock.PrevHash != lastBlock.Hash {
		return false
	}

	if genHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func main() {
	var bc Blockchain
	var err error

	bc = bc.init()
	for i := 0; i < 4; i++ {
		newBlock := genBlock(bc)
		bc, err = bc.addBlock(newBlock)
		if err != nil {
			fmt.Println(err)
		}
	}
	spew.Dump(bc)
}
