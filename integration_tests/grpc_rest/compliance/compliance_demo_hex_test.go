// Bash Source: integration_tests/cli/compliance-demo-hex.sh
//
// This test verifies compliance operations using hex-formatted VID/PID values:
// - ACTION: Add model and model version using hex VID/PID (compliance-demo-hex.sh:L49-L70)
// - ACTION: Certify model with zigbee and matter types using hex values (compliance-demo-hex.sh:L128-L141)
// - ACTION: Query compliance info using hex VID/PID (compliance-demo-hex.sh:L219-L226)
// - ACTION: Revoke certification using hex values (compliance-demo-hex.sh:L268-L273)
// - ACTION: Re-certify using hex values (compliance-demo-hex.sh:L324-L329)
package compliance_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/compliance"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestComplianceDemoTrackComplianceWithHexVidAndPid maps to: integration_tests/cli/compliance-demo-hex.sh
// ACTION: Full compliance lifecycle using hex-encoded VID (0xA13=2579) and PID (0xA11=2577)
func TestComplianceDemoTrackComplianceWithHexVidAndPid(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	compliance.DemoTrackComplianceWithHexVidAndPid(&suite)
}