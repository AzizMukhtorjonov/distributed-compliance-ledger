// Bash Source: integration_tests/cli/model-demo-hex.sh
//
// This test verifies model operations using hex-formatted VID/PID values:
// - ACTION: Create vendor account with hex VID 0xA13 (model-demo-hex.sh:L22-L31)
// - ACTION: Add model using hex VID 0xA13 and PID 0xA11 (model-demo-hex.sh:L37-L46)
// - ACTION: Query model using hex VID/PID and verify (model-demo-hex.sh:L49-L55)
// - ACTION: Get all models and verify count (model-demo-hex.sh:L58-L62)
// - ACTION: Get vendor models by hex VID (model-demo-hex.sh:L65-L72)
// - ACTION: Update model using hex VID/PID (model-demo-hex.sh:L75-L84)
// - ACTION: Delete model using hex VID/PID (model-demo-hex.sh:L95-L103)
// - ACTION: Add and manage model versions with hex IDs (model-demo-hex.sh:L106-L138)
package model_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/model"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestModelDemoWithHexVidAndPid maps to: integration_tests/cli/model-demo-hex.sh
// ACTION: Full model lifecycle (create, query, update, delete) using hex-encoded VID/PID
func TestModelDemoWithHexVidAndPid(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	// ACTION: Execute the DemoWithHexVidAndPid helper covering model CRUD with hex-encoded
	// VID (0xA13) and PID (0xA11) values (model-demo-hex.sh:L16-L138)
	model.DemoWithHexVidAndPid(&suite)
}
