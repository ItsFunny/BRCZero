package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/brc20-collab/brczero/libs/cosmos-sdk/client/context"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/types/rest"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/x/auth/types"
	ttypes "github.com/brc20-collab/brczero/libs/tendermint/types"
)

// BroadcastReq defines a tx broadcasting request.
type BroadcastReq struct {
	Tx    types.StdTx `json:"tx" yaml:"tx"`
	Mode  string      `json:"mode" yaml:"mode"`
	Nonce uint64      `json:"nonce" yaml:"nonce"`
}

// BroadcastTxRequest implements a tx broadcasting handler that is responsible
// for broadcasting a valid and signed tx to a full node. The tx can be
// broadcasted via a sync|async|block mechanism.
func BroadcastTxRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		if req.Nonce != 0 {
			wcmt := &ttypes.WrapCMTx{
				Tx:    txBytes,
				Nonce: req.Nonce,
			}
			data, err := cliCtx.Codec.MarshalJSON(wcmt)
			if err == nil {
				txBytes = data
			}
		}

		cliCtx = cliCtx.WithBroadcastMode(req.Mode)

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}