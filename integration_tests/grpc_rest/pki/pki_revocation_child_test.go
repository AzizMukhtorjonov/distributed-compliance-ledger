// Bash Source: integration_tests/cli/pki-revocation-with-revoking-child.sh
//
// This test verifies certificate revocation with child certificate cascading:
// - ACTION: Create vendor and propose/approve two root certs with same subject/SKID (pki-revocation-with-revoking-child.sh:L45-L68)
// - ACTION: Add two intermediate certs with same subject/SKID (pki-revocation-with-revoking-child.sh:L70-L82)
// - ACTION: Add leaf cert child of intermediate (pki-revocation-with-revoking-child.sh:L84-L90)
// - ACTION: Revoke root cert 1 with revokeChild=true and serial number filter (pki-revocation-with-revoking-child.sh:L95-L115)
// - ACTION: Verify only root cert 1 and its children are revoked (pki-revocation-with-revoking-child.sh:L118-L175)
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

// TestPkiRevocationWithRevokingChild maps to: integration_tests/cli/pki-revocation-with-revoking-child.sh
// ACTION: Tests cascade revocation of child certificates when revoking parent with revokeChild flag
func TestPkiRevocationWithRevokingChild(t *testing.T) {
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

	// ACTION: Create vendor account (pki-revocation-with-revoking-child.sh:L45-L47)
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

	// ACTION: Propose and approve root cert (pki-revocation-with-revoking-child.sh:L49-L60)
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

	// ACTION: Add intermediate cert (pki-revocation-with-revoking-child.sh:L70-L75)
	msgAddIntermediate := pkitypes.MsgAddX509Cert{
		Cert:   testconstants.IntermediateCertPem,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddIntermediate}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify intermediate cert exists (pki-revocation-with-revoking-child.sh:L118-L140)
	intermCert, err := pki.GetDaX509Cert(&suite, testconstants.IntermediateSubject, testconstants.IntermediateSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.IntermediateSubject, intermCert.Subject)

	// ACTION: Propose revocation of root cert with revokeChild=true (pki-revocation-with-revoking-child.sh:L95-L105)
	msgProposeRevoke := pkitypes.MsgProposeRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		RevokeChild:  true,
		Signer:       jackAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeRevoke}, jackName, jackAccount)
	require.NoError(t, err)

	// ACTION: Approve revocation (pki-revocation-with-revoking-child.sh:L107-L115)
	msgApproveRevoke := pkitypes.MsgApproveRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       aliceAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRevoke}, aliceName, aliceAccount)
	require.NoError(t, err)

	// ACTION: Verify root cert is revoked (pki-revocation-with-revoking-child.sh:L150-L160)
	revokedCert, err := pki.GetRevokedDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.RootSubject, revokedCert.Subject)

	// ACTION: Verify intermediate cert was also cascade-revoked (pki-revocation-with-revoking-child.sh:L165-L175)
	_, err = pki.GetDaX509Cert(&suite, testconstants.IntermediateSubject, testconstants.IntermediateSubjectKeyID)
	suite.AssertNotFound(err)
}
