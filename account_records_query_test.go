package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeAccountRecordsQuery(t *testing.T) {
	query := NewAccountRecordsQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptoGetAccountRecords:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestAccountRecordQuery_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(env.Client)
	assert.NoError(t, err)

	recordsQuery, err := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(env.Client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := records.GetCost(env.Client)
	assert.NoError(t, err)

	recordsQuery, err := records.SetQueryPayment(cost).Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_BigMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(env.Client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(100000))

	_, err = records.GetCost(env.Client)
	assert.NoError(t, err)

	recordsQuery, err := records.Execute(env.Client)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_SmallMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(env.Client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetMaxQueryPayment(HbarFromTinybar(1))

	cost, err := records.GetCost(env.Client)
	assert.NoError(t, err)

	recordsQuery, err := records.Execute(env.Client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("cost of AccountRecordsQuery ("+cost.String()+") without explicit payment is greater than the max query payment of 1 tħ"), err.Error())
	}

	assert.Equal(t, 0, len(recordsQuery))
}

func TestAccountRecordQueryCost_InsufficientFee_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	newKey, err := GeneratePrivateKey()
	assert.NoError(t, err)

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey.PublicKey()).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetInitialBalance(NewHbar(1)).
		Execute(env.Client)
	assert.NoError(t, err)

	receipt, err := resp.GetReceipt(env.Client)
	assert.NoError(t, err)

	account := *receipt.AccountID

	_, err = NewTransferTransaction().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		AddHbarTransfer(account, NewHbar(1)).
		AddHbarTransfer(env.Client.GetOperatorAccountID(), NewHbar(-1)).
		Execute(env.Client)
	assert.NoError(t, err)

	records := NewAccountRecordsQuery().
		SetNodeAccountIDs([]AccountID{resp.NodeID}).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err = records.GetCost(env.Client)
	assert.NoError(t, err)

	recordsQuery, err := records.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INSUFFICIENT_TX_FEE"), err.Error())
	}

	assert.Equal(t, 0, len(recordsQuery))
}

func Test_AccountRecord_NoAccountID(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountRecordsQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status INVALID_ACCOUNT_ID"), err.Error())
	}
}
