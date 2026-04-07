// Bash Source: integration_tests/cli/modelversion-demo-hex.sh
//
// This test mirrors modelversion-demo.sh but uses hex-formatted VID/PID:
// - ACTION: Create vendor account with hex VID 0xA13 (modelversion-demo-hex.sh:L28-L29)
// - ACTION: Add model with hex VID/PID (modelversion-demo-hex.sh:L34-L37)
// - ACTION: Create model version with hex VID/PID (modelversion-demo-hex.sh:L42-L46)
// - ACTION: Query model version using hex VID/PID (modelversion-demo-hex.sh:L51-L61)
// - ACTION: Query all model versions using hex VID/PID (modelversion-demo-hex.sh:L66-L72)
// - ACTION: Query non-existent model version (modelversion-demo-hex.sh:L79-L81)
// - ACTION: Update model version using hex VID/PID (modelversion-demo-hex.sh:L95-L98)
// - ACTION: Add second model version (modelversion-demo-hex.sh:L118-L123)
// - ACTION: Cross-vendor creation attempt with hex VID (modelversion-demo-hex.sh:L141-L147)
// - ACTION: Cross-vendor update attempt with hex VID (modelversion-demo-hex.sh:L152-L155)
package model_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	testDclauth "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/model"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestModelVersionDemoHex maps to: integration_tests/cli/modelversion-demo-hex.sh
// ACTION: Full model version lifecycle using hex VID 0xA13 (2579) and PID 0xA11 (2577)
func TestModelVersionDemoHex(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation
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

	// ACTION: Create vendor account with hex VID 0xA13 (modelversion-demo-hex.sh:L28-L29)
	var vid int32 = 0xA13 // 2579
	var pid int32 = 0xA11 // 2577
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

	// ACTION: Add model with hex VID/PID (modelversion-demo-hex.sh:L34-L37)
	createModelMsg := model.NewMsgCreateModel(vid, pid, vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Create model version (modelversion-demo-hex.sh:L42-L46)
	var sv uint32 = 42
	createVersionMsg := model.NewMsgCreateModelVersion(vid, pid, sv, "1", vendorAccount.Address)
	createVersionMsg.MinApplicableSoftwareVersion = 1
	createVersionMsg.MaxApplicableSoftwareVersion = 10
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createVersionMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query model version using hex VID/PID (modelversion-demo-hex.sh:L51-L61)
	receivedVersion, err := model.GetModelVersion(&suite, vid, pid, sv)
	require.NoError(t, err)
	require.Equal(t, vid, receivedVersion.Vid)
	require.Equal(t, pid, receivedVersion.Pid)
	require.Equal(t, sv, receivedVersion.SoftwareVersion)
	require.Equal(t, "1", receivedVersion.SoftwareVersionString)
	require.True(t, receivedVersion.SoftwareVersionValid)
	require.Equal(t, uint32(1), receivedVersion.MinApplicableSoftwareVersion)
	require.Equal(t, uint32(10), receivedVersion.MaxApplicableSoftwareVersion)

	// ACTION: Query all model versions (modelversion-demo-hex.sh:L66-L72)
	allVersions, err := model.GetModelVersions(&suite, vid, pid)
	require.NoError(t, err)
	require.Contains(t, allVersions.SoftwareVersions, sv)

	// ACTION: Query non-existent model version (modelversion-demo-hex.sh:L79-L81)
	_, err = model.GetModelVersion(&suite, vid, pid, 123456)
	require.Error(t, err)
	stat, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())

	// ACTION: Update model version (modelversion-demo-hex.sh:L95-L98)
	updateMsg := model.NewMsgUpdateModelVersion(vid, pid, sv, vendorAccount.Address)
	updateMsg.MinApplicableSoftwareVersion = 2
	updateMsg.MaxApplicableSoftwareVersion = 10
	updateMsg.SoftwareVersionValid = false
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// Verify update
	receivedVersion, err = model.GetModelVersion(&suite, vid, pid, sv)
	require.NoError(t, err)
	require.False(t, receivedVersion.SoftwareVersionValid)
	require.Equal(t, uint32(2), receivedVersion.MinApplicableSoftwareVersion)

	// ACTION: Add second model version (modelversion-demo-hex.sh:L118-L123)
	var sv2 uint32 = 43
	createVersionMsg2 := model.NewMsgCreateModelVersion(vid, pid, sv2, "2", vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createVersionMsg2}, vendorName, vendorAccount)
	require.NoError(t, err)

	allVersions, err = model.GetModelVersions(&suite, vid, pid)
	require.NoError(t, err)
	require.Contains(t, allVersions.SoftwareVersions, sv)
	require.Contains(t, allVersions.SoftwareVersions, sv2)

	// ACTION: Cross-vendor creation attempt (modelversion-demo-hex.sh:L141-L147)
	var newVid int32 = 0xB17
	differentVendorName := utils.RandString()
	differentVendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		differentVendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		newVid,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	badCreateMsg := model.NewMsgCreateModelVersion(vid, pid, sv, "1", differentVendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badCreateMsg}, differentVendorName, differentVendorAccount)
	require.Error(t, err)

	// ACTION: Cross-vendor update attempt (modelversion-demo-hex.sh:L152-L155)
	badUpdateMsg := model.NewMsgUpdateModelVersion(vid, pid, sv, differentVendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badUpdateMsg}, differentVendorName, differentVendorAccount)
	require.Error(t, err)
}
