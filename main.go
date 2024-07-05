package main

import (
	"blockChainWithGo/wallet"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKey())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockChainAddress())
}
