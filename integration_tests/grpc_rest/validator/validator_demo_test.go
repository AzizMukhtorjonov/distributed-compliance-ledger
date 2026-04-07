// Bash Source: integration_tests/cli/validator-demo.sh
//
// This test covers the full validator node lifecycle:
// - ACTION: Create NodeAdmin account (validator-demo.sh:L29-L36)
// - ACTION: Create validator node (validator-demo.sh:L40-L55)
// - ACTION: Disable validator by NodeAdmin (validator-demo.sh:L65-L75)
// - ACTION: Enable validator by NodeAdmin (validator-demo.sh:L85-L95)
// - ACTION: Propose disable validator by Trustee (validator-demo.sh:L105-L120)
// - ACTION: Reject proposed disable by another Trustee (validator-demo.sh:L125-L140)
// - ACTION: Approve proposed disable (validator-demo.sh:L180-L210)
// - ACTION: Verify validator is disabled (validator-demo.sh:L215-L240)
// - ACTION: Re-enable and test disable rejection flow (validator-demo.sh:L250-L480)
// - ACTION: Test invalid pubkey validation (validator-demo.sh:L500-L530)
// - ACTION: Query proposed/disabled/rejected validators lists (validator-demo.sh:L545-L758)
package validator_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/validator"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestValidatorDemo maps to: integration_tests/cli/validator-demo.sh
// ACTION: Full validator lifecycle including create, disable/enable by NodeAdmin,
// propose/approve/reject disable by Trustees, and comprehensive state queries
func TestValidatorDemo(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	validator.Demo(&suite)
}
