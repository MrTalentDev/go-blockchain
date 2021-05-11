package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Nodes struct {
	Address []string
}

type MutexBlockchain struct {
	sync.Mutex
	Blockchain Blockchain
}

func NewMutexBlockchain() MutexBlockchain {
	return MutexBlockchain{Blockchain: NewBlockchain()}
}

var GLOBAL_BLOCKCHAIN MutexBlockchain = NewMutexBlockchain()

func MineHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GLOBAL_BLOCKCHAIN.Mutex.Lock()
	defer GLOBAL_BLOCKCHAIN.Mutex.Unlock()

	last_proof := GLOBAL_BLOCKCHAIN.Blockchain.LastBlock().Proof
	proof := GLOBAL_BLOCKCHAIN.Blockchain.ProofOfWork(last_proof)

	GLOBAL_BLOCKCHAIN.Blockchain.NewTransaction("0", "57e430de001d498fbf6e493a79665d57", 1.0)

	block := GLOBAL_BLOCKCHAIN.Blockchain.NewBlock(proof, "")

	data := struct {
		Message      string
		Index        int64
		Transactions []Transaction
		Proof        int64
		PreviousHash string
	}{
		"new block forged",
		block.Index,
		block.Transactions,
		block.Proof,
		block.PreviousHash,
	}

	json.NewEncoder(w).Encode(data)
}

func ChainHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GLOBAL_BLOCKCHAIN.Mutex.Lock()
	defer GLOBAL_BLOCKCHAIN.Mutex.Unlock()

	data := struct {
		Chain  []Block
		Length int
	}{
		GLOBAL_BLOCKCHAIN.Blockchain.Chain,
		len(GLOBAL_BLOCKCHAIN.Blockchain.Chain),
	}

	json.NewEncoder(w).Encode(data)
}

func NodesResolveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	GLOBAL_BLOCKCHAIN.Mutex.Lock()
	defer GLOBAL_BLOCKCHAIN.Mutex.Unlock()

	var message string

	if GLOBAL_BLOCKCHAIN.Blockchain.ResolveConflicts() {
		message = "Our chain was replaced"
	} else {
		message = "Our chain is authoritative"
	}

	data := struct {
		Message    string
		Chain []Block
	}{
		message,
		GLOBAL_BLOCKCHAIN.Blockchain.Chain,
	}

	json.NewEncoder(w).Encode(data)

}

func NodesRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nodes Nodes
	_ = json.NewDecoder(r.Body).Decode(&nodes)
	if len(nodes.Address) <= 0 {
		panic("send some address")
	}

	GLOBAL_BLOCKCHAIN.Mutex.Lock()
	defer GLOBAL_BLOCKCHAIN.Mutex.Unlock()

	for _, node := range nodes.Address {
		GLOBAL_BLOCKCHAIN.Blockchain.RegisterNodes(node)
	}

	data := struct {
		Message    string
		TotalNodes int
	}{
		"New nodes have been added",
		len(GLOBAL_BLOCKCHAIN.Blockchain.Chain),
	}

	json.NewEncoder(w).Encode(data)
}

func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var transaction Transaction
	_ = json.NewDecoder(r.Body).Decode(&transaction)

	GLOBAL_BLOCKCHAIN.Mutex.Lock()
	defer GLOBAL_BLOCKCHAIN.Mutex.Unlock()

	index := GLOBAL_BLOCKCHAIN.Blockchain.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)

	data := struct {
		Message string
	}{
		fmt.Sprintf("new transaction created, index %d", index),
	}

	json.NewEncoder(w).Encode(data)
}

func main() {
	rourte := mux.NewRouter()
	rourte.HandleFunc("/mine", MineHandler).Methods("GET")
	rourte.HandleFunc("/chain", ChainHandler).Methods("GET")
	rourte.HandleFunc("/nodes/resolve", NodesResolveHandler).Methods("GET")
	rourte.HandleFunc("/nodes/register", NodesRegisterHandler).Methods("POST")
	rourte.HandleFunc("/transaction/new", TransactionsHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", rourte))
}
