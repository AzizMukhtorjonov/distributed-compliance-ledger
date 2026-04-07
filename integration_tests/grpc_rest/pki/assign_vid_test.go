// Bash Source: integration_tests/cli/pki-assign-vid.sh
package pki_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	testDclauth "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
	pkitypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/pki/types"
)

func TestPKIAssignVid(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation of Actors
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

	// ACTION: Create a VendorAdmin account to be used for authorized transactions
	vendorAdminName := utils.RandString()
	vendorAdminAccount := testDclauth.CreateAccount(
		&suite,
		vendorAdminName,
		dclauthtypes.AccountRoles{dclauthtypes.VendorAdmin},
		0,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		jackName,
		jackAccount,
		testconstants.Info,
	)
	require.NotNil(t, vendorAdminAccount)

	// ASSIGN VID TO ROOT CERTIFICATE THAT ALREADY HAS VID
	// Note: In the bash script, this section proposes, approves and then assigns VID.

	// ACTION: Propose a new X.509 root certificate to the ledger
	msgProposeAddX509RootCert := pkitypes.MsgProposeAddX509RootCert{
		Cert:   testconstants.RootCertPem,
		Signer: jackAccount.Address,
		Vid:    testconstants.Vid,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeAddX509RootCert}, jackName, jackAccount)
	require.NoError(t, err)

	// ACTION: Approve the proposed X.509 root certificate by a second trustee
	msgApproveAddX509RootCert := pkitypes.MsgApproveAddX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       aliceAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveAddX509RootCert}, aliceName, aliceAccount)
	require.NoError(t, err)

	// ACTION: Assign a Vendor ID (VID) to the approved root certificate
	msgAssignVid := pkitypes.MsgAssignVid{
		Signer:       vendorAdminAccount.Address,
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Vid:          testconstants.Vid,
	}
	res, err := suite.BuildAndBroadcastTx([]sdk.Msg{&msgAssignVid}, vendorAdminName, vendorAdminAccount)
	// Bash script checks: check_response "$result" "vid is not empty"
	// In the ledger implementation, if the certificate already has a VID (propose with VID),
	// assign-vid would fail with "vid is not empty".
	require.Error(t, err)
	require.Contains(t, err.Error(), "vid is not empty")
	require.NotNil(t, res)
}
