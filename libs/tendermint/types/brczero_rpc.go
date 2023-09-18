package types

import (
	sm "github.com/brc20-collab/brczero/libs/tendermint/state"
	"sync"
)

type UnconfirmedBTCTxRspCache struct {
	rsps map[int64]sm.ABCIResponses
	mtx  sync.RWMutex
}
