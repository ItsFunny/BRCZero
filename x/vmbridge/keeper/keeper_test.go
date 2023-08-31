package keeper_test

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/brc20-collab/brczero/app"
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	abci "github.com/brc20-collab/brczero/libs/tendermint/abci/types"
	"github.com/brc20-collab/brczero/libs/tendermint/types"
	evmtypes "github.com/brc20-collab/brczero/x/evm/types"
	"github.com/brc20-collab/brczero/x/vmbridge/keeper"
	wasmtypes "github.com/brc20-collab/brczero/x/wasm/types"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

//go:embed testdata/erc20abi.json
var erc20abiBytes []byte
var initCoin = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *app.BRCZeroApp

	keeper *keeper.Keeper

	addr         sdk.AccAddress
	wasmContract sdk.WasmAddress
	codeId       uint64

	evmContract common.Address

	freeCallWasmContract sdk.WasmAddress
	freeCallWasmCodeId   uint64
	freeCallEvmContract  common.Address

	evmABI abi.ABI
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.NewContext(checkTx, abci.Header{
		Height:  2,
		ChainID: "ethermint-3",
		Time:    time.Now().UTC(),
	})
	suite.keeper = suite.app.VMBridgeKeeper
	types.UnittestOnlySetMilestoneEarthHeight(1)

	suite.addr = sdk.AccAddress{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20}
	acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, suite.addr)
	err := acc.SetCoins(initCoin)
	suite.Require().NoError(err)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

	suite.app.WasmKeeper.SetParams(suite.ctx, wasmtypes.TestParams())
	evmParams := evmtypes.DefaultParams()
	evmParams.EnableCreate = true
	evmParams.EnableCall = true
	suite.app.EvmKeeper.SetParams(suite.ctx, evmParams)
	wasmcode, err := ioutil.ReadFile("./testdata/cw20.wasm")
	if err != nil {
		panic(err)
	}
	freeCallWasmCode, err := ioutil.ReadFile("./testdata/freecall.wasm")
	if err != nil {
		panic(err)
	}

	suite.codeId, err = suite.app.WasmPermissionKeeper.Create(suite.ctx, sdk.AccToAWasmddress(suite.addr), wasmcode, nil)
	suite.Require().NoError(err)
	suite.freeCallWasmCodeId, err = suite.app.WasmPermissionKeeper.Create(suite.ctx, sdk.AccToAWasmddress(suite.addr), freeCallWasmCode, nil)
	suite.Require().NoError(err)

	initMsg := []byte(fmt.Sprintf("{\"decimals\":10,\"initial_balances\":[{\"address\":\"%s\",\"amount\":\"100000000\"}],\"name\":\"my test token\", \"symbol\":\"MTT\"}", suite.addr.String()))
	suite.wasmContract, _, err = suite.app.WasmPermissionKeeper.Instantiate(suite.ctx, suite.codeId, sdk.AccToAWasmddress(suite.addr), sdk.AccToAWasmddress(suite.addr), initMsg, "label", sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)})
	suite.Require().NoError(err)

	palyload := "60806040526040518060600160405280603d815260200162002355603d913960079080519060200190620000359291906200004a565b503480156200004357600080fd5b506200015f565b8280546200005890620000fa565b90600052602060002090601f0160209004810192826200007c5760008555620000c8565b82601f106200009757805160ff1916838001178555620000c8565b82800160010185558215620000c8579182015b82811115620000c7578251825591602001919060010190620000aa565b5b509050620000d79190620000db565b5090565b5b80821115620000f6576000816000905550600101620000dc565b5090565b600060028204905060018216806200011357607f821691505b602082108114156200012a576200012962000130565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6121e6806200016f6000396000f3fe608060405234801561001057600080fd5b50600436106101215760003560e01c806370a08231116100ad578063cc1207c011610071578063cc1207c014610330578063d069cf761461034c578063d241877c1461037c578063dd62ed3e1461039a578063ee366654146103ca57610121565b806370a08231146102665780638e155cee1461029657806395d89b41146102b2578063a457c2d7146102d0578063a9059cbb1461030057610121565b8063313ce567116100f4578063313ce567146101c257806335b2bd2d146101e057806339509351146101fe5780633a0c76ea1461022e57806340c10f191461024a57610121565b806306fdde0314610126578063095ea7b31461014457806318160ddd1461017457806323b872dd14610192575b600080fd5b61012e6103e8565b60405161013b91906119a9565b60405180910390f35b61015e600480360381019061015991906114c5565b61047a565b60405161016b919061198e565b60405180910390f35b61017c610496565b6040516101899190611b70565b60405180910390f35b6101ac60048036038101906101a79190611472565b6104a0565b6040516101b9919061198e565b60405180910390f35b6101ca6104c8565b6040516101d79190611b8b565b60405180910390f35b6101e86104df565b6040516101f59190611973565b60405180910390f35b610218600480360381019061021391906114c5565b6104f7565b604051610225919061198e565b60405180910390f35b610248600480360381019061024391906115c2565b61059a565b005b610264600480360381019061025f91906114c5565b6105e4565b005b610280600480360381019061027b9190611405565b6105f2565b60405161028d9190611b70565b60405180910390f35b6102b060048036038101906102ab9190611579565b61063b565b005b6102ba610655565b6040516102c791906119a9565b60405180910390f35b6102ea60048036038101906102e591906114c5565b6106e7565b6040516102f7919061198e565b60405180910390f35b61031a600480360381019061031591906114c5565b6107ca565b604051610327919061198e565b60405180910390f35b61034a6004803603810190610345919061164d565b6107e6565b005b61036660048036038101906103619190611505565b6107f5565b604051610373919061198e565b60405180910390f35b6103846108b4565b60405161039191906119a9565b60405180910390f35b6103b460048036038101906103af9190611432565b610942565b6040516103c19190611b70565b60405180910390f35b6103d26109c9565b6040516103df91906119a9565b60405180910390f35b6060600180546103f790611d59565b80601f016020809104026020016040519081016040528092919081815260200182805461042390611d59565b80156104705780601f1061044557610100808354040283529160200191610470565b820191906000526020600020905b81548152906001019060200180831161045357829003601f168201915b5050505050905090565b60008033905061048b8185856109d8565b600191505092915050565b6000600454905090565b6000803390506104b1858285610ba3565b6104bc858585610c2f565b60019150509392505050565b6000600360009054906101000a900460ff16905090565b73c63cf6c8e1f3df41085e9d8af49584dae1432b4f81565b60008033905061058f818585600660008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205461058a9190611c38565b6109d8565b600191505092915050565b6105a43382610e9d565b7f41e4c36823b869e11ae85a7e623a332d31d961ba9ed670a3c9cb71c973c53caa8284836040516105d7939291906119cb565b60405180910390a1505050565b6105ee828261105e565b5050565b6000600560008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b806007908051906020019061065192919061125d565b5050565b60606002805461066490611d59565b80601f016020809104026020016040519081016040528092919081815260200182805461069090611d59565b80156106dd5780601f106106b2576101008083540402835291602001916106dd565b820191906000526020600020905b8154815290600101906020018083116106c057829003601f168201915b5050505050905090565b6000803390506000600660008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050838110156107b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107a890611b30565b60405180910390fd5b6107be82868684036109d8565b60019250505092915050565b6000803390506107db818585610c2f565b600191505092915050565b6107f18283836111a7565b5050565b600073c63cf6c8e1f3df41085e9d8af49584dae1432b4f73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461084357600080fd5b6007604051602001610855919061195c565b60405160208183030381529060405280519060200120858560405160200161087e929190611943565b604051602081830303815290604052805190602001201461089e57600080fd5b6108a8838361105e565b60019050949350505050565b600780546108c190611d59565b80601f01602080910402602001604051908101604052809291908181526020018280546108ed90611d59565b801561093a5780601f1061090f5761010080835404028352916020019161093a565b820191906000526020600020905b81548152906001019060200180831161091d57829003601f168201915b505050505081565b6000600660008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60606109d3610655565b905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610a48576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a3f90611b10565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610ab8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aaf90611a50565b60405180910390fd5b80600660008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92583604051610b969190611b70565b60405180910390a3505050565b6000610baf8484610942565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8114610c295781811015610c1b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c1290611a90565b60405180910390fd5b610c2884848484036109d8565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610c9f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c9690611af0565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610d0f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d0690611a10565b60405180910390fd5b6000600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610d96576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d8d90611ab0565b60405180910390fd5b818103600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610e2b9190611c38565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610e8f9190611b70565b60405180910390a350505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610f0d576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f0490611ad0565b60405180910390fd5b6000600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610f94576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f8b90611a30565b60405180910390fd5b818103600560008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508160046000828254610fec9190611c8e565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516110519190611b70565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156110ce576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110c590611b50565b60405180910390fd5b80600460008282546110e09190611c38565b9250508190555080600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546111369190611c38565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161119b9190611b70565b60405180910390a35050565b60008054906101000a900460ff16156111f5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111ec90611a70565b60405180910390fd5b60016000806101000a81548160ff021916908315150217905550826001908051906020019061122592919061125d565b50816002908051906020019061123c92919061125d565b5080600360006101000a81548160ff021916908360ff160217905550505050565b82805461126990611d59565b90600052602060002090601f01602090048101928261128b57600085556112d2565b82601f106112a457805160ff19168380011785556112d2565b828001600101855582156112d2579182015b828111156112d15782518255916020019190600101906112b6565b5b5090506112df91906112e3565b5090565b5b808211156112fc5760008160009055506001016112e4565b5090565b600061131361130e84611bcb565b611ba6565b90508281526020810184848401111561132f5761132e611e58565b5b61133a848285611d17565b509392505050565b6000813590506113518161216b565b92915050565b60008083601f84011261136d5761136c611e4e565b5b8235905067ffffffffffffffff81111561138a57611389611e49565b5b6020830191508360018202830111156113a6576113a5611e53565b5b9250929050565b600082601f8301126113c2576113c1611e4e565b5b81356113d2848260208601611300565b91505092915050565b6000813590506113ea81612182565b92915050565b6000813590506113ff81612199565b92915050565b60006020828403121561141b5761141a611e62565b5b600061142984828501611342565b91505092915050565b6000806040838503121561144957611448611e62565b5b600061145785828601611342565b925050602061146885828601611342565b9150509250929050565b60008060006060848603121561148b5761148a611e62565b5b600061149986828701611342565b93505060206114aa86828701611342565b92505060406114bb868287016113db565b9150509250925092565b600080604083850312156114dc576114db611e62565b5b60006114ea85828601611342565b92505060206114fb858286016113db565b9150509250929050565b6000806000806060858703121561151f5761151e611e62565b5b600085013567ffffffffffffffff81111561153d5761153c611e5d565b5b61154987828801611357565b9450945050602061155c87828801611342565b925050604061156d878288016113db565b91505092959194509250565b60006020828403121561158f5761158e611e62565b5b600082013567ffffffffffffffff8111156115ad576115ac611e5d565b5b6115b9848285016113ad565b91505092915050565b6000806000606084860312156115db576115da611e62565b5b600084013567ffffffffffffffff8111156115f9576115f8611e5d565b5b611605868287016113ad565b935050602084013567ffffffffffffffff81111561162657611625611e5d565b5b611632868287016113ad565b9250506040611643868287016113db565b9150509250925092565b6000806040838503121561166457611663611e62565b5b600083013567ffffffffffffffff81111561168257611681611e5d565b5b61168e858286016113ad565b925050602061169f858286016113f0565b9150509250929050565b6116b281611cc2565b82525050565b6116c181611cd4565b82525050565b60006116d38385611c2d565b93506116e0838584611d17565b82840190509392505050565b60006116f782611c11565b6117018185611c1c565b9350611711818560208601611d26565b61171a81611e67565b840191505092915050565b6000815461173281611d59565b61173c8186611c2d565b9450600182166000811461175757600181146117685761179b565b60ff1983168652818601935061179b565b61177185611bfc565b60005b8381101561179357815481890152600182019150602081019050611774565b838801955050505b50505092915050565b60006117b1602383611c1c565b91506117bc82611e78565b604082019050919050565b60006117d4602283611c1c565b91506117df82611ec7565b604082019050919050565b60006117f7602283611c1c565b915061180282611f16565b604082019050919050565b600061181a601b83611c1c565b915061182582611f65565b602082019050919050565b600061183d601d83611c1c565b915061184882611f8e565b602082019050919050565b6000611860602683611c1c565b915061186b82611fb7565b604082019050919050565b6000611883602183611c1c565b915061188e82612006565b604082019050919050565b60006118a6602583611c1c565b91506118b182612055565b604082019050919050565b60006118c9602483611c1c565b91506118d4826120a4565b604082019050919050565b60006118ec602583611c1c565b91506118f7826120f3565b604082019050919050565b600061190f601f83611c1c565b915061191a82612142565b602082019050919050565b61192e81611d00565b82525050565b61193d81611d0a565b82525050565b60006119508284866116c7565b91508190509392505050565b60006119688284611725565b915081905092915050565b600060208201905061198860008301846116a9565b92915050565b60006020820190506119a360008301846116b8565b92915050565b600060208201905081810360008301526119c381846116ec565b905092915050565b600060608201905081810360008301526119e581866116ec565b905081810360208301526119f981856116ec565b9050611a086040830184611925565b949350505050565b60006020820190508181036000830152611a29816117a4565b9050919050565b60006020820190508181036000830152611a49816117c7565b9050919050565b60006020820190508181036000830152611a69816117ea565b9050919050565b60006020820190508181036000830152611a898161180d565b9050919050565b60006020820190508181036000830152611aa981611830565b9050919050565b60006020820190508181036000830152611ac981611853565b9050919050565b60006020820190508181036000830152611ae981611876565b9050919050565b60006020820190508181036000830152611b0981611899565b9050919050565b60006020820190508181036000830152611b29816118bc565b9050919050565b60006020820190508181036000830152611b49816118df565b9050919050565b60006020820190508181036000830152611b6981611902565b9050919050565b6000602082019050611b856000830184611925565b92915050565b6000602082019050611ba06000830184611934565b92915050565b6000611bb0611bc1565b9050611bbc8282611d8b565b919050565b6000604051905090565b600067ffffffffffffffff821115611be657611be5611e1a565b5b611bef82611e67565b9050602081019050919050565b60008190508160005260206000209050919050565b600081519050919050565b600082825260208201905092915050565b600081905092915050565b6000611c4382611d00565b9150611c4e83611d00565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611c8357611c82611dbc565b5b828201905092915050565b6000611c9982611d00565b9150611ca483611d00565b925082821015611cb757611cb6611dbc565b5b828203905092915050565b6000611ccd82611ce0565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b82818337600083830152505050565b60005b83811015611d44578082015181840152602081019050611d29565b83811115611d53576000848401525b50505050565b60006002820490506001821680611d7157607f821691505b60208210811415611d8557611d84611deb565b5b50919050565b611d9482611e67565b810181811067ffffffffffffffff82111715611db357611db2611e1a565b5b80604052505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a206275726e20616d6f756e7420657863656564732062616c616e60008201527f6365000000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20616c726561647920696e697469616c697a65643b0000000000600082015250565b7f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000600082015250565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a206275726e2066726f6d20746865207a65726f2061646472657360008201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b61217481611cc2565b811461217f57600080fd5b50565b61218b81611d00565b811461219657600080fd5b50565b6121a281611d0a565b81146121ad57600080fd5b5056fea264697066735822122022151d41dd9654958225e494d6adc336a64e7ff8fe7993ba1764f69906eb984064736f6c6343000807003365783134686a32746176713866706573647778786375343472747933686839307668756a7276636d73746c347a723374786d6676773973366671753237"
	bytescode := common.Hex2Bytes(palyload)
	_, r2, err := suite.app.VMBridgeKeeper.CallEvm(suite.ctx, common.BytesToAddress(suite.addr), nil, big.NewInt(0), bytescode)
	suite.Require().NoError(err)
	suite.evmContract = r2.ContractAddress

	freeCallPalyload := "608060405234801561001057600080fd5b50610893806100206000396000f3fe6080604052600436106100345760003560e01c806335b2bd2d146100395780635d78dad01461006457806382c11cef14610094575b600080fd5b34801561004557600080fd5b5061004e6100d1565b60405161005b91906103aa565b60405180910390f35b61007e6004803603810190610079919061051f565b6100e9565b60405161008b9190610616565b60405180910390f35b3480156100a057600080fd5b506100bb60048036038101906100b6919061066e565b61018b565b6040516100c89190610714565b60405180910390f35b731033796b018b2bf0fc9cb88c0793b2f275edb62481565b6060600061012c6040518060400160405280601381526020017f63616c6c42795761736d2072657475726e3a2000000000000000000000000000815250856101d3565b9050600061016f826040518060400160405280600a81526020017f202d2d2d646174613a20000000000000000000000000000000000000000000008152506101d3565b9050600061017d82866101d3565b905080935050505092915050565b60007fcca73dc0c9131f3d7540642f5b7bc76eceaedddf94108f54b7a7c9e594d967bf8484846040516101c09392919061073e565b60405180910390a1600190509392505050565b6060600083905060008390506000815183516101ef91906107b2565b67ffffffffffffffff811115610208576102076103f4565b5b6040519080825280601f01601f19166020018201604052801561023a5781602001600182028036833780820191505090505b50905060008190506000805b85518110156102ce57858181518110610262576102616107e6565b5b602001015160f81c60f81b83838061027990610815565b94508151811061028c5761028b6107e6565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535080806102c690610815565b915050610246565b5060005b845181101561035a578481815181106102ee576102ed6107e6565b5b602001015160f81c60f81b83838061030590610815565b945081518110610318576103176107e6565b5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350808061035290610815565b9150506102d2565b50829550505050505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061039482610369565b9050919050565b6103a481610389565b82525050565b60006020820190506103bf600083018461039b565b92915050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61042c826103e3565b810181811067ffffffffffffffff8211171561044b5761044a6103f4565b5b80604052505050565b600061045e6103c5565b905061046a8282610423565b919050565b600067ffffffffffffffff82111561048a576104896103f4565b5b610493826103e3565b9050602081019050919050565b82818337600083830152505050565b60006104c26104bd8461046f565b610454565b9050828152602081018484840111156104de576104dd6103de565b5b6104e98482856104a0565b509392505050565b600082601f830112610506576105056103d9565b5b81356105168482602086016104af565b91505092915050565b60008060408385031215610536576105356103cf565b5b600083013567ffffffffffffffff811115610554576105536103d4565b5b610560858286016104f1565b925050602083013567ffffffffffffffff811115610581576105806103d4565b5b61058d858286016104f1565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b60005b838110156105d15780820151818401526020810190506105b6565b60008484015250505050565b60006105e882610597565b6105f281856105a2565b93506106028185602086016105b3565b61060b816103e3565b840191505092915050565b6000602082019050818103600083015261063081846105dd565b905092915050565b6000819050919050565b61064b81610638565b811461065657600080fd5b50565b60008135905061066881610642565b92915050565b600080600060608486031215610687576106866103cf565b5b600084013567ffffffffffffffff8111156106a5576106a46103d4565b5b6106b1868287016104f1565b93505060206106c286828701610659565b925050604084013567ffffffffffffffff8111156106e3576106e26103d4565b5b6106ef868287016104f1565b9150509250925092565b60008115159050919050565b61070e816106f9565b82525050565b60006020820190506107296000830184610705565b92915050565b61073881610638565b82525050565b6000606082019050818103600083015261075881866105dd565b9050610767602083018561072f565b818103604083015261077981846105dd565b9050949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006107bd82610638565b91506107c883610638565b92508282019050808211156107e0576107df610783565b5b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600061082082610638565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361085257610851610783565b5b60018201905091905056fea2646970667358221220b53efdfdf8b5a055fe5af52d3e108cbe12b020fb309c812c303b6ce6a159771964736f6c63430008120033"
	freeCallBytescode := common.Hex2Bytes(freeCallPalyload)
	_, freeCallR2, err := suite.app.VMBridgeKeeper.CallEvm(suite.ctx, common.BytesToAddress(suite.addr), nil, big.NewInt(0), freeCallBytescode)
	suite.Require().NoError(err)
	suite.freeCallEvmContract = freeCallR2.ContractAddress

	initFreeCallMsg := []byte(fmt.Sprintf("{\"decimals\":10,\"initial_balances\":[{\"address\":\"%s\",\"amount\":\"100000000\"},{\"address\":\"%s\",\"amount\":\"100000000\"}],\"name\":\"my test token\", \"symbol\":\"MTT\"}", sdk.AccToAWasmddress(suite.addr).String(), sdk.WasmAddress(suite.freeCallEvmContract.Bytes()).String()))
	suite.freeCallWasmContract, _, err = suite.app.WasmPermissionKeeper.Instantiate(suite.ctx, suite.freeCallWasmCodeId, sdk.AccToAWasmddress(suite.addr), sdk.AccToAWasmddress(suite.addr), initFreeCallMsg, "label", sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)})
	suite.Require().NoError(err)

	suite.evmABI, err = abi.JSON(bytes.NewReader(erc20abiBytes))
	suite.Require().NoError(err)

	//init
	init, err := suite.evmABI.Pack("mint", common.BytesToAddress(suite.addr.Bytes()), big.NewInt(1000))
	suite.Require().NoError(err)
	_, _, err = suite.app.VMBridgeKeeper.CallEvm(suite.ctx, common.BytesToAddress(suite.addr), &suite.evmContract, big.NewInt(0), init)
	suite.Require().NoError(err)

	update, err := suite.evmABI.Pack("updatewasmContractAddress", suite.wasmContract.String())
	suite.Require().NoError(err)
	_, _, err = suite.app.VMBridgeKeeper.CallEvm(suite.ctx, common.BytesToAddress(suite.addr), &suite.evmContract, big.NewInt(0), update)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) queryBalance(addr common.Address) *big.Int {
	update, err := suite.evmABI.Pack("balanceOf", addr)
	suite.Require().NoError(err)
	_, result, err := suite.app.VMBridgeKeeper.CallEvm(suite.ctx, common.BytesToAddress(suite.addr), &suite.evmContract, big.NewInt(0), update)
	r, err := suite.evmABI.Unpack("balanceOf", result.Ret)
	return r[0].(*big.Int)
}

func (suite *KeeperTestSuite) queryCoins(addr sdk.AccAddress) sdk.Coins {
	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, addr)
	if acc == nil {
		return sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0))
	}
	return acc.GetCoins()
}

func (suite *KeeperTestSuite) SetAccountCoins(addr sdk.AccAddress, value sdk.Int) {
	acc := suite.app.AccountKeeper.GetAccount(suite.ctx, addr)
	suite.Require().NotNil(acc)

	err := acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, value.Int64())))
	suite.Require().NoError(err)
	suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
}

func (suite *KeeperTestSuite) deployEvmContract(code string) common.Address {
	freeCallBytescode := common.Hex2Bytes(code)
	_, contract, err := suite.app.VMBridgeKeeper.CallEvm(suite.ctx, common.BytesToAddress(suite.addr), nil, big.NewInt(0), freeCallBytescode)
	suite.Require().NoError(err)
	return contract.ContractAddress
}

func (suite *KeeperTestSuite) deployWasmContract(filename string, initMsg []byte) sdk.WasmAddress {
	wasmcode, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s", filename))
	suite.Require().NoError(err)

	codeid, err := suite.app.WasmPermissionKeeper.Create(suite.ctx, sdk.AccToAWasmddress(suite.addr), wasmcode, nil)
	suite.Require().NoError(err)

	//initMsg := []byte(fmt.Sprintf("{\"decimals\":10,\"initial_balances\":[{\"address\":\"%s\",\"amount\":\"100000000\"}],\"name\":\"my test token\", \"symbol\":\"MTT\"}", suite.addr.String()))
	contract, _, err := suite.app.WasmPermissionKeeper.Instantiate(suite.ctx, codeid, sdk.AccToAWasmddress(suite.addr), sdk.AccToAWasmddress(suite.addr), initMsg, "label", sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)})
	suite.Require().NoError(err)
	return contract
}

func (suite *KeeperTestSuite) executeWasmContract(ctx sdk.Context, caller, wasmContract sdk.WasmAddress, msg []byte, amount sdk.Coins) []byte {
	ret, err := suite.app.WasmPermissionKeeper.Execute(ctx, wasmContract, caller, msg, amount)
	suite.Require().NoError(err)
	return ret
}
func (suite *KeeperTestSuite) queryEvmContract(ctx sdk.Context, addr common.Address, calldata []byte) ([]byte, error) {
	_, contract, err := suite.app.VMBridgeKeeper.CallEvm(ctx, common.BytesToAddress(suite.addr), &addr, big.NewInt(0), calldata)
	return contract.Ret, err
}

func (suite *KeeperTestSuite) queryWasmContract(caller string, calldata []byte) ([]byte, error) {
	subCtx, _ := suite.ctx.CacheContext()
	result, err := suite.app.VMBridgeKeeper.QueryToWasm(subCtx, caller, calldata)
	return result, err
}
