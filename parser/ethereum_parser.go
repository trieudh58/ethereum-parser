package parser

import (
	"github.com/trieudh58/ethereum-parser/store"
	"github.com/trieudh58/ethereum-parser/types"
)

type EthereumParser struct {
	store store.Store
}

func NewEthereumParser(store store.Store) *EthereumParser {
	return &EthereumParser{
		store,
	}
}

// GetCurrentBlock implements Parser.
func (p *EthereumParser) GetCurrentBlock() int {
	blockNo := p.store.GetLastParsedBlock()
	return blockNo
}

// GetTransactions implements Parser.
func (p *EthereumParser) GetTransactions(address string) []types.Transaction {
	return p.store.GetTransactions(address)
}

// Subscribe implements Parser.
func (p *EthereumParser) Subscribe(address string) bool {
	p.store.AddSubscription(address)
	return true
}

var _ Parser = (*EthereumParser)(nil)
