package pki

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	cliputils "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/cli/utils"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestPKIApproval translates pki-approval.sh.
// This test covers approval flows with multiple trustees and quorum requirements.
func TestPKIApproval(t *testing.T) {
	jack := testconstants.JackAccount
	alice := testconstants.AliceAccount
	bob := testconstants.BobAccount

	userAccount := cliputils.CreateAccount(t, "CertificationCenter")

	// Add 3 new Trustee accounts to get 6 total (4 approvals needed for 2/3 quorum)
	fourthTrustee := cliputils.CreateAccount(t, "Trustee")
	fifthTrustee := cliputils.CreateAccount(t, "Trustee")
	sixthTrustee := cliputils.CreateAccount(t, "Trustee")

	t.Run("NonTrustee_CannotProposeRootCert", func(t *testing.T) {
		txResult, err := ProposeAddX509RootCert(rootCertPath, userAccount,
			"--vid", fmt.Sprintf("%d", pkiDemoVid),
		)
		require.NoError(t, err)
		require.NotEqual(t, uint32(0), txResult.Code)
		_, _ = utils.AwaitTxConfirmation(txResult.TxHash)
	})

	t.Run("ProposeAndApproveWithQuorum", func(t *testing.T) {
		// Jack proposes
		txResult, err := ProposeAddX509RootCert(rootCertPath, jack,
			"--vid", fmt.Sprintf("%d", pkiDemoVid),
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// With 6 trustees, need 4 approvals. Jack proposed (counts as approval).
		// Alice approves
		txResult, err = ApproveAddX509RootCert(rootCertSubject, rootCertSubjectKeyID, alice)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Still in proposed (only 2 approvals)
		out, err := QueryProposedX509RootCert(rootCertSubject, rootCertSubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, rootCertSubject))

		// Bob approves
		txResult, err = ApproveAddX509RootCert(rootCertSubject, rootCertSubjectKeyID, bob)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// 4th trustee approves — now should be approved
		txResult, err = ApproveAddX509RootCert(rootCertSubject, rootCertSubjectKeyID, fourthTrustee)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		out, err = QueryX509Cert(rootCertSubject, rootCertSubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), fmt.Sprintf(`"subject": "%s"`, rootCertSubject))
		require.Contains(t, string(out), `"isRoot": true`)

		out, err = QueryProposedX509RootCert(rootCertSubject, rootCertSubjectKeyID)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")
	})

	t.Run("AddIntermediateAndLeaf", func(t *testing.T) {
		txResult, err := AddX509Cert(intermediateCertPath, "jack")
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = AddX509Cert(leafCertPath, "jack")
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)
	})

	_ = fifthTrustee
	_ = sixthTrustee
}
