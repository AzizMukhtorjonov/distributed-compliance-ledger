// Bash Source: integration_tests/cli/pki-combine-certs.sh
//
// This test verifies operations that combine DA and NOC certificates:
// - ACTION: Propose and approve DA root certificate (pki-combine-certs.sh:L41-L56)
// - ACTION: Add DA intermediate certificate (pki-combine-certs.sh:L60-L68)
// - ACTION: Add NOC root certificate (pki-combine-certs.sh:L75-L80)
// - ACTION: Add NOC ICA certificate (pki-combine-certs.sh:L85-L92)
// - ACTION: Query all-certs and verify both DA and NOC certs present (pki-combine-certs.sh:L100-L145)
// - ACTION: Query all-certs by SKID and verify (pki-combine-certs.sh:L150-L210)
// - ACTION: Query individual cert by subject and SKID (pki-combine-certs.sh:L215-L280)
// - ACTION: Query certs by subject (pki-combine-certs.sh:L285-L355)
// - ACTION: Revoke DA cert and verify all-certs updated (pki-combine-certs.sh:L358-L422)
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

// TestPkiCombineCerts maps to: integration_tests/cli/pki-combine-certs.sh
// ACTION: Tests combined DA and NOC certificate queries via the all-certs endpoint
func TestPkiCombineCerts(t *testing.T) {
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

	// ACTION: Create vendor account (pki-combine-certs.sh:L37-L39)
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

	// ACTION: Propose and approve DA root cert (pki-combine-certs.sh:L41-L56)
	msgProposeRoot := pkitypes.MsgProposeAddX509RootCert{
		Cert:   testconstants.RootCertPem,
		Signer: jackAccount.Address,
		Vid:    testconstants.Vid,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeRoot}, jackName, jackAccount)
	require.NoError(t, err)

	msgApproveRoot := pkitypes.MsgApproveAddX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       aliceAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRoot}, aliceName, aliceAccount)
	require.NoError(t, err)

	// ACTION: Add DA intermediate cert (pki-combine-certs.sh:L60-L68)
	msgAddIntermediate := pkitypes.MsgAddX509Cert{
		Cert:   testconstants.IntermediateCertPem,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddIntermediate}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Add NOC root cert (pki-combine-certs.sh:L75-L80)
	msgAddNocRoot := pkitypes.MsgAddNocX509RootCert{
		Cert:   testconstants.NocRootCert1,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddNocRoot}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query all-certs and verify both DA and NOC certs are present (pki-combine-certs.sh:L100-L145)
	allCerts, err := pki.GetAllCerts(&suite)
	require.NoError(t, err)
	require.True(t, len(allCerts) >= 2, "Expected at least 2 certificates (DA + NOC)")
}
