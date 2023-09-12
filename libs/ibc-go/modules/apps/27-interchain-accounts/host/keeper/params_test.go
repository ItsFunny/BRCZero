package keeper_test

import "github.com/brc20-collab/brczero/libs/ibc-go/modules/apps/27-interchain-accounts/host/types"

func (suite *KeeperTestSuite) TestParams() {
	expParams := types.DefaultParams()

	params := suite.chainA.GetSimApp().ICAHostKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)

	expParams.HostEnabled = false
	expParams.AllowMessages = []string{"/cosmos.staking.v1beta1.MsgDelegate"}
	suite.chainA.GetSimApp().ICAHostKeeper.SetParams(suite.chainA.GetContext(), expParams)
	params = suite.chainA.GetSimApp().ICAHostKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)
}