package importer

import (
	"errors"
	"io"

	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
	"github.com/carlosmaniero/ofx"
)

// TransactionOfxImporter is an importer of ofx
type TransactionOfxImporter struct {
	File io.Reader
}

// ErrNoMoreTransaction is the error returned by the next function when
// there is no more transactions to iterate over.
var ErrNoMoreTransaction = errors.New("there is no more transactions")

// Parse return a transactionIterator that iterates into a transaction list
func (importer *TransactionOfxImporter) Parse() *TransactionIterator {
	parsedOfx, err := ofx.Parse(importer.File)

	if err != nil {
		panic(err)
	}

	transactions := parsedOfx.Transactions

	return &TransactionIterator{transactions: transactions}
}

// TransactionIterator iterates into a transaction list return a domain transaction
type TransactionIterator struct {
	transactions []*ofx.Transaction
	current      int
}

// Next returns the next transaction
func (iterator *TransactionIterator) Next(transaction *domain.Transaction) error {
	if iterator.current >= len(iterator.transactions) {
		return ErrNoMoreTransaction
	}

	iterator.ofxToTransaction(transaction, iterator.transactions[iterator.current])
	iterator.current++

	return nil
}

// HasNext checks if has next
func (iterator *TransactionIterator) HasNext() bool {
	return iterator.current < len(iterator.transactions)
}

// Count returns the total of transactions
func (iterator *TransactionIterator) Count() int {
	return len(iterator.transactions)
}

func (iterator *TransactionIterator) ofxToTransaction(transaction *domain.Transaction, ofxTransaction *ofx.Transaction) {
	transaction.Description = ofxTransaction.Description
	if ofxTransaction.Memo != "" {
		transaction.Description = ofxTransaction.Memo
	}
	transaction.Amount, _ = ofxTransaction.Amount.Value.Float64()
	transaction.Date = ofxTransaction.PostedDate
}
