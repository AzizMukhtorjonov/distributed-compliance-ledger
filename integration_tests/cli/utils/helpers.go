package cliputils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

// CreateAccount generates a new key, proposes and approves the account via jack/alice,
// and returns the account name. roles is a comma-separated string e.g. "NodeAdmin" or "Vendor".
func CreateAccount(t *testing.T, roles string) string {
	t.Helper()

	name := utils.RandString()

	_, err := utils.ExecuteCLI("keys", "add", name, "--keyring-backend", "test")
	require.NoError(t, err)

	addrOut, err := utils.ExecuteCLI("keys", "show", name, "-a", "--keyring-backend", "test")
	require.NoError(t, err)
	addr := stripNewline(addrOut)

	pubkeyOut, err := utils.ExecuteCLI("keys", "show", name, "-p", "--keyring-backend", "test")
	require.NoError(t, err)
	pubkey := stripNewline(pubkeyOut)

	txResult, err := utils.ExecuteTx("tx", "auth", "propose-add-account",
		"--address", addr,
		"--pubkey", pubkey,
		"--roles", roles,
		"--from", testconstants.JackAccount,
	)
	require.NoError(t, err)
	require.Equal(t, uint32(0), txResult.Code)

	_, err = utils.AwaitTxConfirmation(txResult.TxHash)
	require.NoError(t, err)

	txResult, err = utils.ExecuteTx("tx", "auth", "approve-add-account",
		"--address", addr,
		"--from", testconstants.AliceAccount,
	)
	require.NoError(t, err)
	require.Equal(t, uint32(0), txResult.Code)

	_, err = utils.AwaitTxConfirmation(txResult.TxHash)
	require.NoError(t, err)

	return name
}

// CreateVendorAccount creates a vendor account with the given name and vid.
// An optional single pid_ranges argument may be passed.
func CreateVendorAccount(t *testing.T, name string, vid int, pidRanges ...string) string {
	t.Helper()

	_, err := utils.ExecuteCLI("keys", "add", name, "--keyring-backend", "test")
	require.NoError(t, err)

	addrOut, err := utils.ExecuteCLI("keys", "show", name, "-a", "--keyring-backend", "test")
	require.NoError(t, err)
	addr := stripNewline(addrOut)

	pubkeyOut, err := utils.ExecuteCLI("keys", "show", name, "-p", "--keyring-backend", "test")
	require.NoError(t, err)
	pubkey := stripNewline(pubkeyOut)

	args := []string{
		"tx", "auth", "propose-add-account",
		"--address", addr,
		"--pubkey", pubkey,
		"--roles", "Vendor",
		"--vid", strconv.Itoa(vid),
		"--from", testconstants.JackAccount,
	}
	if len(pidRanges) > 0 {
		args = append(args, "--pid_ranges", pidRanges[0])
	}

	txResult, err := utils.ExecuteTx(args...)
	require.NoError(t, err)
	require.Equal(t, uint32(0), txResult.Code)

	_, err = utils.AwaitTxConfirmation(txResult.TxHash)
	require.NoError(t, err)

	txResult, err = utils.ExecuteTx("tx", "auth", "approve-add-account",
		"--address", addr,
		"--from", testconstants.AliceAccount,
	)
	require.NoError(t, err)
	require.Equal(t, uint32(0), txResult.Code)

	_, err = utils.AwaitTxConfirmation(txResult.TxHash)
	require.NoError(t, err)

	return name
}

// CreateModelAndVersion adds a model and a model version for the given vid/pid/sv/svs
// using the provided userAddr (account name) as the signer.
func CreateModelAndVersion(t *testing.T, vid, pid, sv int, svs, userAddr string) {
	t.Helper()

	txResult, err := utils.ExecuteTx("tx", "model", "add-model",
		"--vid", strconv.Itoa(vid),
		"--pid", strconv.Itoa(pid),
		"--deviceTypeID", "1",
		"--productName", "TestProduct",
		"--productLabel", "TestingProductLabel",
		"--partNumber", "1",
		"--commissioningCustomFlow", "0",
		"--enhancedSetupFlowOptions", "0",
		"--from", userAddr,
	)
	require.NoError(t, err)
	require.Equal(t, uint32(0), txResult.Code)

	_, err = utils.AwaitTxConfirmation(txResult.TxHash)
	require.NoError(t, err)

	txResult, err = utils.ExecuteTx("tx", "model", "add-model-version",
		"--vid", strconv.Itoa(vid),
		"--pid", strconv.Itoa(pid),
		"--softwareVersion", strconv.Itoa(sv),
		"--softwareVersionString", svs,
		"--cdVersionNumber", "1",
		"--maxApplicableSoftwareVersion", "10",
		"--minApplicableSoftwareVersion", "1",
		"--from", userAddr,
	)
	require.NoError(t, err)
	require.Equal(t, uint32(0), txResult.Code)

	_, err = utils.AwaitTxConfirmation(txResult.TxHash)
	require.NoError(t, err)
}

// GetHeight queries the node status and returns the latest block height.
func GetHeight() (int64, error) {
	out, err := utils.ExecuteCLI("status")
	if err != nil {
		return 0, err
	}

	var status struct {
		SyncInfo struct {
			LatestBlockHeight string `json:"latest_block_height"`
		} `json:"SyncInfo"`
	}
	if err := json.Unmarshal(out, &status); err != nil {
		return 0, fmt.Errorf("failed to parse status: %v, output: %s", err, string(out))
	}

	h, err := strconv.ParseInt(status.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse latest_block_height: %v", err)
	}

	return h, nil
}

// WaitForHeight polls until the chain reaches target height or timeoutSec is exceeded.
func WaitForHeight(t *testing.T, target int64, timeoutSec int) {
	t.Helper()

	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	for {
		h, err := GetHeight()
		if err == nil && h >= target {
			return
		}
		if time.Now().After(deadline) {
			require.Failf(t, "WaitForHeight timed out",
				"height %d not reached within %d seconds (current: %d, err: %v)",
				target, timeoutSec, h, err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// stripNewline removes trailing newline/whitespace from CLI output bytes.
func stripNewline(b []byte) string {
	s := string(b)
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}
