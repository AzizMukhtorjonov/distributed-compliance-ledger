// Bash Source: integration_tests/cli/model-validation-cases.sh
//
// This test covers model validation edge cases:
// - ACTION: Vendor with product IDs adds model matching PID range (model-validation-cases.sh:L34-L45)
// - ACTION: Update model by vendor with matching product IDs (model-validation-cases.sh:L50-L60)
// - ACTION: Delete model with associated model versions (model-validation-cases.sh:L101-L140)
// - ACTION: Delete model with certified model versions - expect error (model-validation-cases.sh:L151-L200)
// - ACTION: Delete model version individually (model-validation-cases.sh:L208-L246)
// - ACTION: Delete model version before deleting model (model-validation-cases.sh:L252-L300)
// - ACTION: Delete model version with different VID - expect error (model-validation-cases.sh:L310-L370)
// - ACTION: Delete model version that doesn't exist - expect error (model-validation-cases.sh:L380-L410)
// - ACTION: Delete model version not by creator - expect error (model-validation-cases.sh:L420-L480)
// - ACTION: Delete model version when certified - expect error (model-validation-cases.sh:L490-L570)
package model_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/model"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestModelAddByVendorProductIDs maps to: integration_tests/cli/model-validation-cases.sh:L34-L45
// ACTION: Add model where vendor's PID range includes the model's PID - should succeed
func TestModelAddByVendorProductIDs(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.AddModelByVendorWithProductIDs(&suite)
}

// TestModelUpdateByVendorProductIDs maps to: integration_tests/cli/model-validation-cases.sh:L50-L60
// ACTION: Update model where vendor's PID range includes the model's PID - should succeed
func TestModelUpdateByVendorProductIDs(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.UpdateByVendorWithProductIDs(&suite)
}

// TestDeleteModelWithAssociatedModelVersions maps to: integration_tests/cli/model-validation-cases.sh:L101-L140
// ACTION: Delete model that has associated model versions - versions should also be cleaned up
func TestDeleteModelWithAssociatedModelVersions(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelWithAssociatedModelVersions(&suite)
}

// TestDeleteModelWithAssociatedModelVersionsCertified maps to: integration_tests/cli/model-validation-cases.sh:L151-L200
// ACTION: Attempt to delete model with certified model versions - expect error
func TestDeleteModelWithAssociatedModelVersionsCertified(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelWithAssociatedModelVersionsCertified(&suite)
}

// TestDeleteModelVersionStandalone maps to: integration_tests/cli/model-validation-cases.sh:L208-L246
// ACTION: Delete a model version independently
func TestDeleteModelVersionStandalone(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelVersion(&suite)
}

// TestDeleteModelVersionBeforeDeletingModel maps to: integration_tests/cli/model-validation-cases.sh:L252-L300
// ACTION: Delete model version before deleting the parent model
func TestDeleteModelVersionBeforeDeletingModel(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelVersionBeforeDeletingModel(&suite)
}

// TestDeleteModelVersionDifferentVid maps to: integration_tests/cli/model-validation-cases.sh:L310-L370
// ACTION: Attempt to delete model version using account with different VID - expect error
func TestDeleteModelVersionDifferentVid(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelVersionDifferentVid(&suite)
}

// TestDeleteModelVersionDoesNotExist maps to: integration_tests/cli/model-validation-cases.sh:L380-L410
// ACTION: Attempt to delete model version that doesn't exist - expect error
func TestDeleteModelVersionDoesNotExist(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelVersionDoesNotExist(&suite)
}

// TestDeleteModelVersionNotByCreator maps to: integration_tests/cli/model-validation-cases.sh:L420-L480
// ACTION: Attempt to delete model version using account other than creator - expect error
func TestDeleteModelVersionNotByCreator(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelVersionNotByCreator(&suite)
}

// TestDeleteModelVersionCertified maps to: integration_tests/cli/model-validation-cases.sh:L490-L570
// ACTION: Attempt to delete certified model version - expect error
func TestDeleteModelVersionCertified(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	model.DeleteModelVersionCertified(&suite)
}
