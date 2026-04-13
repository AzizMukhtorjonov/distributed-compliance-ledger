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
	rootWithSameSubjectAndSkid1Path         = "integration_tests/constants/root_with_same_subject_and_skid_1"
	rootWithSameSubjectAndSkid2Path         = "integration_tests/constants/root_with_same_subject_and_skid_2"
	rootWithSameSubjectAndSkidSubject       = "MFMxCzAJBgNVBAYTAlVTMQ8wDQYDVQQKDAZBbWF6b24xEzARBgNVBAsMAkNBMR4wHAYDVQQDDBVBbWF6b24gUm9vdCBDQSAtIDAxMQ=="
	rootWithSameSubjectAndSkidSubjectKeyID  = "59:A4:84:8E:A4:E0:72:D7:56:3E:AE:9B:AD:33:EF:29:07:06:7A:9B"
)

// TestPKICombineCerts translates pki-combine-certs.sh.
// Tests that multiple certificates with the same subject and SKID can coexist.
func TestPKICombineCerts(t *testing.T) {
	jack := testconstants.JackAccount
	alice := testconstants.AliceAccount
	bob := testconstants.BobAccount

	vid := 65524

	vendorAccount := fmt.Sprintf("vendor_account_%d", vid)
	cliputils.CreateVendorAccount(t, vendorAccount, vid)

	t.Run("ProposeAndApproveFirstRootCert", func(t *testing.T) {
		txResult, err := ProposeAddX509RootCert(rootWithSameSubjectAndSkid1Path, jack,
			"--vid", fmt.Sprintf("%d", vid),
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = ApproveAddX509RootCert(rootWithSameSubjectAndSkidSubject, rootWithSameSubjectAndSkidSubjectKeyID, alice)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = ApproveAddX509RootCert(rootWithSameSubjectAndSkidSubject, rootWithSameSubjectAndSkidSubjectKeyID, bob)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		out, err := QueryX509Cert(rootWithSameSubjectAndSkidSubject, rootWithSameSubjectAndSkidSubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, rootWithSameSubjectAndSkidSubject))
	})

	t.Run("ProposeAndApproveSecondRootCert_SameSubjectSkid", func(t *testing.T) {
		txResult, err := ProposeAddX509RootCert(rootWithSameSubjectAndSkid2Path, jack,
			"--vid", fmt.Sprintf("%d", vid),
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = ApproveAddX509RootCert(rootWithSameSubjectAndSkidSubject, rootWithSameSubjectAndSkidSubjectKeyID, alice)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = ApproveAddX509RootCert(rootWithSameSubjectAndSkidSubject, rootWithSameSubjectAndSkidSubjectKeyID, bob)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Now both certs with same subject+skid should coexist
		out, err := QueryX509Cert(rootWithSameSubjectAndSkidSubject, rootWithSameSubjectAndSkidSubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, rootWithSameSubjectAndSkidSubject))
	})
}
