// Bash Source: integration_tests/cli/model-demo.sh
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
)

func TestModelDemo(t *testing.T) {
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

	vidWithPids := vid + 1
	pidRanges := []*commontypes.Uint16Range{{Min: pid, Max: pid}}
	// ACTION: Create another vendor account with assigned Product ID ranges
	vendorWithPidsName := utils.RandString()
	vendorWithPidsAccount := testDclauth.CreateVendorAccount(
		&suite,
		vendorWithPidsName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		vidWithPids,
		pidRanges,
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	require.NotNil(t, vendorWithPidsAccount)

	// ACTION: Query a model that does not exist and expect "Not Found"
	_, err = model.GetModel(&suite, vid, pid)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Query models for a vendor that has none and expect "Not Found"
	_, err = model.GetVendorModels(&suite, vid)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Query all models and expect an empty list
	allModels, err := model.GetModels(&suite)
	require.NoError(t, err)
	require.Empty(t, allModels)

	// ACTION: Add a new Model record with basic information
	productLabel := "Device #1"
	createMsg := model.NewMsgCreateModel(vid, pid, vendorAccount.Address)
	createMsg.ProductLabel = productLabel
	createMsg.SchemaVersion = 0
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Add another Model record with enhanced setup flow options
	enhancedSetupFlowTCUrl := "https://example.org/file.txt"
	enhancedSetupFlowTCRevision := int32(1)
	enhancedSetupFlowTCDigest := "MWRjNGE0NDA0MWRjYWYxMTU0NWI3NTQzZGZlOTQyZjQ3NDJmNTY4YmU2OGZlZTI3NTQ0MWIwOTJiYjYwZGVlZA=="
	enhancedSetupFlowTCFileSize := uint32(1024)
	maintenanceUrl := "https://example.org"
	commissioningFallbackUrl := "https://url.commissioningfallbackurl.dclmodel"
	discoveryCapabilitiesBitmask := uint32(1)

	createMsg2 := model.NewMsgCreateModel(vidWithPids, pid, vendorWithPidsAccount.Address)
	createMsg2.ProductLabel = productLabel
	createMsg2.SchemaVersion = 0
	createMsg2.EnhancedSetupFlowOptions = 1
	createMsg2.EnhancedSetupFlowTCUrl = enhancedSetupFlowTCUrl
	createMsg2.EnhancedSetupFlowTCRevision = enhancedSetupFlowTCRevision
	createMsg2.EnhancedSetupFlowTCDigest = enhancedSetupFlowTCDigest
	createMsg2.EnhancedSetupFlowTCFileSize = enhancedSetupFlowTCFileSize
	createMsg2.MaintenanceUrl = maintenanceUrl
	createMsg2.CommissioningFallbackUrl = commissioningFallbackUrl
	createMsg2.DiscoveryCapabilitiesBitmask = discoveryCapabilitiesBitmask

	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createMsg2}, vendorWithPidsName, vendorWithPidsAccount)
	require.NoError(t, err)

	// ACTION: Query the first Model and verify its fields
	receivedModel, err := model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, vid, receivedModel.Vid)
	require.Equal(t, pid, receivedModel.Pid)
	require.Equal(t, productLabel, receivedModel.ProductLabel)

	// ACTION: Query the second Model and verify its detailed fields
	receivedModel2, err := model.GetModel(&suite, vidWithPids, pid)
	require.NoError(t, err)
	require.Equal(t, vidWithPids, receivedModel2.Vid)
	require.Equal(t, pid, receivedModel2.Pid)
	require.Equal(t, enhancedSetupFlowTCUrl, receivedModel2.EnhancedSetupFlowTCUrl)
	require.Equal(t, enhancedSetupFlowTCRevision, receivedModel2.EnhancedSetupFlowTCRevision)
	require.Equal(t, maintenanceUrl, receivedModel2.MaintenanceUrl)

	// ACTION: Create a Model Version for the first model
	sv := uint32(1)
	createVersionMsg := model.NewMsgCreateModelVersion(vid, pid, sv, "1", vendorAccount.Address)
	createVersionMsg.CdVersionNumber = 10
	createVersionMsg.MinApplicableSoftwareVersion = 1
	createVersionMsg.MaxApplicableSoftwareVersion = 15
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createVersionMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Create a Model Version for the second model
	createVersionMsg2 := model.NewMsgCreateModelVersion(vidWithPids, pid, sv, "1", vendorWithPidsAccount.Address)
	createVersionMsg2.CdVersionNumber = 10
	createVersionMsg2.MinApplicableSoftwareVersion = 1
	createVersionMsg2.MaxApplicableSoftwareVersion = 15
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createVersionMsg2}, vendorWithPidsName, vendorWithPidsAccount)
	require.NoError(t, err)

	// ACTION: Query all models and verify both are listed
	allModels, err = model.GetModels(&suite)
	require.NoError(t, err)
	require.True(t, len(allModels) >= 2)

	// ACTION: Query models specifically for the first vendor
	vendorModels, err := model.GetVendorModels(&suite, vid)
	require.NoError(t, err)
	require.Len(t, vendorModels.Products, 1)
	require.Equal(t, pid, vendorModels.Products[0].Pid)

	// ACTION: Update several fields of the first model
	description := "New Device Description"
	updateMsg := model.NewMsgUpdateModel(vid, pid, vendorAccount.Address)
	updateMsg.ProductLabel = description
	updateMsg.CommissioningModeInitialStepsHint = 8
	updateMsg.CommissioningModeSecondaryStepsHint = 9
	updateMsg.IcdUserActiveModeTriggerHint = 7
	updateMsg.FactoryResetStepsHint = 6
	updateMsg.EnhancedSetupFlowOptions = 2
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Update advanced fields of the second model
	newTCUrl := "https://example.org/file2.txt"
	updateMsg2 := model.NewMsgUpdateModel(vidWithPids, pid, vendorWithPidsAccount.Address)
	updateMsg2.ProductLabel = description
	updateMsg2.EnhancedSetupFlowTCUrl = newTCUrl
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg2}, vendorWithPidsName, vendorWithPidsAccount)
	require.NoError(t, err)

	// ACTION: Verify fields of the first model after update
	receivedModel, err = model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, description, receivedModel.ProductLabel)
	require.Equal(t, uint32(8), receivedModel.CommissioningModeInitialStepsHint)

	// ACTION: Verify fields of the second model after update
	receivedModel2, err = model.GetModel(&suite, vidWithPids, pid)
	require.NoError(t, err)
	require.Equal(t, newTCUrl, receivedModel2.EnhancedSetupFlowTCUrl)

	// ACTION: Update the support URL of the first model
	supportURL := "https://newsupporturl.test"
	updateMsg3 := model.NewMsgUpdateModel(vid, pid, vendorAccount.Address)
	updateMsg3.SupportUrl = supportURL
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg3}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify the support URL was updated correctly
	receivedModel, err = model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, supportURL, receivedModel.SupportUrl)

	// ACTION: Delete the first model
	deleteMsg := model.NewMsgDeleteModel(vid, pid, vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{deleteMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Delete the second model
	deleteMsg2 := model.NewMsgDeleteModel(vidWithPids, pid, vendorWithPidsAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{deleteMsg2}, vendorWithPidsName, vendorWithPidsAccount)
	require.NoError(t, err)

	// ACTION: Verify the first model is no longer found after deletion
	_, err = model.GetModel(&suite, vid, pid)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Verify the second model is no longer found after deletion
	_, err = model.GetModel(&suite, vidWithPids, pid)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Verify model versions for the first model are gone after model deletion
	_, err = model.GetModelVersion(&suite, vid, pid, sv)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Verify model versions for the second model are gone after model deletion
	_, err = model.GetModelVersion(&suite, vidWithPids, pid, sv)
	require.Error(t, err)
	suite.AssertNotFound(err)
}
