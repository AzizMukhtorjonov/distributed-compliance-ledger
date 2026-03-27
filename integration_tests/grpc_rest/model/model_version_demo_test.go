package model_test

import (
	"context"
	"fmt"
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
	modeltypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/model/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestModelVersionDemo(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation of Actors
	aliceName := testconstants.AliceAccount
	aliceKeyInfo, err := suite.Kr.Key(aliceName)
	require.NoError(t, err)
	address, err := aliceKeyInfo.GetAddress()
	require.NoError(t, err)
	aliceAccount, err := testDclauth.GetAccount(&suite, address)
	require.NoError(t, err)

	bobName := testconstants.BobAccount
	bobKeyInfo, err := suite.Kr.Key(bobName)
	require.NoError(t, err)
	address, err = bobKeyInfo.GetAddress()
	require.NoError(t, err)
	bobAccount, err := testDclauth.GetAccount(&suite, address)
	require.NoError(t, err)

	// Create Vendor account
	vid := int32(tmrand.Uint16())
	pid := int32(tmrand.Uint16())
	vendorName := utils.RandString()
	vendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		vendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		vid,
		[]*commontypes.Uint16Range{{Min: pid, Max: pid}},
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	require.NotNil(t, vendorAccount)

	// Create model (prerequisite)
	createModelMsg := model.NewMsgCreateModel(vid, pid, vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	sv := uint32(tmrand.Uint16())
	createModelVersionMsg := model.NewMsgCreateModelVersion(vid, pid, sv, "1", vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelVersionMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// Query the model version
	receivedModelVersion, err := model.GetModelVersion(&suite, vid, pid, sv)
	require.NoError(t, err)
	require.Equal(t, vid, receivedModelVersion.Vid)
	require.Equal(t, pid, receivedModelVersion.Pid)
	require.Equal(t, sv, receivedModelVersion.SoftwareVersion)
	require.Equal(t, "1", receivedModelVersion.SoftwareVersionString)
	require.Equal(t, uint32(testconstants.CdVersionNumber), receivedModelVersion.CdVersionNumber)
	require.True(t, receivedModelVersion.SoftwareVersionValid)
	require.Equal(t, uint32(testconstants.MinApplicableSoftwareVersion), receivedModelVersion.MinApplicableSoftwareVersion)
	require.Equal(t, uint32(testconstants.MaxApplicableSoftwareVersion), receivedModelVersion.MaxApplicableSoftwareVersion)

	// Query all model versions
	modelVersions, err := GetModelVersions(&suite, vid, pid)
	require.NoError(t, err)
	require.Len(t, modelVersions.SoftwareVersions, 1)
	require.Equal(t, sv, modelVersions.SoftwareVersions[0])

	// Query non existent model version
	_, err = model.GetModelVersion(&suite, vid, pid, 123456)
	require.Error(t, err)
	stat, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())

	// Query non existent model versions
	_, err = GetModelVersions(&suite, int32(tmrand.Uint16()), int32(tmrand.Uint16()))
	require.Error(t, err)
	stat, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())

	// Update the existing model version
	updateModelVersionMsg := model.NewMsgUpdateModelVersion(vid, pid, sv, vendorAccount.Address)
	// Customize update message as per bash script: minApplicable=2, max=10, valid=false
	updateModelVersionMsg.MinApplicableSoftwareVersion = 2
	updateModelVersionMsg.MaxApplicableSoftwareVersion = 10
	updateModelVersionMsg.SoftwareVersionValid = false

	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateModelVersionMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// Query Updated model version
	receivedModelVersion, err = model.GetModelVersion(&suite, vid, pid, sv)
	require.NoError(t, err)
	require.False(t, receivedModelVersion.SoftwareVersionValid)
	require.Equal(t, uint32(2), receivedModelVersion.MinApplicableSoftwareVersion)
	require.Equal(t, uint32(10), receivedModelVersion.MaxApplicableSoftwareVersion)

	// Add second model version
	sv2 := uint32(tmrand.Uint16())
	// Ensure sv2 is different from sv
	for sv2 == sv {
		sv2 = uint32(tmrand.Uint16())
	}
	createModelVersionMsg2 := model.NewMsgCreateModelVersion(vid, pid, sv2, "2", vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelVersionMsg2}, vendorName, vendorAccount)
	require.NoError(t, err)

	// Query all model versions
	modelVersions, err = GetModelVersions(&suite, vid, pid)
	require.NoError(t, err)
	require.Len(t, modelVersions.SoftwareVersions, 2)
	require.Contains(t, modelVersions.SoftwareVersions, sv)
	require.Contains(t, modelVersions.SoftwareVersions, sv2)

	// Create model version with vid belonging to another vendor
	newVid := int32(tmrand.Uint16())
	differentVendorName := utils.RandString()
	differentVendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		differentVendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		newVid,
		[]*commontypes.Uint16Range{}, // No products needed for account creation test
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	// Try to create version for ORIGINAL vid/pid using NEW vendor
	createModelVersionMsg3 := model.NewMsgCreateModelVersion(vid, pid, sv, "1", differentVendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createModelVersionMsg3}, differentVendorName, differentVendorAccount)
	require.Error(t, err)
	// Bash script checks: "transaction should be signed by a vendor account containing the vendorID $vid"
	// We expect an unauthorized error or similar validation error.

	// Update model version with vid belonging to another vendor
	updateModelVersionMsg2 := model.NewMsgUpdateModelVersion(vid, pid, sv, differentVendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateModelVersionMsg2}, differentVendorName, differentVendorAccount)
	require.Error(t, err)

	// Delete existing model version
	deleteModelVersionMsg := model.NewMsgDeleteModelVersion(vid, pid, sv, vendorAccount.Address)
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{deleteModelVersionMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// Query deleted model version
	_, err = model.GetModelVersion(&suite, vid, pid, sv)
	require.Error(t, err)
	stat, ok = status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, stat.Code())
}

func GetModelVersions(
	suite *utils.TestSuite,
	vid int32,
	pid int32,
) (*modeltypes.ModelVersions, error) {
	var res modeltypes.ModelVersions

	if suite.Rest {
		var resp modeltypes.QueryGetModelVersionsResponse
		err := suite.QueryREST(fmt.Sprintf("/dcl/model/versions/%v/%v", vid, pid), &resp)
		if err != nil {
			return nil, err
		}
		res = resp.GetModelVersions()
	} else {
		grpcConn := suite.GetGRPCConn()
		defer grpcConn.Close()

		// This creates a gRPC client to query the x/dclauth service.
		modelClient := modeltypes.NewQueryClient(grpcConn)
		resp, err := modelClient.ModelVersions(
			context.Background(),
			&modeltypes.QueryGetModelVersionsRequest{Vid: vid, Pid: pid},
		)
		if err != nil {
			return nil, err
		}
		res = resp.GetModelVersions()
	}

	return &res, nil
}
