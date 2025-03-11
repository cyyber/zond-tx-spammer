package txbuilder

import (
	"context"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-zond/accounts/abi/bind"
	"github.com/theQRL/go-zond/common"
	"github.com/theQRL/go-zond/core/types"
)

type Wallet struct {
	nonceMutex     sync.Mutex
	balanceMutex   sync.RWMutex
	dilithiumKey   *dilithium.Dilithium
	address        common.Address
	chainid        *big.Int
	pendingNonce   atomic.Uint64
	confirmedNonce uint64
	balance        *big.Int

	txNonceChans     map[uint64]*nonceStatus
	txNonceMutex     sync.Mutex
	lastConfirmation uint64
}

type nonceStatus struct {
	receipt *types.Receipt
	channel chan bool
}

func NewWallet(seed string) (*Wallet, error) {
	wallet := &Wallet{
		txNonceChans: map[uint64]*nonceStatus{},
	}
	err := wallet.loadKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (wallet *Wallet) loadKeyFromSeed(seed string) error {
	var dilithiumKey *dilithium.Dilithium
	if seed == "" {
		var err error
		dilithiumKey, err = dilithium.New()
		if err != nil {
			return err
		}
	} else {
		var err error
		dilithiumKey, err = dilithium.NewDilithiumFromHexSeed(seed)
		if err != nil {
			return err
		}
	}

	wallet.dilithiumKey = dilithiumKey
	wallet.address = dilithiumKey.GetAddress()
	return nil
}

func (wallet *Wallet) GetAddress() common.Address {
	return wallet.address
}

func (wallet *Wallet) GetDilithiumKey() *dilithium.Dilithium {
	return wallet.dilithiumKey
}

func (wallet *Wallet) GetChainId() *big.Int {
	return wallet.chainid
}

func (wallet *Wallet) GetNonce() uint64 {
	return wallet.pendingNonce.Load()
}

func (wallet *Wallet) GetBalance() *big.Int {
	wallet.balanceMutex.RLock()
	defer wallet.balanceMutex.RUnlock()
	return wallet.balance
}

func (wallet *Wallet) SetChainId(chainid *big.Int) {
	wallet.chainid = chainid
}

func (wallet *Wallet) SetNonce(nonce uint64) {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	pendingNonce := wallet.pendingNonce.Load()
	if nonce > pendingNonce {
		wallet.pendingNonce.Store(nonce)
	}

	wallet.confirmedNonce = nonce
}

func (wallet *Wallet) GetNextNonce() uint64 {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()
	return wallet.pendingNonce.Add(1) - 1
}

func (wallet *Wallet) SetBalance(balance *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = balance
}

func (wallet *Wallet) SubBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Sub(wallet.balance, amount)
}

func (wallet *Wallet) AddBalance(amount *big.Int) {
	wallet.balanceMutex.Lock()
	defer wallet.balanceMutex.Unlock()
	wallet.balance = wallet.balance.Add(wallet.balance, amount)
}

func (wallet *Wallet) BuildDynamicFeeTx(txData *types.DynamicFeeTx) (*types.Transaction, error) {
	wallet.nonceMutex.Lock()
	txData.ChainID = wallet.chainid
	txData.Nonce = wallet.pendingNonce.Add(1) - 1
	wallet.nonceMutex.Unlock()
	return wallet.signTx(txData)
}

func (wallet *Wallet) BuildBoundTx(txData *TxMetadata, buildFn func(transactOpts *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	transactor, err := bind.NewKeyedTransactorWithChainID(wallet.dilithiumKey, wallet.chainid)
	if err != nil {
		return nil, err
	}

	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	transactor.Context = context.Background()
	transactor.From = wallet.address
	nonce := wallet.pendingNonce.Add(1) - 1
	transactor.Nonce = big.NewInt(0).SetUint64(nonce)

	transactor.GasTipCap = txData.GasTipCap.ToBig()
	transactor.GasFeeCap = txData.GasFeeCap.ToBig()
	transactor.GasLimit = txData.Gas
	transactor.Value = txData.Value.ToBig()
	transactor.NoSend = true

	tx, err := buildFn(transactor)
	if err != nil {
		wallet.pendingNonce.Store(nonce)
		return nil, err
	}

	return tx, nil
}

func (wallet *Wallet) ReplaceDynamicFeeTx(txData *types.DynamicFeeTx, nonce uint64) (*types.Transaction, error) {
	txData.ChainID = wallet.chainid
	txData.Nonce = nonce
	return wallet.signTx(txData)
}

func (wallet *Wallet) ResetPendingNonce(client *Client) {
	wallet.nonceMutex.Lock()
	defer wallet.nonceMutex.Unlock()

	nonce, err := client.GetPendingNonceAt(wallet.address)
	if nonce < wallet.confirmedNonce {
		nonce = wallet.confirmedNonce
	}

	if err == nil && wallet.pendingNonce.Load() != nonce {
		logrus.Warnf("Resyncing pending nonce for %v from %d to %d", wallet.address.String(), wallet.pendingNonce.Load(), nonce)
		wallet.pendingNonce.Store(nonce)
	}
}

func (wallet *Wallet) signTx(txData types.TxData) (*types.Transaction, error) {
	tx := types.NewTx(txData)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(wallet.chainid), wallet.dilithiumKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (wallet *Wallet) getTxNonceChan(targetNonce uint64) (*nonceStatus, bool) {
	wallet.txNonceMutex.Lock()
	defer wallet.txNonceMutex.Unlock()

	if wallet.confirmedNonce > targetNonce {
		return nil, false
	}

	nonceChan := wallet.txNonceChans[targetNonce]
	if nonceChan != nil {
		return nonceChan, false
	}

	nonceChan = &nonceStatus{
		channel: make(chan bool),
	}
	wallet.txNonceChans[targetNonce] = nonceChan

	return nonceChan, len(wallet.txNonceChans) == 1
}
