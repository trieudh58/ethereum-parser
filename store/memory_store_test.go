package store

import (
	"testing"

	"github.com/trieudh58/ethereum-parser/types"
)

func TestMemoryStore_AddSubscription(t *testing.T) {
	store := NewMemoryStore()
	address := "0x123abc"

	store.AddSubscription(address)

	if !store.IsSubscribing(address) {
		t.Errorf("Expected subscription for address %s, but it was not found", address)
	}
}

func TestMemoryStore_GetLastParsedBlock(t *testing.T) {
	store := NewMemoryStore()
	store.SaveLastParsedBlock(100)

	lastParsedBlock := store.GetLastParsedBlock()

	if lastParsedBlock != 100 {
		t.Errorf("Expected last parsed block 100, but got %d", lastParsedBlock)
	}
}

func TestMemoryStore_GetTransactions(t *testing.T) {
	store := NewMemoryStore()
	address := "0x123abc"
	transactions := []types.Transaction{
		{From: address, To: "0x456def"},
		{From: "0x789ghi", To: address},
		{From: "0xaaa", To: "0xbbb"},
	}

	store.SaveTransactions(transactions)

	retrievedTransactions := store.GetTransactions(address)

	if len(retrievedTransactions) != 2 {
		t.Errorf("Expected 2 transactions for address %s, but got %d", address, len(retrievedTransactions))
	}
}

func TestMemoryStore_IsSubscribing(t *testing.T) {
	store := NewMemoryStore()
	address := "0x123abc"

	if store.IsSubscribing(address) {
		t.Errorf("Expected no subscription for address %s, but it was found", address)
	}

	store.AddSubscription(address)

	if !store.IsSubscribing(address) {
		t.Errorf("Expected subscription for address %s, but it was not found", address)
	}
}

func TestMemoryStore_SaveLastParsedBlock(t *testing.T) {
	store := NewMemoryStore()

	store.SaveLastParsedBlock(200)

	lastParsedBlock := store.GetLastParsedBlock()

	if lastParsedBlock != 200 {
		t.Errorf("Expected last parsed block 200, but got %d", lastParsedBlock)
	}
}

func TestMemoryStore_SaveTransactions(t *testing.T) {
	store := NewMemoryStore()
	transactions := []types.Transaction{
		{From: "0x123abc", To: "0x456def"},
		{From: "0x789ghi", To: "0x123abc"},
	}

	store.SaveTransactions(transactions)

	retrievedTransactions := store.GetTransactions("0x123abc")

	if len(retrievedTransactions) != 2 {
		t.Errorf("Expected 2 transactions for address 0x123abc, but got %d", len(retrievedTransactions))
	}
}
