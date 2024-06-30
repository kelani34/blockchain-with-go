package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	timestamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

const MINING_DIFFICULTY = 9

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	//b := new(Block)
	//b.timestamp = time.Now().UnixNano()
	//return b

	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		previousHash: previousHash,
		transactions: transactions,
	}
}

func (b *Block) Print() {
	fmt.Printf("timestamp: %d\n", b.timestamp)
	fmt.Printf("nonce: %d\n", b.nonce)
	fmt.Printf("previousHash: %x\n", b.previousHash)
	for _, tx := range b.transactions {
		tx.Print()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previousHash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}

type BlockChain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockChain() *BlockChain {
	b := &Block{}
	bc := new(BlockChain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) AddTransaction(sender string, receiver string, value float32) {
	tx := NewTransaction(sender, receiver, value)
	bc.transactionPool = append(bc.transactionPool, tx)
}

func (bc *BlockChain) CopyTransaction() []*Transaction {
	transactions := make([]*Transaction, len(bc.transactionPool))
	for _, tx := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(tx.senderBlockChainAddress, tx.receiverBlockChainAddress, tx.value))
	}
	return transactions
}

func (bc *BlockChain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

type Transaction struct {
	senderBlockChainAddress   string
	receiverBlockChainAddress string
	value                     float32
}

func NewTransaction(sender string, receiver string, value float32) *Transaction {
	return &Transaction{sender, receiver, value}
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	fmt.Printf("senderBlockChainAddress          %s\n", t.senderBlockChainAddress)
	fmt.Printf("receiverBlockChainAddress        %s\n", t.receiverBlockChainAddress)
	fmt.Printf("value                            %.2f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockChainAddress   string  `json:"senderBlockChainAddress"`
		ReceiverBlockChainAddress string  `json:"receiverBlockChainAddress"`
		Value                     float32 `json:"value"`
	}{
		SenderBlockChainAddress:   t.senderBlockChainAddress,
		ReceiverBlockChainAddress: t.receiverBlockChainAddress,
		Value:                     t.value,
	})
}

func (bc *BlockChain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s \n\n\n ", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		fmt.Printf("%s \n", strings.Repeat("*", 25))
		block.Print()
	}
}

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	//block := &Block{nonce: 1}
	//
	//fmt.Println(block.Hash())

	blockChain := NewBlockChain()
	blockChain.Print()

	blockChain.AddTransaction("A", "B", 1.0)
	previousHash := blockChain.LastBlock().Hash()
	blockChain.CreateBlock(5, previousHash)
	blockChain.Print()

	blockChain.AddTransaction("C", "D", 899.0)
	blockChain.AddTransaction("X", "Y", 9999.0)
	previousHash = blockChain.LastBlock().Hash()
	blockChain.CreateBlock(4, previousHash)
	blockChain.Print()
}
