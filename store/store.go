package store

import "github.com/trieudh58/ethereum-parser/types"

type Store interface {
	// block
	GetLastParsedBlock() int
	SaveLastParsedBlock(blockNumber int)

	// transaction
	GetTransactions(address string) []types.Transaction
	SaveTransactions([]types.Transaction)

	// address
	AddSubscription(address string)
	IsSubscribing(address string) bool
}
