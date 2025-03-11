package txbuilder

import (
	"context"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/theQRL/go-zond/common"
	"github.com/theQRL/go-zond/core/types"
	"github.com/theQRL/go-zond/rpc"
	"github.com/theQRL/go-zond/zondclient"
)

type Client struct {
	Timeout time.Duration
	rpchost string
	client  *zondclient.Client
	logger  *logrus.Entry

	gasSuggestionMutex sync.Mutex
	lastGasSuggestion  time.Time
	lastGasCap         *big.Int
	lastTipCap         *big.Int

	blockHeight      uint64
	blockHeightTime  time.Time
	blockHeightMutex sync.Mutex
}

func NewClient(rpchost string) (*Client, error) {
	headers := map[string]string{}

	if strings.HasPrefix(rpchost, "headers(") {

		headersEnd := strings.Index(rpchost, ")")
		headersStr := rpchost[8:headersEnd]
		rpchost = rpchost[headersEnd+1:]

		for _, headerStr := range strings.Split(headersStr, "|") {
			headerParts := strings.Split(headerStr, ":")
			headers[strings.Trim(headerParts[0], " ")] = strings.Trim(headerParts[1], " ")
		}
	}

	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, rpchost)
	if err != nil {
		return nil, err
	}

	for hKey, hVal := range headers {
		rpcClient.SetHeader(hKey, hVal)
	}

	return &Client{
		client:  zondclient.NewClient(rpcClient),
		rpchost: rpchost,
		logger:  logrus.WithField("client", rpchost),
	}, nil
}

func (client *Client) GetName() string {
	url, _ := url.Parse(client.rpchost)
	name := strings.TrimSuffix(url.Host, ".theqrl.org")
	return name
}

func (client *Client) GetZondClient() *zondclient.Client {
	return client.client
}

func (client *Client) GetRPCHost() string {
	return client.rpchost
}

func (client *Client) UpdateWallet(wallet *Wallet) error {
	if wallet.GetChainId() == nil {
		chainId, err := client.GetChainId()
		if err != nil {
			return err
		}
		wallet.SetChainId(chainId)
	}

	nonce, err := client.GetNonceAt(wallet.GetAddress(), nil)
	if err != nil {
		return err
	}
	wallet.SetNonce(nonce)

	balance, err := client.GetBalanceAt(wallet.GetAddress())
	if err != nil {
		return err
	}
	wallet.SetBalance(balance)

	return nil
}

func (client *Client) getContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	if client.Timeout > 0 {
		return context.WithTimeout(ctx, client.Timeout)
	}
	return context.WithCancel(ctx)
}

func (client *Client) GetChainId() (*big.Int, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.ChainID(ctx)
}

func (client *Client) GetNonceAt(wallet common.Address, blockNumber *big.Int) (uint64, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.NonceAt(ctx, wallet, blockNumber)
}

func (client *Client) GetPendingNonceAt(wallet common.Address) (uint64, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.PendingNonceAt(ctx, wallet)
}

func (client *Client) GetBalanceAt(wallet common.Address) (*big.Int, error) {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.client.BalanceAt(ctx, wallet, nil)
}

func (client *Client) GetSuggestedFee() (*big.Int, *big.Int, error) {
	client.gasSuggestionMutex.Lock()
	defer client.gasSuggestionMutex.Unlock()

	if time.Since(client.lastGasSuggestion) < 12*time.Second {
		return client.lastGasCap, client.lastTipCap, nil
	}

	ctx, cancel := client.getContext()
	defer cancel()

	gasCap, err := client.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}
	tipCap, err := client.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}

	client.lastGasSuggestion = time.Now()
	client.lastGasCap = gasCap
	client.lastTipCap = tipCap
	return gasCap, tipCap, nil
}

func (client *Client) SendTransaction2(tx *types.Transaction) error {
	ctx, cancel := client.getContext()
	defer cancel()

	return client.SendTransactionCtx(ctx, tx)
}

func (client *Client) SendTransactionCtx(ctx context.Context, tx *types.Transaction) error {
	client.logger.Tracef("submitted transaction %v", tx.Hash().String())

	return client.client.SendTransaction(ctx, tx)
}

func (client *Client) GetTransactionReceiptCtx(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	client.logger.Tracef("get receipt: 0x%x", txHash.Bytes())

	return client.client.TransactionReceipt(ctx, txHash)
}

func (client *Client) GetBlockHeight() (uint64, error) {
	client.blockHeightMutex.Lock()
	defer client.blockHeightMutex.Unlock()

	if time.Since(client.blockHeightTime) < 12*time.Second {
		return client.blockHeight, nil
	}

	client.logger.Tracef("get block number")

	ctx, cancel := client.getContext()
	defer cancel()

	blockHeight, err := client.client.BlockNumber(ctx)
	if err != nil {
		return blockHeight, err
	}
	if blockHeight > client.blockHeight {
		client.blockHeight = blockHeight
		client.blockHeightTime = time.Now()
	}
	return client.blockHeight, nil
}
