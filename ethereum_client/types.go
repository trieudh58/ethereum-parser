package ethereum_client

// Request structure for the JSON-RPC call
type EthRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// Response structure for the JSON-RPC call
type EthGenericResponse struct {
	Result interface{} `json:"result"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

type EthBlockNumberResponseResponse struct {
	Result string `json:"result"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

type EthGetBlockByNumberResponseResponse struct {
	Result struct {
		Number       string           `json:"number"`
		Hash         string           `json:"hash"`
		Transactions []EthTransaction `json:"transactions"`
		// ignore others, see: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbyhash
	} `json:"result"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	ID int `json:"id"`
}

type EthTransaction struct {
	TransactionIndex string `json:"transactionIndex"`
	Hash             string `json:"hash"`
	BlockNumber      string
	From             string
	To               string
	Value            string
	// ignore others, see: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionbyhash
}
