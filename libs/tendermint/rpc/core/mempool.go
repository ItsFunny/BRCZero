package core

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"time"

	"github.com/pkg/errors"

	"github.com/brc20-collab/brczero/libs/cosmos-sdk/baseapp"
	abci "github.com/brc20-collab/brczero/libs/tendermint/abci/types"
	"github.com/brc20-collab/brczero/libs/tendermint/config"
	mempl "github.com/brc20-collab/brczero/libs/tendermint/mempool"
	ctypes "github.com/brc20-collab/brczero/libs/tendermint/rpc/core/types"
	rpctypes "github.com/brc20-collab/brczero/libs/tendermint/rpc/jsonrpc/types"
	"github.com/brc20-collab/brczero/libs/tendermint/types"
)

//-----------------------------------------------------------------------------
// NOTE: tx should be signed, but this is only checked at the app level (not by Tendermint!)

// BroadcastTxAsync returns right away, with no response. Does not wait for
// CheckTx nor DeliverTx results.
// More: https://docs.tendermint.com/master/rpc/#/Tx/broadcast_tx_async
func BroadcastTxAsync(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	rtx := mempl.GetRealTxFromWrapCMTx(tx)
	err := env.Mempool.CheckTx(tx, nil, mempl.TxInfo{})

	if err != nil {
		return nil, err
	}
	return &ctypes.ResultBroadcastTx{Hash: rtx.Hash()}, nil
}

// BroadcastTxSync returns with the response from CheckTx. Does not wait for
// DeliverTx result.
// More: https://docs.tendermint.com/master/rpc/#/Tx/broadcast_tx_sync
func BroadcastTxSync(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	resCh := make(chan *abci.Response, 1)
	rtx := mempl.GetRealTxFromWrapCMTx(tx)
	err := env.Mempool.CheckTx(tx, func(res *abci.Response) {
		resCh <- res
	}, mempl.TxInfo{})
	if err != nil {
		return nil, err
	}
	res := <-resCh
	r := res.GetCheckTx()
	// reset r.Data for compatibility with cosmwasmJS
	r.Data = nil
	return &ctypes.ResultBroadcastTx{
		Code:      r.Code,
		Data:      r.Data,
		Log:       r.Log,
		Codespace: r.Codespace,
		Hash:      rtx.Hash(),
	}, nil
}

// BroadcastTxCommit returns with the responses from CheckTx and DeliverTx.
// More: https://docs.tendermint.com/master/rpc/#/Tx/broadcast_tx_commit
func BroadcastTxCommit(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	subscriber := ctx.RemoteAddr()

	if env.EventBus.NumClients() >= config.DynamicConfig.GetMaxSubscriptionClients() {
		return nil, fmt.Errorf("max_subscription_clients %d reached", config.DynamicConfig.GetMaxSubscriptionClients())
	} else if env.EventBus.NumClientSubscriptions(subscriber) >= env.Config.MaxSubscriptionsPerClient {
		return nil, fmt.Errorf("max_subscriptions_per_client %d reached", env.Config.MaxSubscriptionsPerClient)
	}

	// Subscribe to tx being committed in block.
	subCtx, cancel := context.WithTimeout(ctx.Context(), SubscribeTimeout)
	defer cancel()
	rtx := mempl.GetRealTxFromWrapCMTx(tx)
	q := types.EventQueryTxFor(rtx)
	deliverTxSub, err := env.EventBus.Subscribe(subCtx, subscriber, q)
	if err != nil {
		err = fmt.Errorf("failed to subscribe to tx: %w", err)
		env.Logger.Error("Error on broadcast_tx_commit", "err", err)
		return nil, err
	}
	defer env.EventBus.Unsubscribe(context.Background(), subscriber, q)

	// Broadcast tx and wait for CheckTx result
	checkTxResCh := make(chan *abci.Response, 1)
	err = env.Mempool.CheckTx(tx, func(res *abci.Response) {
		checkTxResCh <- res
	}, mempl.TxInfo{})
	if err != nil {
		env.Logger.Error("Error on broadcastTxCommit", "err", err)
		return nil, fmt.Errorf("error on broadcastTxCommit: %v", err)
	}
	checkTxResMsg := <-checkTxResCh
	checkTxRes := checkTxResMsg.GetCheckTx()
	if checkTxRes.Code != abci.CodeTypeOK {
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxRes,
			DeliverTx: abci.ResponseDeliverTx{},
			Hash:      rtx.Hash(),
		}, nil
	}

	// Wait for the tx to be included in a block or timeout.
	select {
	case msg := <-deliverTxSub.Out(): // The tx was included in a block.
		deliverTxRes := msg.Data().(types.EventDataTx)
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxRes,
			DeliverTx: deliverTxRes.Result,
			Hash:      rtx.Hash(),
			Height:    deliverTxRes.Height,
		}, nil
	case <-deliverTxSub.Cancelled():
		var reason string
		if deliverTxSub.Err() == nil {
			reason = "Tendermint exited"
		} else {
			reason = deliverTxSub.Err().Error()
		}
		err = fmt.Errorf("deliverTxSub was cancelled (reason: %s)", reason)
		env.Logger.Error("Error on broadcastTxCommit", "err", err)
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxRes,
			DeliverTx: abci.ResponseDeliverTx{},
			Hash:      rtx.Hash(),
		}, err
	case <-time.After(env.Config.TimeoutBroadcastTxCommit):
		err = errors.New("timed out waiting for tx to be included in a block")
		env.Logger.Error("Error on broadcastTxCommit", "err", err)
		return &ctypes.ResultBroadcastTxCommit{
			CheckTx:   *checkTxRes,
			DeliverTx: abci.ResponseDeliverTx{},
			Hash:      rtx.Hash(),
		}, err
	}
}

func BroadcastBrczeroTxsAsync(ctx *rpctypes.Context, btcHeight int64, brczeroTxs []types.BRCZeroRequestTx) (*ctypes.ResultBroadcastTx, error) {
	txs := make([]types.Tx, 0)
	for _, s := range brczeroTxs {
		tx, err := rlp.EncodeToBytes(s)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	err := env.Mempool.AddBrczeroData(btcHeight, txs)

	if err != nil {
		return nil, err
	}
	//todo: construct resp
	return &ctypes.ResultBroadcastTx{}, nil
}

// UnconfirmedTxs gets unconfirmed transactions (maximum ?limit entries)
// including their number.
// More: https://docs.tendermint.com/master/rpc/#/Info/unconfirmed_txs
func UnconfirmedTxs(ctx *rpctypes.Context, limit int) (*ctypes.ResultUnconfirmedTxs, error) {

	txs := env.Mempool.ReapMaxTxs(limit)
	return &ctypes.ResultUnconfirmedTxs{
		Count:      len(txs),
		Total:      env.Mempool.Size(),
		TotalBytes: env.Mempool.TxsBytes(),
		Txs:        txs}, nil
}

// NumUnconfirmedTxs gets number of unconfirmed transactions.
// More: https://docs.tendermint.com/master/rpc/#/Info/num_unconfirmed_txs
func NumUnconfirmedTxs(ctx *rpctypes.Context) (*ctypes.ResultUnconfirmedTxs, error) {
	return &ctypes.ResultUnconfirmedTxs{
		Count:      env.Mempool.Size(),
		Total:      env.Mempool.Size(),
		TotalBytes: env.Mempool.TxsBytes()}, nil
}

func TxSimulateGasCost(ctx *rpctypes.Context, hash string) (*ctypes.ResponseTxSimulateGas, error) {
	return &ctypes.ResponseTxSimulateGas{
		GasCost: env.Mempool.GetTxSimulateGas(hash),
	}, nil
}

func UserUnconfirmedTxs(address string, limit int) (*ctypes.ResultUserUnconfirmedTxs, error) {
	txs := env.Mempool.ReapUserTxs(address, limit)
	return &ctypes.ResultUserUnconfirmedTxs{
		Count: len(txs),
		Txs:   txs}, nil
}

func TmUserUnconfirmedTxs(ctx *rpctypes.Context, address string, limit int) (*ctypes.ResultUserUnconfirmedTxs, error) {
	return UserUnconfirmedTxs(address, limit)
}

func UserNumUnconfirmedTxs(address string) (*ctypes.ResultUserUnconfirmedTxs, error) {
	nums := env.Mempool.ReapUserTxsCnt(address)
	return &ctypes.ResultUserUnconfirmedTxs{
		Count: nums}, nil
}

func TmUserNumUnconfirmedTxs(ctx *rpctypes.Context, address string) (*ctypes.ResultUserUnconfirmedTxs, error) {
	return UserNumUnconfirmedTxs(address)
}

func GetUnconfirmedTxByHash(hash [sha256.Size]byte) (types.Tx, error) {
	return env.Mempool.GetTxByHash(hash)
}

func GetAddressList() (*ctypes.ResultUnconfirmedAddresses, error) {
	addressList := env.Mempool.GetAddressList()
	return &ctypes.ResultUnconfirmedAddresses{
		Addresses: addressList,
	}, nil
}

func TmGetAddressList(ctx *rpctypes.Context) (*ctypes.ResultUnconfirmedAddresses, error) {
	return GetAddressList()
}

func GetPendingNonce(address string) (*ctypes.ResultPendingNonce, bool) {
	nonce, ok := env.Mempool.GetPendingNonce(address)
	if !ok {
		return nil, false
	}
	return &ctypes.ResultPendingNonce{
		Nonce: nonce,
	}, true
}

func GetEnableDeleteMinGPTx(ctx *rpctypes.Context) (*ctypes.ResultEnableDeleteMinGPTx, error) {
	status := env.Mempool.GetEnableDeleteMinGPTx()
	return &ctypes.ResultEnableDeleteMinGPTx{Enable: status}, nil
}

func GetPendingTxs(ctx *rpctypes.Context) (*ctypes.ResultPendingTxs, error) {
	pendingTx := make(map[string]map[string]types.WrappedMempoolTx)
	if baseapp.IsMempoolEnablePendingPool() {
		pendingTx = env.Mempool.GetPendingPoolTxsBytes()
	}
	return &ctypes.ResultPendingTxs{Txs: pendingTx}, nil
}
