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
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
)

func TestModelUpdate(t *testing.T) {
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

	// ACTION: Query a model that does not exist and expect "Not Found"
	_, err = model.GetModel(&suite, vid, pid)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Add a new Model record
	productLabel := "Device #1"
	createMsg := model.NewMsgCreateModel(vid, pid, vendorAccount.Address)
	createMsg.ProductLabel = productLabel
	createMsg.SchemaVersion = 0
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query the newly created Model and verify default hint values
	receivedModel, err := model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, uint32(1), receivedModel.CommissioningModeInitialStepsHint)
	require.Equal(t, uint32(4), receivedModel.CommissioningModeSecondaryStepsHint)
	require.Equal(t, uint32(1), receivedModel.IcdUserActiveModeTriggerHint)
	require.Equal(t, uint32(1), receivedModel.FactoryResetStepsHint)

	// ACTION: Update the model with custom hint values
	description := "New Device Description"
	updateMsg := model.NewMsgUpdateModel(vid, pid, vendorAccount.Address)
	updateMsg.ProductLabel = description
	updateMsg.CommissioningModeInitialStepsHint = 3
	updateMsg.CommissioningModeSecondaryStepsHint = 7
	updateMsg.IcdUserActiveModeTriggerHint = 6
	updateMsg.FactoryResetStepsHint = 5
	updateMsg.EnhancedSetupFlowOptions = 2
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify fields reflect the updated hint values
	receivedModel, err = model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, uint32(3), receivedModel.CommissioningModeInitialStepsHint)
	require.Equal(t, uint32(7), receivedModel.CommissioningModeSecondaryStepsHint)
	require.Equal(t, uint32(6), receivedModel.IcdUserActiveModeTriggerHint)
	require.Equal(t, uint32(5), receivedModel.FactoryResetStepsHint)

	// ACTION: Update only the description of the model
	description2 := "New Device Description 2"
	updateMsg2 := model.NewMsgUpdateModel(vid, pid, vendorAccount.Address)
	updateMsg2.ProductLabel = description2
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg2}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify hints remained unchanged after description-only update
	receivedModel, err = model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, uint32(3), receivedModel.CommissioningModeInitialStepsHint)

	// ACTION: Attempt to update hints with 0 (expecting them to remain unchanged or use previous values)
	description3 := "New Device Description 3"
	updateMsg3 := model.NewMsgUpdateModel(vid, pid, vendorAccount.Address)
	updateMsg3.ProductLabel = description3
	updateMsg3.CommissioningModeInitialStepsHint = 0
	updateMsg3.CommissioningModeSecondaryStepsHint = 0
	updateMsg3.IcdUserActiveModeTriggerHint = 0
	updateMsg3.FactoryResetStepsHint = 0
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg3}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify hints remained unchanged after setting to 0 in update
	receivedModel, err = model.GetModel(&suite, vid, pid)
	require.NoError(t, err)
	require.Equal(t, uint32(3), receivedModel.CommissioningModeInitialStepsHint)
	require.Equal(t, description3, receivedModel.ProductLabel)
}
