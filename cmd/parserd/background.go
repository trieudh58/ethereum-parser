package main

import (
	"context"
	"log"

	"github.com/trieudh58/ethereum-parser/ethereum_client"
	"github.com/trieudh58/ethereum-parser/store"
	"github.com/trieudh58/ethereum-parser/types"
)

type Worker struct {
	id     int
	store  store.Store
	client *ethereum_client.EthereumJsonRpcClient
}

func NewWorker(id int, store store.Store, client *ethereum_client.EthereumJsonRpcClient) *Worker {
	return &Worker{
		id:     id,
		store:  store,
		client: client,
	}
}

func (w *Worker) Work(ctx context.Context, blockNo int) error {
	log.Printf("start reading block %d from worker %d\n", blockNo, w.id)
	txs, err := w.client.GetTransactionsInBlock(ctx, blockNo)
	if err != nil {
		return err
	}

	savingTxs := make([]types.Transaction, 0)
	for _, tx := range txs {
		if w.store.IsSubscribing(tx.From) || w.store.IsSubscribing(tx.To) {
			savingTxs = append(savingTxs, types.Transaction{
				TransactionIndex: tx.TransactionIndex,
				Hash:             tx.Hash,
				BlockNumber:      tx.BlockNumber,
				From:             tx.From,
				To:               tx.To,
				Value:            tx.Value,
			})
		}
	}
	if len(savingTxs) > 0 {
		w.store.SaveTransactions(savingTxs)
		log.Printf("%d transaction(s) stored.\n", len(savingTxs))
	}

	if blockNo > w.store.GetLastParsedBlock() {
		w.store.SaveLastParsedBlock(blockNo)
	}
	return nil
}
