package model_test

import (
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	testDclauth "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/model"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
	commontypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/common/types"
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestModelVersionDemoRewrite(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation of Actors
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

	vid := int32(tmrand.Uint16())
	pid := int32(tmrand.Uint16())

	// ACTION: Create a vendor account for the specific VID
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

	// ACTION: Add a base Model record as a prerequisite for versioning
	createModelMsg := model.NewMsgCreateModel(vid, pid, vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	sv := uint32(tmrand.Uint16())
	// ACTION: Create a first Model Version record
	createVersionMsg := model.NewMsgCreateModelVersion(vid, pid, sv, "1", vendorAccount.Address)
	createVersionMsg.MinApplicableSoftwareVersion = 1
	createVersionMsg.MaxApplicableSoftwareVersion = 10
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createVersionMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query the newly created Model Version and verify its fields
	receivedVersion, err := model.GetModelVersion(&suite, vid, pid, sv)
	require.NoError(t, err)
	require.Equal(t, vid, receivedVersion.Vid)
	require.Equal(t, pid, receivedVersion.Pid)
	require.Equal(t, sv, receivedVersion.SoftwareVersion)
	require.Equal(t, "1", receivedVersion.SoftwareVersionString)

	// ACTION: Query all versions for this Model and verify the version is present
	allVersions, err := model.GetModelVersions(&suite, vid, pid)
	require.NoError(t, err)
	require.Contains(t, allVersions.SoftwareVersions, sv)

	// ACTION: Query a non-existent version and expect "Not Found"
	_, err = model.GetModelVersion(&suite, vid, pid, 123456)
	require.Error(t, err)
	stat, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())

	// ACTION: Query versions for a non-existent model and expect "Not Found"
	_, err = model.GetModelVersions(&suite, int32(tmrand.Uint16()), int32(tmrand.Uint16()))
	require.Error(t, err)
	stat, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())

	// ACTION: Update the existing Model Version record
	updateMsg := model.NewMsgUpdateModelVersion(vid, pid, sv, vendorAccount.Address)
	updateMsg.MinApplicableSoftwareVersion = 2
	updateMsg.MaxApplicableSoftwareVersion = 10
	updateMsg.SoftwareVersionValid = false
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify the Model Version fields after update
	receivedVersion, err = model.GetModelVersion(&suite, vid, pid, sv)
	require.NoError(t, err)
	require.False(t, receivedVersion.SoftwareVersionValid)
	require.Equal(t, uint32(2), receivedVersion.MinApplicableSoftwareVersion)

	// ACTION: Create a second Model Version for the same Model
	sv2 := uint32(tmrand.Uint16())
	for sv2 == sv {
	    sv2 = uint32(tmrand.Uint16())
	}
	createVersionMsg2 := model.NewMsgCreateModelVersion(vid, pid, sv2, "2", vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createVersionMsg2}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query all versions and verify both are listed
	allVersions, err = model.GetModelVersions(&suite, vid, pid)
	require.NoError(t, err)
	require.Contains(t, allVersions.SoftwareVersions, sv)
	require.Contains(t, allVersions.SoftwareVersions, sv2)

	// ACTION: Create another vendor account (not authorized for original VID)
	newVid := int32(tmrand.Uint16())
	differentVendorName := utils.RandString()
	differentVendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		differentVendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		newVid,
		[]*commontypes.Uint16Range{},
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)

	// ACTION: Attempt to create a version for original VID using unauthorized account
	badCreateMsg := model.NewMsgCreateModelVersion(vid, pid, sv, "1", differentVendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badCreateMsg}, differentVendorName, differentVendorAccount)
	require.Error(t, err)

	// ACTION: Attempt to update version for original VID using unauthorized account
	badUpdateMsg := model.NewMsgUpdateModelVersion(vid, pid, sv, differentVendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badUpdateMsg}, differentVendorName, differentVendorAccount)
	require.Error(t, err)

	// ACTION: Delete the first Model Version
	deleteMsg := model.NewMsgDeleteModelVersion(vid, pid, sv, vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{deleteMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify the Model Version is no longer found after deletion
	_, err = model.GetModelVersion(&suite, vid, pid, sv)
	require.Error(t, err)
	stat, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())
}
