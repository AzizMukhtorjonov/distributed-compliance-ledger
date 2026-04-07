// Bash Source: integration_tests/cli/compliance-provisioning.sh
//
// This test covers the provisioning track of the compliance module:
// - ACTION: Provision an uncertified model for ZB and Matter (compliance-provisioning.sh:L50-L62)
// - ACTION: Attempt re-provisioning (expect error) (compliance-provisioning.sh:L67-L81)
// - ACTION: Query provisional models (compliance-provisioning.sh:L85-L101)
// - ACTION: Verify certified/revoked models are not found (compliance-provisioning.sh:L105-L142)
// - ACTION: Get compliance info showing provisional status (compliance-provisioning.sh:L147-L172)
// - ACTION: Certify a second model and verify provisioning rejects certified models (compliance-provisioning.sh:L236-L264)
// - ACTION: Revoke a third model and verify provisioning rejects revoked models (compliance-provisioning.sh:L277-L292)
// - ACTION: Certify the provisioned model and revoke the matter provisioned model (compliance-provisioning.sh:L298-L314)
// - ACTION: Provision with all optional fields (compliance-provisioning.sh:L466-L511)
// - ACTION: Certify with some optional fields and verify field merging (compliance-provisioning.sh:L517-L575)
package compliance_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/compliance"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestComplianceDemoTrackProvision maps to: integration_tests/cli/compliance-provisioning.sh
// ACTION: Full provisioning lifecycle including provision, query, certify provisioned, revoke provisioned,
// and field merging between provision and certify operations (compliance-provisioning.sh:L16-L580)
func TestComplianceDemoTrackProvision(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	compliance.DemoTrackProvision(&suite)
}
