package hedera

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSerializeStakersQuery(t *testing.T) {
	query := NewAccountStakersQuery().
		SetAccountID(AccountID{Account: 3}).
		Query

	assert.Equal(t, `cryptoGetProxyStakers:{header:{}accountID:{accountNum:3}}`, strings.ReplaceAll(query.pb.String(), " ", ""))
}

func TestAccountStakersQuery_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountStakersQuery().
		SetAccountID(env.Client.GetOperatorAccountID()).
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetQueryPayment(HbarFromTinybar(25)).
		Execute(env.Client)
	assert.Error(t, err)
}

func TestAccountStakersQueryCost_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)
}

func TestAccountStakersQueryCost_BigMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(100000)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)
}

func TestAccountStakersQueryCost_SmallMax_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(HbarFromTinybar(25)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	cost, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(cost).Execute(env.Client)
	assert.Error(t, err)
}

func TestAccountStakersQueryCost_InsufficientFee_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	accountStakers := NewAccountStakersQuery().
		SetMaxQueryPayment(NewHbar(1)).
		SetNodeAccountIDs(env.NodeAccountIDs).
		SetAccountID(env.Client.GetOperatorAccountID())

	_, err := accountStakers.GetCost(env.Client)
	assert.Error(t, err)

	_, err = accountStakers.SetQueryPayment(HbarFromTinybar(1)).Execute(env.Client)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}
}

func TestAccountStakersNoAccountID_Execute(t *testing.T) {
	env := NewIntegrationTestEnv(t)

	_, err := NewAccountStakersQuery().
		SetNodeAccountIDs(env.NodeAccountIDs).
		Execute(env.Client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, fmt.Sprintf("exceptional precheck status NOT_SUPPORTED"), err.Error())
	}
}
