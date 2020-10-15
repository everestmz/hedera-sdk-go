package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// AccountCreateTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountCreateTransaction struct {
	Transaction
	pb *proto.CryptoCreateTransactionBody
}

// NewAccountCreateTransaction creates an AccountCreateTransaction transaction which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	pb := &proto.CryptoCreateTransactionBody{}

	transaction := AccountCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	transaction.SetReceiveRecordThreshold(MaxHbar)
	transaction.SetSendRecordThreshold(MaxHbar)

	return &transaction
}

// SetKey sets the key that must sign each transfer out of the account. If RecieverSignatureRequired is true, then it
// must also sign any transfer into the account.
func (transaction *AccountCreateTransaction) SetKey(publicKey PublicKey) *AccountCreateTransaction {
	transaction.pb.Key = publicKey.toProtoKey()
	return transaction
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *AccountCreateTransaction) SetInitialBalance(initialBalance Hbar) *AccountCreateTransaction {
	transaction.pb.InitialBalance = uint64(initialBalance.AsTinybar())
	return transaction
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *AccountCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *AccountCreateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return transaction
}

// SetSendRecordThreshold sets the threshold amount for which an account record is created for any send/withdraw
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction *AccountCreateTransaction) SetSendRecordThreshold(recordThreshold Hbar) *AccountCreateTransaction {
	transaction.pb.SendRecordThreshold = uint64(recordThreshold.AsTinybar())
	return transaction
}

// SetReceiveRecordThreshold sets the threshold amount for which an account record is created for any receive/deposit
// transaction
//
// Deprecated: No longer used by Hedera
func (transaction *AccountCreateTransaction) SetReceiveRecordThreshold(recordThreshold Hbar) *AccountCreateTransaction {
	transaction.pb.ReceiveRecordThreshold = uint64(recordThreshold.AsTinybar())
	return transaction
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a node, then it will behave as if proxyAccountID was not set.
func (transaction *AccountCreateTransaction) SetProxyAccountID(id AccountID) *AccountCreateTransaction {
	transaction.pb.ProxyAccountID = id.toProtobuf()
	return transaction
}

// SetReceiverSignatureRequired sets the receiverSigRequired flag. If the receiverSigRequired flag is set to true, then
// all cryptocurrency transfers must be signed by this account's key, both for transfers in and out. If it is false,
// then only transfers out have to be signed by it. This transaction must be signed by the
// payer account. If receiverSigRequired is false, then the transaction does not have to be signed by the keys in the
// keys field. If it is true, then it must be signed by them, in addition to the keys of the payer account.
func (transaction *AccountCreateTransaction) SetReceiverSignatureRequired(required bool) *AccountCreateTransaction {
	transaction.pb.ReceiverSigRequired = required
	return transaction
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func accountCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CreateAccount,
	}
}

func (transaction *AccountCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountCreateTransaction) Sign(
	privateKey PrivateKey,
) *AccountCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *AccountCreateTransaction) SignWithOperator(
	client *Client,
) (*AccountCreateTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *AccountCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountCreateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountCreateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
			client.operator.signer,
		)
	}

	_, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		accountCreateTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
		query_makePaymentTransaction,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *AccountCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoCreateAccount{
		CryptoCreateAccount: transaction.pb,
	}

	return true
}

func (transaction *AccountCreateTransaction) Freeze() (*AccountCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountCreateTransaction) FreezeWith(client *Client) (*AccountCreateTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *AccountCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetMaxTransactionFee(fee Hbar) *AccountCreateTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *AccountCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionMemo(memo string) *AccountCreateTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionValidDuration(duration time.Duration) *AccountCreateTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetTransactionID(transactionID TransactionID) *AccountCreateTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *AccountCreateTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this AccountCreateTransaction.
func (transaction *AccountCreateTransaction) SetNodeID(nodeID AccountID) *AccountCreateTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
