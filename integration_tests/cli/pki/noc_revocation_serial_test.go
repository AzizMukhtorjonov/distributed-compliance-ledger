package pki

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	cliputils "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/cli/utils"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestPKINocRevocationWithSerialNumber translates pki-noc-revocation-with-serial-number.sh.
func TestPKINocRevocationWithSerialNumber(t *testing.T) {
	vendorAccount := fmt.Sprintf("vendor_account_%d", nocVid)
	cliputils.CreateVendorAccount(t, vendorAccount, nocVid)

	t.Run("RevokeNocRootCertBySerial", func(t *testing.T) {
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

		// Try to revoke with invalid serial number
		txResult, err = RevokeNocRootCert(nocRootCert1Subject, nocRootCert1SubjectKeyID, vendorAccount,
			"--serial-number", "invalid",
		)
		require.NoError(t, err)
		require.Equal(t, uint32(404), txResult.Code)

		// Revoke only first root cert by serial number (child certs should NOT be revoked)
		txResult, err = RevokeNocRootCert(nocRootCert1Subject, nocRootCert1SubjectKeyID, vendorAccount,
			"--serial-number", nocRootCert1SerialNumber,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Only first root cert should be in revoked list
		out, err := QueryAllRevokedNocRootCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), nocRootCert1SerialNumber)
		require.NotContains(t, string(out), nocRootCert1CopySerialNumber)
		require.NotContains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocCert1Subject))
		require.NotContains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocLeafCert1Subject))

		// Second root cert should still be active
		out, err = QueryNocRootCerts(nocVid)
		require.NoError(t, err)
		require.Contains(t, string(out), nocRootCert1CopySerialNumber)
		require.NotContains(t, string(out), nocRootCert1SerialNumber)

		// ICA and leaf should still be active
		out, err = QueryNocX509IcaCerts(nocVid)
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocCert1Subject))
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocLeafCert1Subject))

		// Revoke second root cert with revoke-child=true
		txResult, err = RevokeNocRootCert(nocRootCert1Subject, nocRootCert1SubjectKeyID, vendorAccount,
			"--serial-number", nocRootCert1CopySerialNumber,
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

		// No active NOC certs
		out, err = QueryAllNocX509Certs()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")
	})

	t.Run("RevokeNocIcaCertBySerial", func(t *testing.T) {
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

		// Try to revoke with invalid serial number
		txResult, err = RevokeNocX509IcaCert(nocCert2Subject, nocCert2SubjectKeyID, vendorAccount,
			"--serial-number", "invalid",
		)
		require.NoError(t, err)
		require.Equal(t, uint32(404), txResult.Code)

		// Revoke only first ICA cert by serial (child should not be revoked)
		txResult, err = RevokeNocX509IcaCert(nocCert2Subject, nocCert2SubjectKeyID, vendorAccount,
			"--serial-number", nocCert2SerialNumber,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		out, err := QueryAllRevokedNocX509IcaCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), nocCert2SerialNumber)
		require.NotContains(t, string(out), nocRevChildCert2CopySerialNumber)
		require.NotContains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocLeafCert2Subject))

		// Second ICA cert should still be active
		out, err = QueryNocCert("--subject-key-id", nocCert2SubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), nocRevChildCert2CopySerialNumber)
		require.NotContains(t, string(out), nocCert2SerialNumber)

		// Revoke second ICA cert with revoke-child=true
		txResult, err = RevokeNocX509IcaCert(nocCert2Subject, nocCert2SubjectKeyID, vendorAccount,
			"--serial-number", nocRevChildCert2CopySerialNumber,
			"--revoke-child", "true",
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Both ICA certs and leaf should be revoked
		out, err = QueryAllRevokedNocX509IcaCerts()
		require.NoError(t, err)
		require.Contains(t, string(out), nocCert2SerialNumber)
		require.Contains(t, string(out), nocRevChildCert2CopySerialNumber)
		require.Contains(t, string(out), nocLeafCert2SerialNumber)

		// Only root cert should remain
		out, err = QueryAllNocX509Certs()
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, nocRootCert2Subject))
		require.NotContains(t, string(out), nocCert2Subject)
		require.NotContains(t, string(out), nocLeafCert2Subject)
	})
}
