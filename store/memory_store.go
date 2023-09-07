package store

import (
	"strings"
	"sync"

	"github.com/trieudh58/ethereum-parser/types"
)

type MemoryStore struct {
	lastParsedBlock int
	subs            map[string]struct{}
	txs             []types.Transaction

	bMu sync.RWMutex // for lastParsedBlock
	sMu sync.RWMutex // for subs
	tMu sync.RWMutex // for txs
}

// AddSubscription implements Store.
func (s *MemoryStore) AddSubscription(address string) {
	s.sMu.Lock()
	defer s.sMu.Unlock()
	s.subs[strings.ToLower(address)] = struct{}{}
}

// GetLastParsedBlock implements Store.
func (s *MemoryStore) GetLastParsedBlock() int {
	s.bMu.RLock()
	defer s.bMu.RUnlock()
	return s.lastParsedBlock
}

// GetTransactions implements Store.
func (s *MemoryStore) GetTransactions(address string) []types.Transaction {
	s.tMu.RLock()
	defer s.tMu.RUnlock()
	lAddress := strings.ToLower(address)
	txs := make([]types.Transaction, 0)
	for _, tx := range s.txs {
		if tx.To == lAddress || tx.From == lAddress {
			txs = append(txs, tx)
		}
	}
	return txs
}

// IsSubscribing implements Store.
func (s *MemoryStore) IsSubscribing(address string) bool {
	s.sMu.RLock()
	defer s.sMu.RUnlock()
	_, ok := s.subs[strings.ToLower(address)]
	return ok
}

// SaveLastParsedBlock implements Store.
func (s *MemoryStore) SaveLastParsedBlock(blockNumber int) {
	s.bMu.Lock()
	defer s.bMu.Unlock()
	s.lastParsedBlock = blockNumber
}

// SaveTransactions implements Store.
func (s *MemoryStore) SaveTransactions(txs []types.Transaction) {
	s.tMu.Lock()
	defer s.tMu.Unlock()
	s.txs = append(s.txs, txs...)
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		lastParsedBlock: 0,
		subs:            make(map[string]struct{}),
		txs:             make([]types.Transaction, 0),
	}
}

var _ Store = (*MemoryStore)(nil)
