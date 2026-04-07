// Bash Source: integration_tests/cli/auth-demo.sh
//
// This test exercises the full account lifecycle in the dclauth module:
// - ACTION: Generate keys and propose a new account (auth-demo.sh:L19-L88)
// - ACTION: Approve the proposed account (auth-demo.sh:L124-L127)
// - ACTION: Verify account appears in all-accounts (auth-demo.sh:L138-L140)
// - ACTION: Propose revocation of the account (auth-demo.sh:L188-L191)
// - ACTION: Approve revocation (auth-demo.sh:L232-L235)
// - ACTION: Verify account appears in revoked list (auth-demo.sh:L241-L244)
// - ACTION: Re-propose the revoked account (auth-demo.sh:L285-L288)
// - ACTION: Approve the re-proposed account (auth-demo.sh:L330-L333)
// - ACTION: 2nd revocation cycle (auth-demo.sh:L385-L432)
// - ACTION: Re-propose again and test rejection flow (auth-demo.sh:L472-L606)
// - ACTION: Vendor account approval with 1/3 threshold (auth-demo.sh:L636-L776)
// - ACTION: Vendor account with PID ranges (auth-demo.sh:L780-L1416)
package dclauth_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/dclauth"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestAuthDemo maps to: integration_tests/cli/auth-demo.sh
// ACTION: Full account lifecycle - propose, approve, revoke, reject, vendor roles, PID ranges
func TestAuthDemo(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	// ACTION: Execute the full auth demo which covers account proposal, approval,
	// revocation, rejection, vendor account creation with 1/3 approval threshold,
	// and vendor accounts with PID ranges (auth-demo.sh:L16-L1416)
	dclauth.AuthDemo(&suite)
}
