// Bash Source: integration_tests/cli/compliance-demo.sh
//
// This test covers the full compliance certification lifecycle:
// - ACTION: Certify model with zigbee and matter types (compliance-demo.sh:L190-L203)
// - ACTION: Query certified/revoked/provisional models (compliance-demo.sh:L227-L277)
// - ACTION: Revoke certification (compliance-demo.sh:L383-L388)
// - ACTION: Re-certify a revoked model (compliance-demo.sh:L460-L464)
// - ACTION: Certify with all optional fields (compliance-demo.sh:L570-L574)
// - ACTION: Update compliance info (compliance-demo.sh:L643-L646)
// - ACTION: Delete compliance info (compliance-demo.sh:L790-L793)
package compliance_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/compliance"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestComplianceDemoTrackCompliance maps to: integration_tests/cli/compliance-demo.sh (certification flow)
// ACTION: Tests full certification lifecycle including certify, query, re-certify by different/same account,
// update compliance info, and delete compliance info (compliance-demo.sh:L16-L816)
func TestComplianceDemoTrackCompliance(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	compliance.DemoTrackCompliance(&suite)
}

// TestDeleteComplianceInfoForAllCertStatuses maps to: integration_tests/cli/compliance-demo.sh (deletion section)
// ACTION: Tests deleting compliance info across all certification statuses (compliance-demo.sh:L790-L800)
func TestDeleteComplianceInfoForAllCertStatuses(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	compliance.DeleteComplianceInfoForAllCertStatuses(&suite)
}
