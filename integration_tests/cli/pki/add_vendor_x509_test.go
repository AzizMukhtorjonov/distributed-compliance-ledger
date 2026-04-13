package pki

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	cliputils "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/cli/utils"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

const (
	leafCertWithVid65522Path         = "integration_tests/constants/leaf_cert_with_vid_65522"
	leafCertWithVid65522Subject      = "MDExCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApzb21lLXN0YXRlMQ0wCwYDVQQKDARsZWFm"
	leafCertWithVid65522SubjectKeyID = "30:F4:65:75:14:20:B2:AF:3D:14:71:17:AC:49:90:93:3E:24:A0:1F"
	intermCertWithVid1Path           = "integration_tests/constants/intermediate_cert_with_vid_1"
	intermCertWithVid1Subject        = "MDwxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApzb21lLXN0YXRlMRgwFgYDVQQKDA9pbnRlcm1lZGlhdGUtY2E="
	intermCertWithVid1SubjectKeyID   = "4E:3B:73:F4:70:4D:C2:98:0D:DB:C8:5A:5F:02:3B:BF:86:25:56:2B"
	addVendorVid                     = 65522
)

// TestPKIAddVendorX509Certificates translates pki-add-vendor-x509-certificates.sh.
func TestPKIAddVendorX509Certificates(t *testing.T) {
	jack := testconstants.JackAccount
	alice := testconstants.AliceAccount

	vendorAccount := fmt.Sprintf("vendor_account_%d", addVendorVid)
	cliputils.CreateVendorAccount(t, vendorAccount, addVendorVid)

	vendorAccount1 := fmt.Sprintf("vendor_account_%d", pkiDemoVid)
	cliputils.CreateVendorAccount(t, vendorAccount1, pkiDemoVid)

	t.Run("ProposeAndApproveRootCert", func(t *testing.T) {
		txResult, err := ProposeAddX509RootCert(rootCertPath, jack,
			"--vid", fmt.Sprintf("%d", pkiDemoVid),
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = ApproveAddX509RootCert(rootCertSubject, rootCertSubjectKeyID, alice)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)
	})

	t.Run("AddIntermediateCert_WithVendorVid", func(t *testing.T) {
		txResult, err := AddX509Cert(intermediateCertPath, vendorAccount1)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		out, err := QueryX509Cert(intermediateCertSubject, intermediateCertSubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, intermediateCertSubject))
	})

	t.Run("AddLeafCert_WrongVendor_Fails", func(t *testing.T) {
		// Adding leaf cert signed by intermediate belonging to vid=1 from a vendor with vid=65522
		txResult, err := AddX509Cert(leafCertPath, vendorAccount)
		// This should fail — vendor VID mismatch
		if err == nil && txResult != nil {
			require.NotEqual(t, uint32(0), txResult.Code)
		}
	})
}
