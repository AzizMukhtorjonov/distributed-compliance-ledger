package vendorinfo_test

import (
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	testDclauth "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/vendorinfo"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
	vendorinfotypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/vendorinfo/types"
)

func TestVendorInfoDemo(t *testing.T) {
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
	vid2 := int32(tmrand.Uint16())

	// ACTION: Create a vendor account with a random Vendor ID
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

	// ACTION: Create a second vendor account with another random Vendor ID
	secondVendorName := utils.RandString()
	secondVendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		secondVendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		vid2,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	require.NotNil(t, secondVendorAccount)

	// ACTION: Create a VendorAdmin account for administrative operations
	vendorAdminName := utils.RandString()
	vendorAdminAccount := testDclauth.CreateAccount(
		&suite,
		vendorAdminName,
		dclauthtypes.AccountRoles{dclauthtypes.VendorAdmin},
		0,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)
	require.NotNil(t, vendorAdminAccount)

	// ACTION: Query VendorInfo for a non-existent VID and expect it to not be found
	_, err = vendorinfo.GetVendorInfo(&suite, vid)
	require.Error(t, err)
	suite.AssertNotFound(err)

	// ACTION: Query all vendors and expect an empty list initially
	allVendors, err := vendorinfo.GetVendorInfos(&suite)
	require.NoError(t, err)
	require.Empty(t, allVendors)

	// ACTION: Add VendorInfo record for the first VID using its vendor account
	companyLegalName := "XYZ IOT Devices Inc"
	vendorNameValue := "XYZ Devices"
	createMsg := &vendorinfotypes.MsgCreateVendorInfo{
		Creator:          vendorAccount.Address,
		VendorID:         vid,
		CompanyLegalName: companyLegalName,
		VendorName:       vendorNameValue,
		SchemaVersion:    0,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{createMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query the newly created VendorInfo and verify its fields
	receivedVendorInfo, err := vendorinfo.GetVendorInfo(&suite, vid)
	require.NoError(t, err)
	require.Equal(t, vid, receivedVendorInfo.VendorID)
	require.Equal(t, companyLegalName, receivedVendorInfo.CompanyLegalName)
	require.Equal(t, vendorNameValue, receivedVendorInfo.VendorName)

	// ACTION: Query all vendors and verify the created record is present
	allVendors, err = vendorinfo.GetVendorInfos(&suite)
	require.NoError(t, err)
	require.Len(t, allVendors, 1)
	require.Equal(t, vid, allVendors[0].VendorID)

	// ACTION: Update VendorInfo record with only required fields (omitting optional ones)
	updateMsg := &vendorinfotypes.MsgUpdateVendorInfo{
		Creator:  vendorAccount.Address,
		VendorID: vid,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{updateMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify fields after partial update to ensure omitted fields weren't wiped
	receivedVendorInfo, err = vendorinfo.GetVendorInfo(&suite, vid)
	require.NoError(t, err)
	require.Equal(t, companyLegalName, receivedVendorInfo.CompanyLegalName)
	require.Equal(t, vendorNameValue, receivedVendorInfo.VendorName)

	// ACTION: Perform a full update of the VendorInfo record including optional fields
	companyLegalName = "ABC Subsidiary Corporation"
	vendorLandingPageURL := "https://www.w3.org/"
	fullUpdateMsg := &vendorinfotypes.MsgUpdateVendorInfo{
		Creator:              vendorAccount.Address,
		VendorID:             vid,
		CompanyLegalName:     companyLegalName,
		VendorLandingPageURL: vendorLandingPageURL,
		VendorName:           vendorNameValue,
		SchemaVersion:        0,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{fullUpdateMsg}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify the fields reflect the full update
	receivedVendorInfo, err = vendorinfo.GetVendorInfo(&suite, vid)
	require.NoError(t, err)
	require.Equal(t, companyLegalName, receivedVendorInfo.CompanyLegalName)
	require.Equal(t, vendorLandingPageURL, receivedVendorInfo.VendorLandingPageURL)

	// ACTION: Attempt to add VendorInfo for a different VID and expect an authorization error
	vid1 := int32(tmrand.Uint16())
	badCreateMsg := &vendorinfotypes.MsgCreateVendorInfo{
		Creator:          vendorAccount.Address,
		VendorID:         vid1,
		CompanyLegalName: companyLegalName,
		VendorName:       vendorNameValue,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badCreateMsg}, vendorName, vendorAccount)
	require.Error(t, err)
	require.True(t, sdkerrors.ErrUnauthorized.Is(err))

	// ACTION: Attempt to update VendorInfo using an account from another vendor and expect an error
	badUpdateMsg := &vendorinfotypes.MsgUpdateVendorInfo{
		Creator:          secondVendorAccount.Address,
		VendorID:         vid,
		CompanyLegalName: companyLegalName,
		VendorName:       vendorNameValue,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{badUpdateMsg}, secondVendorName, secondVendorAccount)
	require.Error(t, err)
	require.True(t, sdkerrors.ErrUnauthorized.Is(err))

	// ACTION: Create a VendorInfo record using a VendorAdmin account (allowed)
	vid3 := int32(tmrand.Uint16())
	adminCreateMsg := &vendorinfotypes.MsgCreateVendorInfo{
		Creator:          vendorAdminAccount.Address,
		VendorID:         vid3,
		CompanyLegalName: companyLegalName,
		VendorName:       vendorNameValue,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{adminCreateMsg}, vendorAdminName, vendorAdminAccount)
	require.NoError(t, err)

	// ACTION: Update the VendorInfo record using a VendorAdmin account
	companyLegalName1 := "New Corp"
	vendorName1 := "New Vendor Name"
	adminUpdateMsg := &vendorinfotypes.MsgUpdateVendorInfo{
		Creator:          vendorAdminAccount.Address,
		VendorID:         vid3,
		CompanyLegalName: companyLegalName1,
		VendorName:       vendorName1,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{adminUpdateMsg}, vendorAdminName, vendorAdminAccount)
	require.NoError(t, err)
}
