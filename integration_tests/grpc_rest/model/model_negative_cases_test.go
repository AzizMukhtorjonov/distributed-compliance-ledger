// Bash Source: integration_tests/cli/model-negative-cases.sh
//
// This test verifies negative/error cases for the model module:
// - ACTION: Non-Vendor adds model - should fail (model-negative-cases.sh:L31-L40)
// - ACTION: Vendor adds model for a VID they don't own - should fail (model-negative-cases.sh:L52-L62)
// - ACTION: Vendor with non-associated Product IDs adds model - should fail (model-negative-cases.sh:L72-L82)
// - ACTION: Vendor updates model for non-associated Product IDs - should fail (model-negative-cases.sh:L113-L112)
// - ACTION: Vendor deletes model for non-associated Product IDs - should fail (model-negative-cases.sh:L125-L134)
// - ACTION: Same model added twice - should fail (model-negative-cases.sh:L145-L155)
// - ACTION: Get model for non-existent VID/PID (model-negative-cases.sh:L158-L161)
// - ACTION: Get model for invalid VID/PID (model-negative-cases.sh:L164-L166)
package model_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/model"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestModelAddByNonVendor maps to: integration_tests/cli/model-negative-cases.sh:L31-L40
// ACTION: Non-Vendor account attempts to add model - expect ErrUnauthorized
func TestModelAddByNonVendor(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.AddModelByNonVendor(&suite)
}

// TestModelAddByDifferentVendor maps to: integration_tests/cli/model-negative-cases.sh:L52-L62
// ACTION: Vendor adds model for a different VID - expect ErrUnauthorized
func TestModelAddByDifferentVendor(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.AddModelByDifferentVendor(&suite)
}

// TestModelAddByVendorWithNonAssociatedProductIDs maps to: integration_tests/cli/model-negative-cases.sh:L72-L82
// ACTION: Vendor with specific PID range adds model outside that range - expect error
func TestModelAddByVendorWithNonAssociatedProductIDs(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.AddModelByVendorWithNonAssociatedProductIDs(&suite)
}

// TestModelUpdateByVendorWithNonAssociatedProductIDs maps to: integration_tests/cli/model-negative-cases.sh:L113-L112
// ACTION: Vendor with specific PID range updates model outside that range - expect error
func TestModelUpdateByVendorWithNonAssociatedProductIDs(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.UpdateModelByVendorWithNonAssociatedProductIDs(&suite)
}

// TestModelDeleteByVendorWithNonAssociatedProductIDs maps to: integration_tests/cli/model-negative-cases.sh:L125-L134
// ACTION: Vendor with specific PID range deletes model outside that range - expect error
func TestModelDeleteByVendorWithNonAssociatedProductIDs(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelByVendorWithNonAssociatedProductIDs(&suite)
}

// TestModelAddTwice maps to: integration_tests/cli/model-negative-cases.sh:L145-L155
// ACTION: Adding the same model twice - expect error on duplicate
func TestModelAddTwice(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.AddModelTwice(&suite)
}

// TestModelGetForUnknown maps to: integration_tests/cli/model-negative-cases.sh:L158-L161
// ACTION: Querying model for non-existent VID/PID - expect NotFound
func TestModelGetForUnknown(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.GetModelForUnknown(&suite)
}

// TestModelGetForInvalidVidPid maps to: integration_tests/cli/model-negative-cases.sh:L164-L166
// ACTION: Querying model with invalid VID/PID values - expect error
func TestModelGetForInvalidVidPid(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.GetModelForInvalidVidPid(&suite)
}
