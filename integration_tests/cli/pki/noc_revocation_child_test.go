package pki

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	cliputils "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/cli/utils"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

const (
	nocRevChildCert2CopyPath         = "integration_tests/constants/noc_cert_2_copy"
	nocRevChildCert2CopySerialNumber = "157351092243199289154908179633004790674818411696"

	nocLeafCert2Path         = "integration_tests/constants/noc_leaf_cert_2"
	nocLeafCert2Subject      = "MIGBMQswCQYDVQQGEwJVWjETMBEGA1UECAwKU29tZSBTdGF0ZTETMBEGA1UEBwwKU29tZSBTdGF0ZTEYMBYGA1UECgwPRXhhbXBsZSBDb21wYW55MRkwFwYDVQQLDBBUZXN0aW5nIERpdmlzaW9uMRMwEQYDVQQDDApOT0MtbGVhZi0y"
	nocLeafCert2SubjectKeyID = "F7:2D:E5:60:05:1E:06:45:E6:17:09:DE:1A:0C:B7:AE:19:66:EA:D5"
	nocLeafCert2SerialNumber = "628585745496304216074570439204763956375973944746"
)

// TestPKINocRevocationWithRevokingChild translates pki-noc-revocation-with-revoking-child.sh.
func TestPKINocRevocationWithRevokingChild(t *testing.T) {
	vendorAccount := fmt.Sprintf("vendor_account_%d", nocVid)
	cliputils.CreateVendorAccount(t, vendorAccount, nocVid)

	t.Run("RevokeNocRootCertWithChildFlag", func(t *testing.T) {
		// Add root certs
		txResult, err := AddNocRootCert(nocRootCert1Path, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = AddNocRootCert(nocRootCert1CopyPath, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Add ICA and leaf certs
		txResult, err = AddNocX509IcaCert(nocCert1Path, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = AddNocX509IcaCert(nocLeafCert1Path, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Verify certs exist
		out, err := QueryAllNocRootCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), nocRootCert1SerialNumber)
		require.Contains(t, string(out), nocRootCert1CopySerialNumber)

		// Revoke root cert with revoke-child=true
		txResult, err = RevokeNocRootCert(nocRootCert1Subject, nocRootCert1SubjectKeyID, vendorAccount,
			"--revoke-child", "true",
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Both root certs should be revoked
		out, err = QueryAllRevokedNocRootCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), nocRootCert1SerialNumber)
		require.Contains(t, string(out), nocRootCert1CopySerialNumber)

		// ICA and leaf should also be revoked
		out, err = QueryAllRevokedNocX509IcaCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocCert1Subject))
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocLeafCert1Subject))

		// NOC certs by VID should be empty
		out, err = QueryNocRootCerts(nocVid)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryNocX509IcaCerts(nocVid)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		// Approved x509 certs should be empty
		out, err = QueryAllX509Certs()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")
	})

	t.Run("RevokeNocIcaCertWithChildFlag", func(t *testing.T) {
		// Add root cert
		txResult, err := AddNocRootCert(nocRootCert2Path, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Add ICA certs
		txResult, err = AddNocX509IcaCert(nocCert2Path, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = AddNocX509IcaCert(nocRevChildCert2CopyPath, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = AddNocX509IcaCert(nocLeafCert2Path, vendorAccount)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Revoke ICA cert with revoke-child=true
		txResult, err = RevokeNocX509IcaCert(nocCert2Subject, nocCert2SubjectKeyID, vendorAccount,
			"--revoke-child", "true",
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Both ICA and leaf should be revoked
		out, err := QueryAllRevokedNocX509IcaCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocCert2Subject))
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocLeafCert2Subject))
		require.Contains(t, string(out), nocCert2SerialNumber)
		require.Contains(t, string(out), nocRevChildCert2CopySerialNumber)
		require.Contains(t, string(out), nocLeafCert2SerialNumber)

		// Root should not be in revoked ICA list
		require.NotContains(t, string(out), fmt.Sprintf(`"subject": "%s`, nocRootCert2Subject))

		// NOC certs by VID should not contain ICA/leaf
		out, err = QueryNocX509IcaCerts(nocVid)
		require.NoError(t, err)
		require.NotContains(t, string(out), nocCert2Subject)
		require.NotContains(t, string(out), nocLeafCert2Subject)

		// All NOC certs should not contain revoked ICA/leaf but should still have root
		out, err = QueryAllNocX509Certs()
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocRootCert2Subject))
		require.NotContains(t, string(out), nocCert2Subject)
		require.NotContains(t, string(out), nocLeafCert2Subject)
	})
}
