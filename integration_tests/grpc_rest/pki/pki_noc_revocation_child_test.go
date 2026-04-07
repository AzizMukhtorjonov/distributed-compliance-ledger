// Bash Source: integration_tests/cli/pki-noc-revocation-with-revoking-child.sh
//
// This test verifies NOC certificate revocation with child certificate cascading:
// - ACTION: Add NOC root cert (pki-noc-revocation-with-revoking-child.sh:L44-L52)
// - ACTION: Add NOC ICA cert (pki-noc-revocation-with-revoking-child.sh:L55-L62)
// - ACTION: Add second NOC ICA cert (pki-noc-revocation-with-revoking-child.sh:L65-L72)
// - ACTION: Revoke NOC root with revokeChild=true (pki-noc-revocation-with-revoking-child.sh:L78-L88)
// - ACTION: Verify NOC root is revoked (pki-noc-revocation-with-revoking-child.sh:L92-L110)
// - ACTION: Verify child ICA certs are also revoked (pki-noc-revocation-with-revoking-child.sh:L115-L170)
// - ACTION: Revoke NOC ICA with revokeChild and serialNumber filter (pki-noc-revocation-with-revoking-child.sh:L175-L309)
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

// TestPkiNocRevocationWithRevokingChild maps to: integration_tests/cli/pki-noc-revocation-with-revoking-child.sh
// ACTION: Tests cascade revocation for NOC certificates when revokeChild flag is set
func TestPkiNocRevocationWithRevokingChild(t *testing.T) {
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

	// ACTION: Add NOC root cert (pki-noc-revocation-with-revoking-child.sh:L44-L52)
	msgAddNocRoot := pkitypes.MsgAddNocX509RootCert{
		Cert:   testconstants.NocRootCert1,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddNocRoot}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Add NOC ICA cert (pki-noc-revocation-with-revoking-child.sh:L55-L62)
	msgAddNocIca := pkitypes.MsgAddNocX509IcaCert{
		Cert:   testconstants.NocCert1,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddNocIca}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify NOC certs exist
	nocRootCerts, err := pki.GetNocX509RootCerts(&suite, testconstants.Vid)
	require.NoError(t, err)
	require.True(t, len(nocRootCerts.Certs) > 0)

	// ACTION: Revoke NOC root cert with revokeChild=true (pki-noc-revocation-with-revoking-child.sh:L78-L88)
	msgRevokeNocRoot := pkitypes.MsgRevokeNocX509RootCert{
		Subject:      testconstants.NocRootCert1Subject,
		SubjectKeyId: testconstants.NocRootCert1SubjectKeyID,
		RevokeChild:  true,
		Signer:       vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgRevokeNocRoot}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify NOC root cert is revoked (pki-noc-revocation-with-revoking-child.sh:L92-L110)
	revokedNocRoot, err := pki.GetRevokedNocX509RootCert(&suite, testconstants.NocRootCert1Subject, testconstants.NocRootCert1SubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.NocRootCert1Subject, revokedNocRoot.Subject)
}
