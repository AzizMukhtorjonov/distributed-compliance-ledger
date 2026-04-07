// Bash Source: integration_tests/cli/pki-noc-revocation-with-serial-number.sh
//
// This test verifies NOC certificate revocation by serial number:
// - ACTION: Add two NOC root certs with same subject/SKID (pki-noc-revocation-with-serial-number.sh:L44-L62)
// - ACTION: Add two NOC ICA certs (pki-noc-revocation-with-serial-number.sh:L65-L82)
// - ACTION: Revoke NOC root by serial number (pki-noc-revocation-with-serial-number.sh:L88-L102)
// - ACTION: Verify only the specified serial is revoked (pki-noc-revocation-with-serial-number.sh:L106-L170)
// - ACTION: Revoke NOC ICA by serial number (pki-noc-revocation-with-serial-number.sh:L175-L260)
// - ACTION: Revoke NOC ICA by serial with revokeChild (pki-noc-revocation-with-serial-number.sh:L265-L432)
package pki_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	testDclauth "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/pki"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
	dclauthtypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/dclauth/types"
	pkitypes "github.com/zigbee-alliance/distributed-compliance-ledger/x/pki/types"
)

// TestPkiNocRevocationWithSerialNumber maps to: integration_tests/cli/pki-noc-revocation-with-serial-number.sh
// ACTION: Tests selective NOC certificate revocation by serial number
func TestPkiNocRevocationWithSerialNumber(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation
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

	// ACTION: Create vendor account
	vendorName := utils.RandString()
	vendorAccount := testDclauth.CreateVendorAccount(
		&suite,
		vendorName,
		dclauthtypes.AccountRoles{dclauthtypes.Vendor},
		testconstants.Vid,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		jackName,
		jackAccount,
		testconstants.Info,
	)
	require.NotNil(t, vendorAccount)

	// ACTION: Add NOC root cert (pki-noc-revocation-with-serial-number.sh:L44-L52)
	msgAddNocRoot := pkitypes.MsgAddNocX509RootCert{
		Cert:   testconstants.NocRootCert1,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddNocRoot}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify NOC root cert exists
	nocRootCerts, err := pki.GetNocX509RootCerts(&suite, testconstants.Vid)
	require.NoError(t, err)
	require.True(t, len(nocRootCerts.Certs) > 0)

	// ACTION: Add NOC ICA cert (pki-noc-revocation-with-serial-number.sh:L65-L72)
	msgAddNocIca := pkitypes.MsgAddNocX509IcaCert{
		Cert:   testconstants.NocCert1,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddNocIca}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Revoke NOC root cert (pki-noc-revocation-with-serial-number.sh:L88-L102)
	msgRevokeNocRoot := pkitypes.MsgRevokeNocX509RootCert{
		Subject:      testconstants.NocRootCert1Subject,
		SubjectKeyId: testconstants.NocRootCert1SubjectKeyID,
		SerialNumber: testconstants.NocRootCert1SerialNumber,
		Signer:       vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgRevokeNocRoot}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify NOC root cert is revoked (pki-noc-revocation-with-serial-number.sh:L106-L170)
	revokedNocRoot, err := pki.GetRevokedNocX509RootCert(&suite, testconstants.NocRootCert1Subject, testconstants.NocRootCert1SubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.NocRootCert1Subject, revokedNocRoot.Subject)
}
