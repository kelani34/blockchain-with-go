package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	//	"github.com/btcsuite/btcutil/base58"
	//"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockChainAddress string
}

func NewWallet() *Wallet {
	wallet := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	wallet.privateKey = privateKey
	wallet.publicKey = &wallet.privateKey.PublicKey
	h2 := sha256.New()
	h2.Write(wallet.publicKey.X.Bytes())
	h2.Write(wallet.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	checkSum := digest6[:4]

	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], checkSum[:])
	address := base58.Encode(dc8)

	wallet.blockChainAddress = address

	return wallet
}

func (wallet *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", wallet.privateKey.D.Bytes())
}

func (wallet *Wallet) PublicKey() *ecdsa.PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}

func (wallet *Wallet) BlockChainAddress() string {
	return wallet.blockChainAddress
}

type Transaction struct {
	senderPrivateKey          *ecdsa.PrivateKey
	senderPublicKey           *ecdsa.PublicKey
	senderBlockChainAddress   string
	receiverBlockChainAddress string
	value                     float32
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, receiver string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, receiver, value}
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (transaction *Transaction) GenerateSignature() *Signature {
	m, _ := json.Marshal(transaction)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.senderPrivateKey, h[:])
	return &Signature{r, s}

}

func (transaction *Transaction) MarshallJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"senderBlockChainAddress"`
		Receiver string  `json:"receiverBlockChainAddress"`
		Value    float32 `json:"value"`
	}{
		Sender:   transaction.senderBlockChainAddress,
		Receiver: transaction.receiverBlockChainAddress,
		Value:    transaction.value,
	})
}

func (signature *Signature) String() string {
	return fmt.Sprintf("%x%x", signature.R, signature.S)
}
