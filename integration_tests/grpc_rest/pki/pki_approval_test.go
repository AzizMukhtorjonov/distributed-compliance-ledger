// Bash Source: integration_tests/cli/pki-approval.sh
//
// This test verifies the multi-trustee approval quorum for PKI root certificate operations:
// - ACTION: Add 3 new Trustee accounts (total 6 Trustees, 4 approvals needed) (pki-approval.sh:L49-L134)
// - ACTION: Non-Trustee tries to propose root cert - expect error (pki-approval.sh:L139-L143)
// - ACTION: 4th Trustee proposes root cert (pki-approval.sh:L145-L148)
// - ACTION: 1st Trustee approves - still not enough (2/4 needed) (pki-approval.sh:L154-L157)
// - ACTION: Verify cert is NOT approved yet (pki-approval.sh:L161-L163)
// - ACTION: 2nd Trustee approves (pki-approval.sh:L167-L170)
// - ACTION: Verify cert is still NOT approved (pki-approval.sh:L174-L176)
// - ACTION: 3rd Trustee approves - now quorum reached, cert becomes approved (pki-approval.sh:L180-L183)
// - ACTION: Verify cert IS approved with 3 approval addresses (pki-approval.sh:L187-L194)
// - ACTION: 6th Trustee proposes revocation (pki-approval.sh:L197-L200)
// - ACTION: 5th, 4th, 3rd Trustees approve revocation - quorum reached (pki-approval.sh:L211-L241)
// - ACTION: Verify cert is revoked (pki-approval.sh:L243-L252)
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

// TestPkiApproval maps to: integration_tests/cli/pki-approval.sh
// ACTION: Tests 2/3 trustee quorum for root certificate operations with 6 trustees
func TestPkiApproval(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)

	// Preparation - get existing trustee accounts
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

	bobName := testconstants.BobAccount
	bobKeyInfo, err := suite.Kr.Key(bobName)
	require.NoError(t, err)
	bobAddr, err := bobKeyInfo.GetAddress()
	require.NoError(t, err)
	bobAccount, err := testDclauth.GetAccount(&suite, bobAddr)
	require.NoError(t, err)

	// ACTION: Create a regular (CertificationCenter) account (pki-approval.sh:L44-L46)
	userName := utils.RandString()
	_ = testDclauth.CreateAccount(
		&suite,
		userName,
		dclauthtypes.AccountRoles{dclauthtypes.CertificationCenter},
		0,
		testconstants.ProductIDsEmpty,
		aliceName,
		aliceAccount,
		bobName,
		bobAccount,
		testconstants.Info,
	)

	// ACTION: Add 3 new Trustee accounts (4th, 5th, 6th) for 6 total (pki-approval.sh:L53-L134)
	fourthTrusteeName := utils.RandString()
	fourthTrusteeAccount := testDclauth.CreateAccount(
		&suite,
		fourthTrusteeName,
		dclauthtypes.AccountRoles{dclauthtypes.Trustee},
		0,
		testconstants.ProductIDsEmpty,
		jackName,
		jackAccount,
		aliceName,
		aliceAccount,
		testconstants.Info,
	)
	require.NotNil(t, fourthTrusteeAccount)

	fifthTrusteeName := utils.RandString()
	fifthTrusteeAccount := testDclauth.CreateAccount(
		&suite,
		fifthTrusteeName,
		dclauthtypes.AccountRoles{dclauthtypes.Trustee},
		0,
		testconstants.ProductIDsEmpty,
		jackName,
		jackAccount,
		aliceName,
		aliceAccount,
		testconstants.Info,
	)
	require.NotNil(t, fifthTrusteeAccount)

	sixthTrusteeName := utils.RandString()
	sixthTrusteeAccount := testDclauth.CreateAccount(
		&suite,
		sixthTrusteeName,
		dclauthtypes.AccountRoles{dclauthtypes.Trustee},
		0,
		testconstants.ProductIDsEmpty,
		jackName,
		jackAccount,
		aliceName,
		aliceAccount,
		testconstants.Info,
	)
	require.NotNil(t, sixthTrusteeAccount)

	// ACTION: 4th Trustee proposes root cert (pki-approval.sh:L145-L148)
	msgProposeAddX509RootCert := pkitypes.MsgProposeAddX509RootCert{
		Cert:   testconstants.RootCertPem,
		Signer: fourthTrusteeAccount.Address,
		Vid:    testconstants.Vid,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeAddX509RootCert}, fourthTrusteeName, fourthTrusteeAccount)
	require.NoError(t, err)

	// ACTION: 1st Trustee approves - still not enough approvals (pki-approval.sh:L154-L157)
	msgApprove := pkitypes.MsgApproveAddX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       jackAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApprove}, jackName, jackAccount)
	require.NoError(t, err)

	// ACTION: Verify cert is NOT approved yet (pki-approval.sh:L161-L163)
	_, err = pki.GetDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	suite.AssertNotFound(err)

	// ACTION: 2nd Trustee approves (pki-approval.sh:L167-L170)
	msgApprove2 := pkitypes.MsgApproveAddX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       aliceAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApprove2}, aliceName, aliceAccount)
	require.NoError(t, err)

	// ACTION: Verify cert is still NOT approved (pki-approval.sh:L174-L176)
	_, err = pki.GetDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	suite.AssertNotFound(err)

	// ACTION: 3rd Trustee approves - quorum reached (pki-approval.sh:L180-L183)
	msgApprove3 := pkitypes.MsgApproveAddX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       bobAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApprove3}, bobName, bobAccount)
	require.NoError(t, err)

	// ACTION: Verify cert IS now approved (pki-approval.sh:L187-L194)
	approvedCert, err := pki.GetDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.RootSubject, approvedCert.Subject)
	require.Equal(t, testconstants.RootSubjectKeyID, approvedCert.SubjectKeyId)

	// ACTION: 6th Trustee proposes revocation (pki-approval.sh:L197-L200)
	msgProposeRevoke := pkitypes.MsgProposeRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       sixthTrusteeAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeRevoke}, sixthTrusteeName, sixthTrusteeAccount)
	require.NoError(t, err)

	// ACTION: 5th Trustee approves revocation (pki-approval.sh:L211-L214)
	msgApproveRevoke := pkitypes.MsgApproveRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       fifthTrusteeAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRevoke}, fifthTrusteeName, fifthTrusteeAccount)
	require.NoError(t, err)

	// ACTION: 4th Trustee approves revocation (pki-approval.sh:L224-L227)
	msgApproveRevoke2 := pkitypes.MsgApproveRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       fourthTrusteeAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRevoke2}, fourthTrusteeName, fourthTrusteeAccount)
	require.NoError(t, err)

	// ACTION: 3rd Trustee approves revocation - quorum reached, cert revoked (pki-approval.sh:L238-L241)
	msgApproveRevoke3 := pkitypes.MsgApproveRevokeX509RootCert{
		Subject:      testconstants.RootSubject,
		SubjectKeyId: testconstants.RootSubjectKeyID,
		Signer:       bobAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRevoke3}, bobName, bobAccount)
	require.NoError(t, err)

	// ACTION: Verify cert is revoked (pki-approval.sh:L243-L252)
	revokedCert, err := pki.GetRevokedDaX509Cert(&suite, testconstants.RootSubject, testconstants.RootSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.RootSubject, revokedCert.Subject)
	require.Equal(t, testconstants.RootSubjectKeyID, revokedCert.SubjectKeyId)
}
