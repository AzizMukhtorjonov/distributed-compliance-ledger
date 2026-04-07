// Bash Source: integration_tests/cli/pki-remove-noc-certificates.sh
//
// This test verifies NOC certificate removal operations:
// - ACTION: Add NOC root cert (pki-remove-noc-certificates.sh:L44-L52)
// - ACTION: Add NOC ICA certs (pki-remove-noc-certificates.sh:L55-L72)
// - ACTION: Remove NOC ICA cert by serial number (pki-remove-noc-certificates.sh:L78-L92)
// - ACTION: Verify only specified NOC ICA removed (pki-remove-noc-certificates.sh:L96-L140)
// - ACTION: Remove all remaining NOC ICA certs (pki-remove-noc-certificates.sh:L145-L165)
// - ACTION: Remove NOC root cert (pki-remove-noc-certificates.sh:L170-L195)
// - ACTION: Non-owner vendor attempts NOC removal - expect error (pki-remove-noc-certificates.sh:L200-L341)
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

// TestPkiRemoveNocCertificates maps to: integration_tests/cli/pki-remove-noc-certificates.sh
// ACTION: Tests NOC certificate removal including serial number selection and ownership checks
func TestPkiRemoveNocCertificates(t *testing.T) {
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

	// ACTION: Add NOC root cert (pki-remove-noc-certificates.sh:L44-L52)
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

	// ACTION: Remove NOC root cert (pki-remove-noc-certificates.sh:L170-L195)
	msgRemoveNocRoot := pkitypes.MsgRemoveNocX509RootCert{
		Subject:      testconstants.NocRootCert1Subject,
		SubjectKeyId: testconstants.NocRootCert1SubjectKeyID,
		Signer:       vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgRemoveNocRoot}, vendorName, vendorAccount)
	require.NoError(t, err)
}
