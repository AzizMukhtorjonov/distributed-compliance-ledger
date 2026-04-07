// Bash Source: integration_tests/cli/vendorinfo-demo-hex.sh
//
// This test verifies vendor info operations using hex-formatted VID values:
// - ACTION: Create vendor account with hex VID 0xA13 (vendorinfo-demo-hex.sh:L22-L28)
// - ACTION: Query VendorInfo for non-existent hex VID (vendorinfo-demo-hex.sh:L31-L33)
// - ACTION: Add VendorInfo using hex VID (vendorinfo-demo-hex.sh:L37-L44)
// - ACTION: Query VendorInfo using hex VID and verify (vendorinfo-demo-hex.sh:L47-L55)
// - ACTION: Update VendorInfo using hex VID (vendorinfo-demo-hex.sh:L58-L68)
// - ACTION: Query updated VendorInfo (vendorinfo-demo-hex.sh:L71-L80)
// - ACTION: Query all vendors (vendorinfo-demo-hex.sh:L83-L92)
package vendorinfo_test

import (
	"testing"

	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/grpc_rest/vendorinfo"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// TestVendorInfoDemoWithHexVid maps to: integration_tests/cli/vendorinfo-demo-hex.sh
// ACTION: Full VendorInfo lifecycle using hex-encoded VID (0xA13 = 2579)
func TestVendorInfoDemoWithHexVid(t *testing.T) {
	suite := utils.SetupTest(t, testconstants.ChainID, false)
	vendorinfo.DemoWithHexVid(&suite)
}
