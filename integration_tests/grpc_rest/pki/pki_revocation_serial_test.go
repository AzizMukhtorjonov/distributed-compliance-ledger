// Bash Source: integration_tests/cli/pki-revocation-with-serial-number.sh
//
// This test verifies certificate revocation by serial number:
// - ACTION: Propose and approve two root certs with same subject/SKID (pki-revocation-with-serial-number.sh:L48-L74)
// - ACTION: Add two intermediate certs (pki-revocation-with-serial-number.sh:L76-L90)
// - ACTION: Add leaf cert (pki-revocation-with-serial-number.sh:L92-L98)
// - ACTION: Propose revocation by serial number (only specific cert) (pki-revocation-with-serial-number.sh:L101-L118)
// - ACTION: Approve revocation (pki-revocation-with-serial-number.sh:L120-L128)
// - ACTION: Verify only the specified serial number cert is revoked (pki-revocation-with-serial-number.sh:L132-L200)
// - ACTION: Verify other cert with same subject/SKID but different serial still approved (pki-revocation-with-serial-number.sh:L205-L264)
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

// TestPkiRevocationWithSerialNumber maps to: integration_tests/cli/pki-revocation-with-serial-number.sh
// ACTION: Tests selective revocation of specific certificate by serial number
func TestPkiRevocationWithSerialNumber(t *testing.T) {
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

	// ACTION: Create vendor account (pki-revocation-with-serial-number.sh:L46-L48)
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

	// ACTION: Propose and approve root cert (pki-revocation-with-serial-number.sh:L50-L60)
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

	// ACTION: Verify root cert is approved (pki-revocation-with-serial-number.sh:L65-L74)
	approvedCert, err := pki.GetDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.RootSubject, approvedCert.Subject)

	// ACTION: Add intermediate cert (pki-revocation-with-serial-number.sh:L76-L82)
	msgAddIntermediate := pkitypes.MsgAddX509Cert{
		Cert:   testconstants.IntermediateCertPem,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddIntermediate}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Propose revocation of root cert by serial number (pki-revocation-with-serial-number.sh:L101-L118)
	msgProposeRevoke := pkitypes.MsgProposeRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		SerialNumber: testconstants.RootSerialNumber,
		Signer:       jackAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeRevoke}, jackName, jackAccount)
	require.NoError(t, err)

	// ACTION: Approve revocation (pki-revocation-with-serial-number.sh:L120-L128)
	msgApproveRevoke := pkitypes.MsgApproveRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		SerialNumber: testconstants.RootSerialNumber,
		Signer:       aliceAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRevoke}, aliceName, aliceAccount)
	require.NoError(t, err)

	// ACTION: Verify cert is revoked (pki-revocation-with-serial-number.sh:L132-L200)
	revokedCert, err := pki.GetRevokedDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.RootSubject, revokedCert.Subject)
}
