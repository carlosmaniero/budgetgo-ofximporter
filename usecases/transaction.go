package usecases

import (
	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
)

// TransactionIterator contains all transactions usecases
type TransactionIterator struct {
	Repository       TransactionRepository
	FundingID        string
	ConcurrencyLevel int
}

func (iterator *TransactionIterator) getConcurrentyLevel() int {
	if iterator.ConcurrencyLevel == 0 {
		return 100
	}

	return iterator.ConcurrencyLevel
}

// Register one transaction inside the repository
func (iterator *TransactionIterator) Register(transaction *domain.Transaction) error {
	return iterator.Repository.Store(transaction)
}

func (iterator *TransactionIterator) registerAsync(total int, transactions TransactionList, sent chan bool) {
	for i := 0; i < total; i++ {
		go func() {
			transaction := domain.Transaction{}

			if err := transactions.Next(&transaction); err == nil {
				transaction.FundingID = iterator.FundingID
				iterator.Register(&transaction)
			}

			sent <- true
		}()
	}
}

func (iterator *TransactionIterator) waitForTransactions(total int, sent chan bool) {
	totalSented := 0
	for {
		<-sent
		totalSented++

		if totalSented == total {
			return
		}
	}
}

// RegisterMany will register many transactions
func (iterator *TransactionIterator) RegisterMany(list TransactionList) {
	sent := make(chan bool, iterator.getConcurrentyLevel())
	totalSent := 0

	for list.HasNext() {
		toSend := iterator.getConcurrentyLevel()

		if toSend > list.Remaining() {
			toSend = list.Remaining()
		}

		iterator.registerAsync(toSend, list, sent)
		iterator.waitForTransactions(toSend, sent)
		totalSent += iterator.getConcurrentyLevel()
	}
}

// TransactionRepository is the place where transactions are stored
type TransactionRepository interface {
	Store(*domain.Transaction) error
}

// TransactionList is a dynamic list of transactions
type TransactionList interface {
	HasNext() bool
	Next(*domain.Transaction) error
	Count() int
	Remaining() int
}
