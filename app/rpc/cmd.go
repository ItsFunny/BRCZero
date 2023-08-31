package rpc

import (
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/client/flags"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/client/lcd"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec/types"
	"github.com/spf13/cobra"
)

// ServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func ServeCmd(cdc *codec.CodecProxy, reg types.InterfaceRegistry) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, reg, RegisterRoutes)
	cmd.Flags().String(flagUnlockKey, "", "Select a key to unlock on the RPC server")
	cmd.Flags().String(FlagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	return cmd
}
