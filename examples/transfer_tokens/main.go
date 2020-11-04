package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"os"
)

func main() {
	client := hedera.ClientForPreviewnet()

	configOperatorID := os.Getenv("OPERATOR_ID")
	configOperatorKey := os.Getenv("OPERATOR_KEY")

	var err error

	if configOperatorID != "" && configOperatorKey != "" {
		operatorAccountID, err := hedera.AccountIDFromString(configOperatorID)
		if err != nil {
			panic(err)
		}

		operatorPrivateKey, err := hedera.PrivateKeyFromString(configOperatorKey)
		if err != nil {
			panic(err)
		}

		client.SetOperator(operatorAccountID, operatorPrivateKey)
	}

	key1, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	key2, err := hedera.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	fmt.Printf("privateKey = %v\n", key1.String())
	fmt.Printf("publicKey = %v\n", key1.PublicKey().String())
	fmt.Printf("privateKey = %v\n", key2.String())
	fmt.Printf("publicKey = %v\n", key2.PublicKey().String())

	transactionResponse, err := hedera.NewAccountCreateTransaction().
		SetKey(key1.PublicKey()).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err := transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	accountID1 := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", accountID1.String())

	transactionResponse, err = hedera.NewAccountCreateTransaction().
		SetKey(key1.PublicKey()).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	accountID2 := *transactionReceipt.AccountID

	fmt.Printf("account = %v\n", accountID2.String())

	transactionResponse, err = hedera.NewTokenCreateTransaction().
		SetName("ffff").
		SetSymbol("F").
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetDecimals(3).
		SetInitialSupply(1000000).
		SetTreasury(client.GetOperatorID()).
		SetAdminKey(client.GetOperatorKey()).
		SetFreezeKey(client.GetOperatorKey()).
		SetWipeKey(client.GetOperatorKey()).
		SetKycKey(client.GetOperatorKey()).
		SetSupplyKey(client.GetOperatorKey()).
		SetFreezeDefault(false).
		SetMaxTransactionFee(hedera.NewHbar(1000)).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	tokenID := *transactionReceipt.TokenID

	fmt.Printf("token = %v\n", tokenID.String())

	transaction, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(accountID1).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	transactionResponse, err = transaction.
		Sign(key1).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Associated account %v with token %v\n", accountID1.String(), tokenID.String())

	transaction, err = hedera.NewTokenAssociateTransaction().
		SetAccountID(accountID2).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	transactionResponse, err = transaction.
		Sign(key2).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Associated account %v with token %v\n", accountID2.String(), tokenID.String())

	transactionResponse, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID1).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID1.String(), tokenID.String())

	transactionResponse, err = hedera.NewTokenGrantKycTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID2).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Granted KYC for account %v on token %v\n", accountID2.String(), tokenID.String())

	transactionResponse, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, client.GetOperatorID(), -10).
		AddTokenTransfer(tokenID, accountID1, 10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Sent 10 tokens from account %v to account %v on token %v\n",
		client.GetOperatorID().String(),
		accountID1.String(),
		tokenID.String(),
	)

	transactionResponse, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID1, -10).
		AddTokenTransfer(tokenID, accountID2, 10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Sent 10 tokens from account %v to account %v on token %v\n",
		accountID1.String(),
		accountID2.String(),
		tokenID.String(),
	)

	transactionResponse, err = hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, accountID2, -10).
		AddTokenTransfer(tokenID, accountID1, 10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Sent 10 tokens from account %v to account %v on token %v\n",
		accountID2.String(),
		accountID1.String(),
		tokenID.String(),
	)

	transactionResponse, err = hedera.NewTokenWipeTransaction().
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetAccountID(accountID1).
		SetTokenID(tokenID).
		SetAmount(10).
		Execute(client)
	if err != nil {
		panic(err)
	}

	_, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wiped account %v on token %v\n", accountID1.String(), tokenID.String())

	transactionResponse, err = hedera.NewTokenDeleteTransaction().
		SetTokenID(tokenID).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted token %v\n", tokenID.String())

	accountDeleteTx, err := hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID1).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorID()).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	transactionResponse, err = accountDeleteTx.
		Sign(key1).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted account %v\n", accountID1.String())

	accountDeleteTx, err = hedera.NewAccountDeleteTransaction().
		SetAccountID(accountID2).
		SetNodeAccountIDs([]hedera.AccountID{transactionResponse.NodeID}).
		SetTransferAccountID(client.GetOperatorID()).
		FreezeWith(client)
	if err != nil {
		panic(err)
	}

	transactionResponse, err = accountDeleteTx.
		Sign(key2).
		Execute(client)
	if err != nil {
		panic(err)
	}

	transactionReceipt, err = transactionResponse.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Deleted account %v\n", accountID2.String())
}
