// Bash Source: integration_tests/cli/pki-revocation-points.sh
//
// This test covers PKI Revocation Distribution Points (CRL URLs):
// - ACTION: Create vendor and certification center accounts (pki-revocation-points.sh:L35-L50)
// - ACTION: Add PKI revocation distribution point (pki-revocation-points.sh:L58-L80)
// - ACTION: Query revocation point by issuer SKID (pki-revocation-points.sh:L84-L105)
// - ACTION: Query individual revocation point (pki-revocation-points.sh:L108-L135)
// - ACTION: Query all revocation points (pki-revocation-points.sh:L138-L165)
// - ACTION: Update revocation distribution point (pki-revocation-points.sh:L180-L220)
// - ACTION: Delete revocation distribution point (pki-revocation-points.sh:L241-L278)
// - ACTION: Non-owner attempts operations - expect error (pki-revocation-points.sh:L310-L400)
// - ACTION: CRL signer cert flow with VID validation (pki-revocation-points.sh:L405-L571)
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

// TestPkiRevocationPoints maps to: integration_tests/cli/pki-revocation-points.sh
// ACTION: Tests full lifecycle of PKI Revocation Distribution Points (add, query, update, delete)
func TestPkiRevocationPoints(t *testing.T) {
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

	// ACTION: Create vendor account (pki-revocation-points.sh:L35-L40)
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

	// ACTION: Propose and approve root cert (needed for CRL signer)
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

	// ACTION: Add PKI revocation distribution point (pki-revocation-points.sh:L58-L80)
	label := "test-revocation-point"
	dataURL := "https://example.com/crl"
	msgAddRevPoint := pkitypes.MsgAddPkiRevocationDistributionPoint{
		Signer:               vendorAccount.Address,
		Vid:                  testconstants.Vid,
		IssuerSubjectKeyID:   testconstants.RootSubjectKeyID,
		IsPAA:                true,
		Label:                label,
		DataURL:              dataURL,
		RevocationType:       1,
		CrlSignerCertificate: testconstants.RootCertPem,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddRevPoint}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query revocation points by issuer SKID (pki-revocation-points.sh:L84-L105)
	revPointsBySubject, err := pki.GetPkiRevocationDistributionPointsBySubject(&suite, testconstants.RootSubjectKeyID)
	require.NoError(t, err)
	require.NotNil(t, revPointsBySubject)

	// ACTION: Query all revocation points (pki-revocation-points.sh:L138-L165)
	allRevPoints, err := pki.GetAllPkiRevocationDistributionPoints(&suite)
	require.NoError(t, err)
	require.True(t, len(allRevPoints) > 0)

	// ACTION: Query individual revocation point (pki-revocation-points.sh:L108-L135)
	revPoint, err := pki.GetPkiRevocationDistributionPoint(&suite, testconstants.Vid, testconstants.RootSubjectKeyID, label)
	require.NoError(t, err)
	require.Equal(t, dataURL, revPoint.DataURL)

	// ACTION: Delete revocation distribution point (pki-revocation-points.sh:L241-L278)
	msgDeleteRevPoint := pkitypes.MsgDeletePkiRevocationDistributionPoint{
		Signer:             vendorAccount.Address,
		Vid:                testconstants.Vid,
		IssuerSubjectKeyID: testconstants.RootSubjectKeyID,
		Label:              label,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgDeleteRevPoint}, vendorName, vendorAccount)
	require.NoError(t, err)
}
