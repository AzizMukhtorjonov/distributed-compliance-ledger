// Bash Source: integration_tests/cli/pki-demo.sh
//
// This is the largest PKI test covering the full X.509 certificate lifecycle:
// - ACTION: Propose X.509 root certificate (pki-demo.sh:L37-L48)
// - ACTION: Approve root certificate (pki-demo.sh:L55-L64)
// - ACTION: Query approved certificates (pki-demo.sh:L70-L85)
// - ACTION: Add intermediate certificate (pki-demo.sh:L115-L120)
// - ACTION: Add leaf certificate (pki-demo.sh:L125-L130)
// - ACTION: Query certificate chains (pki-demo.sh:L195-L290)
// - ACTION: Revoke leaf certificate (pki-demo.sh:L340-L352)
// - ACTION: Propose revocation of root cert (pki-demo.sh:L480-L520)
// - ACTION: Approve root revocation (pki-demo.sh:L530-L550)
// - ACTION: Assign VID to root certificate (pki-demo.sh:L617-L632)
// - ACTION: Reject proposed root certificate (pki-demo.sh:L780-L820)
// - ACTION: Delete certificate by serial number (pki-demo.sh:L900-L940)
// - ACTION: Certificate by subject queries (pki-demo.sh:L1400-L1818)
package pki_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/pki"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestPkiDemo maps to: integration_tests/cli/pki-demo.sh
// ACTION: Full PKI CRUD lifecycle for DA (Device Attestation) certificates including propose,
// approve, reject, revoke, assign VID, and certificate chain management (pki-demo.sh:L16-L1818)
func TestPkiDemo(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	pki.Demo(&suite)
}
