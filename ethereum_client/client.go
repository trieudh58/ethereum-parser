package ethereum_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type EthereumJsonRpcClient struct {
	endpoint string
	client   *http.Client
}

func NewEthereumClient(endpoint string) EthereumJsonRpcClient {
	return EthereumJsonRpcClient{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

func (c EthereumJsonRpcClient) GetBlockNumber(ctx context.Context) (int, error) {
	// see: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_blocknumber
	request := EthRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{},
		ID:      1,
	}

	var res EthBlockNumberResponseResponse
	if err := c.makeRequest(ctx, request, &res); err != nil {
		return 0, err
	}

	if res.Error.Code != 0 {
		return 0, fmt.Errorf("JSON-RPC error with code %d: %s", res.Error.Code, res.Error.Message)
	}
	blockNumber, err := hexStringToDecimal(res.Result[2:])
	if err != nil {
		return 0, err
	}
	return int(blockNumber), nil
}

func (c EthereumJsonRpcClient) GetTransactionsInBlock(ctx context.Context, blockNumber int) ([]EthTransaction, error) {
	// see: https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbynumber
	request := EthRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			fmt.Sprintf("0x%s", baseDecimalToHexString(int64(blockNumber))),
			true,
		},
		ID: 1,
	}
	var res EthGetBlockByNumberResponseResponse
	if err := c.makeRequest(ctx, request, &res); err != nil {
		return []EthTransaction{}, err
	}

	if res.Error.Code != 0 {
		return []EthTransaction{}, fmt.Errorf("JSON-RPC error with code %d: %s", res.Error.Code, res.Error.Message)
	}

	return res.Result.Transactions, nil
}

func (c EthereumJsonRpcClient) makeRequest(ctx context.Context, req EthRequest, res interface{}) error {
	requestBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, _ := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(requestBytes))
	resp, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&res); err != nil {
		return err
	}

	return nil
}
