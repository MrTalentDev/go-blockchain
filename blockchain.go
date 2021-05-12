package main

import (
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Blockchain struct {
	Chain               []Block
	CurrentTransactions []Transaction
	Nodes               map[string]bool
}

func NewBlockchain() Blockchain {
	blockchain := Blockchain{Nodes: map[string]bool{}}

	blockchain.NewBlock(100, "1")
	return blockchain
}

func (b *Blockchain) NewBlock(proof int64, previousHash string) Block {
	date := time.Now()

	var hash string

	if len(previousHash) > 0 {
		hash = previousHash
	} else {
		hash = Hash(b.LastBlock())
	}

	block := Block{
		Index:        int64(len(b.Chain) + 1),
		Timestamp:    date.Unix(),
		Transactions: b.CurrentTransactions,
		Proof:        proof,
		PreviousHash: hash,
	}

	b.Chain = append(b.Chain, block)

	return block
}

func (b *Blockchain) NewTransaction(sender string, recipient string, amount float64) int {
	transaction := Transaction{Sender: sender, Recipient: recipient, Amount: amount}
	b.CurrentTransactions = append(b.CurrentTransactions, transaction)

	return len(b.Chain)
}

func Hash(block *Block) string {
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(fmt.Sprintf("%v", block)))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (b *Blockchain) LastBlock() *Block {
	lastIndex := len(b.Chain)
	return &b.Chain[lastIndex-1]
}

func (_ *Blockchain) ProofOfWork(lastProof int64) int64 {
	var proof int64 = 0

	for !ValidProof(lastProof, proof) {
		proof++
	}

	return proof
}

func (b *Blockchain) RegisterNodes(address string) {
	u, _ := url.ParseRequestURI(address)
	hostPort := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	b.Nodes[hostPort] = true
}

func (b *Blockchain) ResolveConflicts() bool {
	max_length := int64(len(b.Chain))
	new_chain := []Block{}

	for key := range b.Nodes {
		url := fmt.Sprintf("http://%s/chain", key)

		resp, err := http.Get(url)
		if err != nil {
			panic("Error trying get url")
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic("Error parsing json to content")
		}
		var content ChainResponse

		json.Unmarshal([]byte(body), &content)

		if content.Length > max_length && ValidChain(content.Chain) {
			max_length = content.Length
			new_chain = content.Chain
		}
	}

	if len(new_chain) > 0 {
		b.Chain = new_chain
		return true
	}

	return false
}


func ValidChain(chain []Block) bool {
	last_block := &chain[0]
	current_index := 1

	for current_index < len(chain) {
		block := chain[current_index]
		fmt.Printf("[last block] {%v}", last_block)
		fmt.Printf("[current block] {%v}", block)

		if block.PreviousHash != Hash(last_block) {
			return false
		}

		if !ValidProof(last_block.Proof, block.Proof) {
			return false
		}

		last_block = &block
		current_index++
	}

	return true
}

func ValidProof(lastProof int64, proof int64) bool {
	hasher := crypto.SHA256.New()

	guess := fmt.Sprintf("%d%d", lastProof, proof)

	hasher.Write([]byte(fmt.Sprintf("%v", guess)))

	output := fmt.Sprintf("%x", hasher.Sum(nil))

	return output[0:4] == "0000"
}