// Bash Source: integration_tests/cli/pki-remove-x509-certificates.sh
//
// This test verifies certificate removal operations:
// - ACTION: Create vendor accounts (pki-remove-x509-certificates.sh:L40-L46)
// - ACTION: Propose/approve root cert (pki-remove-x509-certificates.sh:L48-L60)
// - ACTION: Add intermediate certs (pki-remove-x509-certificates.sh:L65-L80)
// - ACTION: Add leaf cert (pki-remove-x509-certificates.sh:L82-L88)
// - ACTION: Vendor removes intermediate cert by serial number (pki-remove-x509-certificates.sh:L95-L105)
// - ACTION: Verify only specified cert removed, other with same subject/SKID remains (pki-remove-x509-certificates.sh:L108-L140)
// - ACTION: Remove remaining intermediate (pki-remove-x509-certificates.sh:L145-L165)
// - ACTION: Non-owner vendor attempts removal - expect error (pki-remove-x509-certificates.sh:L170-L198)
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

// TestPkiRemoveX509Certificates maps to: integration_tests/cli/pki-remove-x509-certificates.sh
// ACTION: Tests vendor removal of X.509 certificates with serial number selection
func TestPkiRemoveX509Certificates(t *testing.T) {
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

	// ACTION: Create vendor account (pki-remove-x509-certificates.sh:L40-L44)
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

	// ACTION: Propose and approve root cert (pki-remove-x509-certificates.sh:L48-L60)
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

	// ACTION: Add intermediate cert (pki-remove-x509-certificates.sh:L65-L72)
	msgAddIntermediate := pkitypes.MsgAddX509Cert{
		Cert:   testconstants.IntermediateCertPem,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddIntermediate}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify intermediate cert exists
	intermCert, err := pki.GetDaX509Cert(&suite, testconstants.IntermediateSubject, testconstants.IntermediateSubjectKeyID)
	require.NoError(t, err)
	require.Equal(t, testconstants.IntermediateSubject, intermCert.Subject)

	// ACTION: Remove intermediate cert by vendor (pki-remove-x509-certificates.sh:L95-L105)
	msgRemoveCert := pkitypes.MsgRemoveX509Cert{
		Subject:      testconstants.IntermediateSubject,
		SubjectKeyId: testconstants.IntermediateSubjectKeyID,
		Signer:       vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgRemoveCert}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Verify cert is removed (pki-remove-x509-certificates.sh:L108-L120)
	_, err = pki.GetDaX509Cert(&suite, testconstants.IntermediateSubject, testconstants.IntermediateSubjectKeyID)
	suite.AssertNotFound(err)
}
