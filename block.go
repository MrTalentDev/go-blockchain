package main

type Block struct {
	Index int64
	Timestamp int64
	Transactions []Transaction
	Proof int64
	PreviousHash string
}