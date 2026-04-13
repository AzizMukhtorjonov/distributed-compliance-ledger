package dclauth

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	testconstants "github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/constants"
	"github.com/zigbee-alliance/distributed-compliance-ledger/integration_tests/utils"
)

func TestAuthDemoNodeAdmin(t *testing.T) {
	jack := testconstants.JackAccount
	alice := testconstants.AliceAccount
	bob := testconstants.BobAccount

	// Generate new user key
	name := fmt.Sprintf("user%d", rand.Intn(99999))
	err := AddKey(name)
	require.NoError(t, err)

	userAddr, err := GetAddress(name)
	require.NoError(t, err)
	userPubkey, err := GetPubkey(name)
	require.NoError(t, err)

	jackAddr, err := GetAddress(jack)
	require.NoError(t, err)
	aliceAddr, err := GetAddress(alice)
	require.NoError(t, err)
	bobAddr, err := GetAddress(bob)
	require.NoError(t, err)

	t.Run("InitialState_NotFound", func(t *testing.T) {
		out, err := QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryProposedAccountToRevoke(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryRevokedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")
	})

	t.Run("InitialListsEmpty", func(t *testing.T) {
		out, err := QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")

		out, err = QueryAllProposedAccountsToRevoke()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")

		out, err = QueryAllRevokedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")
	})

	t.Run("JackProposes", func(t *testing.T) {
		txResult, err := utils.ExecuteTx("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "NodeAdmin",
			"--from", jack,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Account not yet active
		out, err := QueryAllAccountsRaw()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		// Now in proposed list
		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), jackAddr)
		require.Contains(t, string(out), "Jack is proposing this account")
		require.NotContains(t, string(out), aliceAddr)

		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
	})

	t.Run("AliceApproves", func(t *testing.T) {
		txResult, err := utils.ExecuteTx("tx", "auth", "approve-add-account",
			"--address", userAddr,
			"--info", "Alice is approving this account",
			"--from", alice,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Alice cannot reject after approving
		txBad, err := utils.ExecuteTx("tx", "auth", "reject-add-account",
			"--address", userAddr,
			"--info", "Alice is rejecting this account",
			"--from", alice,
		)
		// Either error or non-zero code
		if err == nil {
			require.NotEqual(t, uint32(0), txBad.Code)
		}

		// Account is now active
		out, err := QueryAllAccountsRaw()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), jackAddr)
		require.Contains(t, string(out), aliceAddr)
		require.Contains(t, string(out), "Alice is approving this account")
		require.Contains(t, string(out), "Jack is proposing this account")

		// No longer in proposed list
		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")

		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")
	})

	t.Run("AliceProposeRevoke", func(t *testing.T) {
		txResult, err := utils.ExecuteTx("tx", "auth", "propose-revoke-account",
			"--address", userAddr,
			"--info", "Alice proposes to revoke account",
			"--from", alice,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Still in active accounts (not enough approvals to revoke)
		out, err := QueryAllAccountsRaw()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		// In proposed-to-revoke list
		out, err = QueryAllProposedAccountsToRevoke()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryProposedAccountToRevoke(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), aliceAddr)
		require.Contains(t, string(out), "Alice proposes to revoke account")

		// Not yet revoked
		out, err = QueryAllRevokedAccounts()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)
	})

	t.Run("BobApprovesRevoke", func(t *testing.T) {
		txResult, err := utils.ExecuteTx("tx", "auth", "approve-revoke-account",
			"--address", userAddr,
			"--from", bob,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Now revoked
		out, err := QueryAllRevokedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), "TrusteeVoting")

		// No longer in active accounts
		out, err = QueryAllAccountsRaw()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryAllProposedAccountsToRevoke()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")

		out, err = QueryProposedAccountToRevoke(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryRevokedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), "TrusteeVoting")
	})

	t.Run("ReAddAfterRevoke_ProposeApprove", func(t *testing.T) {
		// Jack proposes again
		txResult, err := utils.ExecuteTx("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "NodeAdmin",
			"--from", jack,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Still revoked in revoked list
		out, err := QueryAllRevokedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		// Not yet active
		out, err = QueryAllAccountsRaw()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		// Alice approves
		txResult, err = utils.ExecuteTx("tx", "auth", "approve-add-account",
			"--address", userAddr,
			"--info", "Alice is approving this account",
			"--from", alice,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// No longer in revoked list
		out, err = QueryAllRevokedAccounts()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		// Active again
		out, err = QueryAllAccountsRaw()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
	})

	t.Run("RejectScenario", func(t *testing.T) {
		// Propose-revoke again then approve-revoke
		txResult, err := utils.ExecuteTx("tx", "auth", "propose-revoke-account",
			"--address", userAddr,
			"--info", "Alice proposes to revoke account",
			"--from", alice,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		txResult, err = utils.ExecuteTx("tx", "auth", "approve-revoke-account",
			"--address", userAddr,
			"--from", bob,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Jack proposes again for rejection test
		txResult, err = utils.ExecuteTx("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "NodeAdmin",
			"--from", jack,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Alice rejects
		txResult, err = utils.ExecuteTx("tx", "auth", "reject-add-account",
			"--address", userAddr,
			"--info", "Alice is rejecting this account",
			"--from", alice,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Still in proposed (alice rejected but not enough rejections)
		out, err := QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryAllAccountsRaw()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryAllRejectedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")

		// Bob rejects
		txResult, err = utils.ExecuteTx("tx", "auth", "reject-add-account",
			"--address", userAddr,
			"--info", "Bob is rejecting this account",
			"--from", bob,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// Bob cannot reject again
		txBad, errBad := utils.ExecuteTx("tx", "auth", "reject-add-account",
			"--address", userAddr,
			"--info", "Bob is rejecting this account",
			"--from", bob,
		)
		if errBad == nil {
			require.NotEqual(t, uint32(0), txBad.Code)
		}

		// Now in rejected list (enough rejections)
		out, err = QueryAllRejectedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryAllAccountsRaw()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryRejectedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), jackAddr)
		require.Contains(t, string(out), "Jack is proposing this account")
		require.Contains(t, string(out), aliceAddr)
		require.Contains(t, string(out), "Alice is rejecting this account")
		require.Contains(t, string(out), bobAddr)
		require.Contains(t, string(out), "Bob is rejecting this account")
	})

	// Unused variables referenced to avoid compiler errors
	_ = strings.TrimSpace
}

func TestAuthDemoVendorAccount(t *testing.T) {
	jack := testconstants.JackAccount
	alice := testconstants.AliceAccount

	vid := rand.Intn(65535)

	t.Run("VendorApprovedByOneApproval", func(t *testing.T) {
		name := fmt.Sprintf("vendor%d", rand.Intn(99999))
		err := AddKey(name)
		require.NoError(t, err)

		userAddr, err := GetAddress(name)
		require.NoError(t, err)
		userPubkey, err := GetPubkey(name)
		require.NoError(t, err)

		jackAddr, err := GetAddress(jack)
		require.NoError(t, err)
		aliceAddr, err := GetAddress(alice)
		require.NoError(t, err)

		// Jack proposes Vendor — only needs 1/3 trustee approvals, so jack's proposal is enough
		txResult, err := utils.ExecuteTx("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "Vendor",
			"--vid", fmt.Sprintf("%d", vid),
			"--from", jack,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// With Vendor role, 1 approval is sufficient so account is already active
		out, err := QueryAllAccountsRaw()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), jackAddr)
		require.Contains(t, string(out), "Jack is proposing this account")

		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")

		_ = aliceAddr
	})

	t.Run("VendorWithPidRanges_Success", func(t *testing.T) {
		pid := rand.Intn(65535)
		pidRanges := fmt.Sprintf("%d-%d", pid, pid)
		vidWithPids := vid + 1

		name := fmt.Sprintf("vendor%d", rand.Intn(99999))
		err := AddKey(name)
		require.NoError(t, err)

		userAddr, err := GetAddress(name)
		require.NoError(t, err)
		userPubkey, err := GetPubkey(name)
		require.NoError(t, err)

		jackAddr, err := GetAddress(jack)
		require.NoError(t, err)

		txResult, err := utils.ExecuteTx("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "Vendor",
			"--vid", fmt.Sprintf("%d", vidWithPids),
			"--pid_ranges", pidRanges,
			"--from", jack,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		out, err := QueryAllAccountsRaw()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), jackAddr)

		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), "[]")
	})

	t.Run("VendorWithInvalidPidRanges_Fails", func(t *testing.T) {
		invalidPidRanges := "100-101,1-200"

		name := fmt.Sprintf("vendor%d", rand.Intn(99999))
		err := AddKey(name)
		require.NoError(t, err)

		userAddr, err := GetAddress(name)
		require.NoError(t, err)
		userPubkey, err := GetPubkey(name)
		require.NoError(t, err)

		out, err := utils.ExecuteCLI("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "Vendor",
			"--vid", fmt.Sprintf("%d", vid),
			"--pid_ranges", invalidPidRanges,
			"--from", jack,
			"--yes", "-o", "json", "--keyring-backend", "test",
		)
		// Expect error about invalid PID range
		combined := string(out)
		if err != nil {
			combined = combined + err.Error()
		}
		require.Contains(t, combined, "invalid PID Range is provided")

		out2, _ := QueryProposedAccount(userAddr)
		require.Contains(t, string(out2), "Not Found")

		out3, _ := QueryAllProposedAccounts()
		require.Contains(t, string(out3), "[]")

		out4, _ := QueryAllAccountsRaw()
		require.NotContains(t, string(out4), userAddr)
	})

	t.Run("NodeAdminWithVendorRole_NeedsMoreApprovals", func(t *testing.T) {
		name := fmt.Sprintf("user%d", rand.Intn(99999))
		err := AddKey(name)
		require.NoError(t, err)

		userAddr, err := GetAddress(name)
		require.NoError(t, err)
		userPubkey, err := GetPubkey(name)
		require.NoError(t, err)

		jackAddr, err := GetAddress(jack)
		require.NoError(t, err)
		aliceAddr, err := GetAddress(alice)
		require.NoError(t, err)

		txResult, err := utils.ExecuteTx("tx", "auth", "propose-add-account",
			"--info", "Jack is proposing this account",
			"--address", userAddr,
			"--pubkey", userPubkey,
			"--roles", "Vendor,NodeAdmin",
			"--vid", fmt.Sprintf("%d", vid),
			"--from", jack,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		// NodeAdmin requires 2/3 approval so not yet active
		out, err := QueryAllAccountsRaw()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")

		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		// Alice approves — now has 2 approvals, should become active
		txResult, err = utils.ExecuteTx("tx", "auth", "approve-add-account",
			"--address", userAddr,
			"--info", "Alice is approving this account",
			"--from", alice,
		)
		require.NoError(t, err)
		require.Equal(t, uint32(0), txResult.Code)
		_, err = utils.AwaitTxConfirmation(txResult.TxHash)
		require.NoError(t, err)

		out, err = QueryAllAccountsRaw()
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)

		out, err = QueryAccountRaw(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), userAddr)
		require.Contains(t, string(out), jackAddr)
		require.Contains(t, string(out), aliceAddr)
		require.Contains(t, string(out), "Alice is approving this account")

		out, err = QueryAllProposedAccounts()
		require.NoError(t, err)
		require.NotContains(t, string(out), userAddr)

		out, err = QueryProposedAccount(userAddr)
		require.NoError(t, err)
		require.Contains(t, string(out), "Not Found")
	})
}
