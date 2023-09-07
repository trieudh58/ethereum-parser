package ethereum_client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBlockNumber(t *testing.T) {
	// Create a mock HTTP server for testing.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a JSON-RPC response for eth_blockNumber.
		responseJSON := `{"jsonrpc":"2.0","id":1,"result":"0x10"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client := NewEthereumClient(server.URL)
	ctx := context.Background()

	blockNumber, err := client.GetBlockNumber(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if blockNumber != 16 {
		t.Errorf("Expected block number 16, but got %d", blockNumber)
	}
}

func TestGetTransactionsInBlock(t *testing.T) {
	// Create a mock HTTP server for testing.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a JSON-RPC response for eth_getBlockByNumber.
		responseJSON := `{"jsonrpc":"2.0","id":1,"result":{"transactions":[{"hash":"0x123","from":"0x456","to":"0x789","value":"0x100"}]}}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	client := NewEthereumClient(server.URL)
	ctx := context.Background()

	blockNumber := 123
	transactions, err := client.GetTransactionsInBlock(ctx, blockNumber)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("Expected 1 transaction, but got %d", len(transactions))
	}

	expectedTransaction := EthTransaction{
		Hash:  "0x123",
		From:  "0x456",
		To:    "0x789",
		Value: "0x100",
	}

	if transactions[0] != expectedTransaction {
		t.Errorf("Expected transaction %+v, but got %+v", expectedTransaction, transactions[0])
	}
}
