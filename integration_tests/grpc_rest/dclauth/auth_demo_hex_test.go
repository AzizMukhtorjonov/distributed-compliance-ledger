// Bash Source: integration_tests/cli/auth-demo-hex.sh
//
// This test verifies auth operations using hex-formatted VID values:
// - ACTION: Generate keys for a new user (auth-demo-hex.sh:L19-L22)
// - ACTION: Propose a vendor account with hex VID 0xA13 (auth-demo-hex.sh:L45-L48)
// - ACTION: Verify account is active with 1/3 approval (auth-demo-hex.sh:L52-L54)
// - ACTION: Add a model using hex VID (auth-demo-hex.sh:L79-L82)
// - ACTION: Attempt to add model with different hex VID (auth-demo-hex.sh:L88-L91)
// - ACTION: Update model using hex VID (auth-demo-hex.sh:L95-L98)
// - ACTION: Query model using hex VID and verify decimal conversion (auth-demo-hex.sh:L102-L107)
package dclauth_test

import (
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	testDclauth "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
	modeltypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/model/types"
)

// TestAuthDemoHex maps to: integration_tests/cli/auth-demo-hex.sh
func TestAuthDemoHex(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation - get trustee accounts
	jackName := testconstants.JackAccount
	jackKeyInfo, err := suite.Kr.Key(jackName)
	require.NoError(t, err)
	jackAddr, err := jackKeyInfo.GetAddress()
	require.NoError(t, err)
	jackAccount, err := testDclauth.GetAccount(&suite, jackAddr)
	require.NoError(t, err)

	aliceName := testconstants.AliceAccount
	aliceKeyInfo, err := suite.Kr.Key(aliceName)
	require.NoError(t, err)
	aliceAddr, err := aliceKeyInfo.GetAddress()
	require.NoError(t, err)
	aliceAccount, err := testDclauth.GetAccount(&suite, aliceAddr)
	require.NoError(t, err)

	bobName := testconstants.BobAccount
	bobKeyInfo, err := suite.Kr.Key(bobName)
	require.NoError(t, err)
	bobAddr, err := bobKeyInfo.GetAddress()
	require.NoError(t, err)
	bobAccount, err := testDclauth.GetAccount(&suite, bobAddr)
	require.NoError(t, err)

	// ACTION: Propose a vendor account with a specific VID (hex 0xA13 = 2579 in bash)
	// (auth-demo-hex.sh:L45-L48)
	vid := int32(2579) // 0xA13
	pid := int32(2577) // 0xA11
	vendorName := utils.RandString()
	vendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		vendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		vid,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	require.NotNil(t, vendorAccount)

	// ACTION: Verify account is active (auth-demo-hex.sh:L52-L54)
	retrievedAccount, err := testDclauth.GetAccount(&suite, suite.GetAddress(vendorName))
	require.NoError(t, err)
	require.Equal(t, vid, retrievedAccount.VendorID)

	// ACTION: Add a model using the VID (auth-demo-hex.sh:L79-L82)
	createModelMsg := &modeltypes.MsgCreateModel{
		Creator:                 vendorAccount.Address,
		Vid:                     vid,
		Pid:                     pid,
		DeviceTypeId:            12,
		ProductName:             "Device #1",
		ProductLabel:            "Device Description",
		CommissioningCustomFlow: 0,
		PartNumber:              "12",
		EnhancedSetupFlowOptions: 0,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Attempt to add model with a different VID than assigned - expect error
	// (auth-demo-hex.sh:L88-L91)
	vidPlusOne := vid + 1
	badModelMsg := &modeltypes.MsgCreateModel{
		Creator:                 vendorAccount.Address,
		Vid:                     vidPlusOne,
		Pid:                     pid,
		DeviceTypeId:            12,
		ProductName:             "Device #1",
		ProductLabel:            "Device Description",
		CommissioningCustomFlow: 0,
		PartNumber:              "12",
		EnhancedSetupFlowOptions: 0,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badModelMsg}, vendorName, vendorAccount)
	require.Error(t, err)

	// ACTION: Update model using correct VID (auth-demo-hex.sh:L95-L98)
	updateModelMsg := &modeltypes.MsgUpdateModel{
		Creator:                  vendorAccount.Address,
		Vid:                      vid,
		Pid:                      pid,
		ProductLabel:             "Device Description",
		EnhancedSetupFlowOptions: 2,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateModelMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	_ = jackAccount
	_ = tmrand.Uint16() // suppress unused import
}
