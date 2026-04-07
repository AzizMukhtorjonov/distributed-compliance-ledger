// Bash Source: integration_tests/cli/compliance-revocation.sh
//
// This test covers the revocation track of the compliance module:
// - ACTION: Revoke certification for an uncertified model (compliance-revocation.sh:L48-L55)
// - ACTION: Attempt re-revocation (expect error) (compliance-revocation.sh:L59-L73)
// - ACTION: Revoke for matter certification type (compliance-revocation.sh:L77-L83)
// - ACTION: Query revoked/certified models (compliance-revocation.sh:L87-L124)
// - ACTION: Get compliance info showing revoked status=3 (compliance-revocation.sh:L129-L153)
// - ACTION: Certify the revoked model (compliance-revocation.sh:L219-L224)
// - ACTION: Verify model state transitions (compliance-revocation.sh:L228-L265)
// - ACTION: Query lists for correct state (compliance-revocation.sh:L297-L347)
package compliance_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/compliance"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestComplianceDemoTrackRevocation maps to: integration_tests/cli/compliance-revocation.sh
// ACTION: Full revocation lifecycle including revoke uncertified, re-revoke (error),
// certify revoked, and state verification (compliance-revocation.sh:L16-L351)
func TestComplianceDemoTrackRevocation(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	compliance.DemoTrackRevocation(&suite)
}
