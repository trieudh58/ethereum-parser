package types

type Transaction struct {
	TransactionIndex string `json:"transactionIndex"`
	Hash             string `json:"hash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
}
