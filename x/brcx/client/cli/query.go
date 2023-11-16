package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/brc20-collab/brczero/libs/cosmos-sdk/client"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/client/context"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/client/flags"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec"
	"github.com/brc20-collab/brczero/x/brcx/internal/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group brcx queries under a subcommand
	slashingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the brcx module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	slashingQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQuerySigningInfo(queryRoute, cdc),
		)...,
	)

	return slashingQueryCmd

}

// GetCmdQuerySigningInfo implements the command to query signing info.
func GetCmdQuerySigningInfo(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		//todo
		Use:   "signing-info [validator-conspub]",
		Short: "Query a validator's signing information",
		Long: strings.TrimSpace(`Use a validators' consensus public key to find the signing-info for that validator:

$ <appcli> query slashing signing-info exvalconspub1zcjduepqfhvwcmt7p06fvdgexxhmz0l8c7sgswl7ulv7aulk364x4g5xsw7sr0k2g5
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			//
			//pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, args[0])
			//if err != nil {
			//	return err
			//}

			//consAddr := sdk.ConsAddress(pk.Address())
			//key := types.GetValidatorSigningInfoKey(consAddr)
			//
			//res, _, err := cliCtx.QueryStore(key, storeName)
			//if err != nil {
			//	return err
			//}
			//
			//if len(res) == 0 {
			//	return fmt.Errorf("validator %s not found in slashing store", consAddr)
			//}
			//
			//var signingInfo types.ValidatorSigningInfo
			//cdc.MustUnmarshalBinaryLengthPrefixed(res, &signingInfo)
			return cliCtx.PrintOutput(nil)
		},
	}
}
