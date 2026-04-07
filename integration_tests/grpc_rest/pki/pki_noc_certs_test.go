// Bash Source: integration_tests/cli/pki-noc-certs.sh
//
// This test covers NOC (Network Operator Certificate) lifecycle:
// - ACTION: Add NOC root certificate (pki-noc-certs.sh:L44-L55)
// - ACTION: Query NOC root certificates by VID (pki-noc-certs.sh:L60-L86)
// - ACTION: Add NOC ICA certificate (pki-noc-certs.sh:L100-L110)
// - ACTION: Query NOC ICA certificates by VID (pki-noc-certs.sh:L115-L141)
// - ACTION: Query all NOC certificates (pki-noc-certs.sh:L160-L195)
// - ACTION: Remove NOC certificates (pki-noc-certs.sh:L310-L420)
// - ACTION: Revoke NOC certificates (pki-noc-certs.sh:L430-L583)
package pki_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/pki"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestPkiNocCertDemo maps to: integration_tests/cli/pki-noc-certs.sh
// ACTION: Full NOC certificate lifecycle including add/query/remove/revoke for root and ICA NOC certs
func TestPkiNocCertDemo(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	pki.NocCertDemo(&suite)
}
