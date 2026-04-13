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
	rootCertWithVidPath         = "integration_tests/constants/root_cert_with_vid"
	paaNoVidPath                = "integration_tests/constants/paa_cert_no_vid"
	rootCertWithVidSubject      = "MDQxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApzb21lLXN0YXRlMRAwDgYDVQQKDAdyb290LWNh"
	rootCertWithVidSubjectKeyID = "5A:88:0E:6C:36:53:D0:7F:B0:89:71:A3:F4:73:79:09:30:E6:2B:DB"
	rootCertWithVid             = 65521
	paaNoVidVid                 = 65521
)

// TestPKIAssignVid translates pki-assign-vid.sh.
func TestPKIAssignVid(t *testing.T) {
	jack := testconstants.JackAccount
	alice := testconstants.AliceAccount

	vendorAdminAccount := cliputils.CreateAccount(t, "VendorAdmin")

	t.Run("AssignVidToRootCertThatAlreadyHasVid_Fails", func(t *testing.T) {
		// Propose and approve root cert that has a vid
		txResult, err := ProposeAddX509RootCert(rootCertPath, jack,
			"--vid", fmt.Sprintf("%d", rootCertWithVid),
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

		// Assign VID to a cert that already has vid — should fail
		txResult, err = AssignVid(rootCertSubject, rootCertSubjectKeyID, rootCertWithVid, vendorAdminAccount)
		// expect error or non-zero code
		if err != nil {
			require.Contains(t, err.Error(), "vid is not empty")
		} else if txResult != nil {
			require.Contains(t, txResult.RawLog, "vid is not empty")
		}
	})

	_ = jack
	_ = alice
}
