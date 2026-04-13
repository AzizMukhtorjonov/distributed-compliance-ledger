package utils

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// TxResult represents the core standard JSON output emitted by `dcld` node
type TxResult struct {
	Code   uint32 `json:"code"`
	TxHash string `json:"txhash"`
	RawLog string `json:"raw_log"`
}

// ExecuteCLI runs a `dcld` command as a subprocess.
// It sets standard cosmos flags like --output json and keyring-backend test
// if they are not already provided in the args for transaction commands.
func ExecuteCLI(args ...string) ([]byte, error) {
	cmd := exec.Command("dcld", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("cmd execution failed: %v, output: %s", err, string(out))
	}
	// Many cosmos CLI commands output warnings before JSON.
	// Strip out everything before the first `{` to ensure clean parsing.
	strOut := string(out)
	idx := strings.Index(strOut, "{")
	if idx > 0 {
		strOut = strOut[idx:]
	}
	return []byte(strOut), nil
}

// ExecuteTx runs a specific `dcld tx` command with standard testing flags 
// and returns the parsed TxResult for immediate error checking.
func ExecuteTx(args ...string) (*TxResult, error) {
	args = append(args, "--yes", "-o", "json", "--keyring-backend", "test")
	out, err := ExecuteCLI(args...)
	if err != nil {
		return nil, err
	}

	var res TxResult
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, fmt.Errorf("failed to parse TxResult: %v, output: %s", err, string(out))
	}

	return &res, nil
}

// QueryTx runs `dcld query tx [txHash]` mimicking exactly how the REST /grpc tests
// await for block inclusion.
func QueryTx(txHash string) ([]byte, error) {
	return ExecuteCLI("query", "tx", txHash, "-o", "json")
}

// AwaitTxConfirmation replaces shell `get_txn_result` retry looping, resolving brittleness.
func AwaitTxConfirmation(txHash string) ([]byte, error) {
	var result []byte
	var err error

	for i := 0; i < 20; i++ {
		result, err = QueryTx(txHash)
		if err == nil && !strings.Contains(string(result), "not found") {
			return result, nil // Success, confirmed
		}
		time.Sleep(2 * time.Second)
	}
	
	if err != nil {
		return result, fmt.Errorf("transaction %s not confirmed after 20 attempts. Last error: %v", txHash, err)
	}
	return result, fmt.Errorf("transaction %s not confirmed after 40 seconds (not found)", txHash)
}
