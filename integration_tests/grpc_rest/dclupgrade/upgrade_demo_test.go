// Bash Source: integration_tests/cli/upgrade-demo.sh
//
// This test covers the full upgrade module lifecycle:
// - ACTION: Trustee proposes upgrade (upgrade-demo.sh:L37-L50)
// - ACTION: Query proposed upgrade and verify fields (upgrade-demo.sh:L52-L70)
// - ACTION: Another trustee approves upgrade (upgrade-demo.sh:L75-L85)
// - ACTION: Verify upgrade is approved state (upgrade-demo.sh:L88-L115)
// - ACTION: Non-Trustee proposes upgrade - expect error (upgrade-demo.sh:L125-L150)
// - ACTION: Non-Trustee approves upgrade - expect error (upgrade-demo.sh:L155-L185)
// - ACTION: Same upgrade proposed twice - expect error (upgrade-demo.sh:L195-L225)
// - ACTION: Propose and reject upgrade flow (upgrade-demo.sh:L232-L265)
// - ACTION: Verify rejected upgrade queries (upgrade-demo.sh:L270-L354)
package dclupgrade_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclupgrade"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestUpgradeDemo maps to: integration_tests/cli/upgrade-demo.sh
// ACTION: Full upgrade lifecycle - propose, approve, reject, error cases
func TestUpgradeDemo(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	dclupgrade.Demo(&suite)
}

// TestUpgradeProposeByNonTrustee maps to: integration_tests/cli/upgrade-demo.sh:L125-L150
// ACTION: Non-Trustee attempts to propose upgrade - expect ErrUnauthorized
func TestUpgradeProposeByNonTrustee(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	dclupgrade.ProposeUpgradeByNonTrustee(&suite)
}

// TestUpgradeApproveByNonTrustee maps to: integration_tests/cli/upgrade-demo.sh:L155-L185
// ACTION: Non-Trustee attempts to approve upgrade - expect ErrUnauthorized
func TestUpgradeApproveByNonTrustee(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	dclupgrade.ApproveUpgradeByNonTrustee(&suite)
}

// TestUpgradeProposeTwice maps to: integration_tests/cli/upgrade-demo.sh:L195-L225
// ACTION: Same upgrade proposed twice - expect ErrProposedUpgradeAlreadyExists
func TestUpgradeProposeTwice(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	dclupgrade.ProposeUpgradeTwice(&suite)
}

// TestUpgradeProposeAndReject maps to: integration_tests/cli/upgrade-demo.sh:L232-L354
// ACTION: Propose and reject flow - verify state clears after rejection
func TestUpgradeProposeAndReject(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	dclupgrade.ProposeAndRejectUpgrade(&suite)
}
