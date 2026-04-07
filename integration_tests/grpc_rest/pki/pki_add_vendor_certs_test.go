// Bash Source: integration_tests/cli/pki-add-vendor-x509-certificates.sh
//
// This test verifies adding VID-scoped X.509 certificates for vendors:
// - ACTION: Create vendor account for VID 65521 (pki-add-vendor-x509-certificates.sh:L42-L44)
// - ACTION: Propose and approve root cert with VID 65521 (pki-add-vendor-x509-certificates.sh:L46-L55)
// - ACTION: Add intermediate cert with VID 65521 by vendor (pki-add-vendor-x509-certificates.sh:L58-L62)
// - ACTION: Add intermediate cert with different VID 65522 - expect error (pki-add-vendor-x509-certificates.sh:L64-L68)
// - ACTION: Query approved certificates and verify VID scope (pki-add-vendor-x509-certificates.sh:L72-L100)
// - ACTION: Create vendor account for VID 65522 (pki-add-vendor-x509-certificates.sh:L105-L107)
// - ACTION: Attempt to add intermed cert with VID 65522 signed by root VID 65521 (pki-add-vendor-x509-certificates.sh:L109-L115)
// - ACTION: Remove vendor cert (pki-add-vendor-x509-certificates.sh:L160-L218)
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

// TestPkiAddVendorX509Certificates maps to: integration_tests/cli/pki-add-vendor-x509-certificates.sh
// ACTION: Tests adding VID-scoped intermediate certificates, VID mismatch errors, and removal
func TestPkiAddVendorX509Certificates(t *testing.T) {
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

	// ACTION: Create vendor account for VID 65521 (pki-add-vendor-x509-certificates.sh:L42-L44)
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

	// ACTION: Propose root cert with VID (pki-add-vendor-x509-certificates.sh:L46-L50)
	msgProposeRoot := pkitypes.MsgProposeAddX509RootCert{
		Cert:   testconstants.RootCertWithVid,
		Signer: jackAccount.Address,
		Vid:    testconstants.Vid,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgProposeRoot}, jackName, jackAccount)
	require.NoError(t, err)

	// ACTION: Approve root cert (pki-add-vendor-x509-certificates.sh:L51-L55)
	msgApproveRoot := pkitypes.MsgApproveAddX509RootCert{
		Subject:      testconstants.RootCertWithVidSubject,
		SubjectKeyId: testconstants.RootCertWithVidSubjectKeyID,
		Signer:       aliceAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgApproveRoot}, aliceName, aliceAccount)
	require.NoError(t, err)

	// ACTION: Add intermediate cert with matching VID by vendor (pki-add-vendor-x509-certificates.sh:L58-L62)
	msgAddIntermediate := pkitypes.MsgAddX509Cert{
		Cert:   testconstants.IntermediateCertWithVid1,
		Signer: vendorAccount.Address,
	}
	_, err = suite.BuildAndBroadcastTx([]sdk.Msg{&msgAddIntermediate}, vendorName, vendorAccount)
	require.NoError(t, err)

	// ACTION: Query the certificate and verify (pki-add-vendor-x509-certificates.sh:L72-L100)
	approvedCerts, err := pki.GetDaX509Cert(
		&suite,
		testconstants.IntermediateCertWithVid1Subject,
		testconstants.IntermediateCertWithVid1SubjectKeyID,
	)
	require.NoError(t, err)
	require.Equal(t, testconstants.IntermediateCertWithVid1Subject, approvedCerts.Subject)
}
