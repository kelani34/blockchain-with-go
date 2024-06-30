package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timestamp    int64 ``
	transactions []string
}

func NewBlock(nonce int, previousHash string) *Block {
	//b := new(Block)
	//b.timestamp = time.Now().UnixNano()
	//return b

	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
	}
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	b := NewBlock(0, "init hash")
	fmt.Println(b)
}
