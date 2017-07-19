package usecases

import (
	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
)

// TransactionIterator contains all transactions usecases
type TransactionIterator struct {
	Repository       TransactionRepository
	ConcurrencyLevel int
}

// Register one transaction inside the repository
func (iterator *TransactionIterator) Register(transaction *domain.Transaction) error {
	return iterator.Repository.Store(transaction)
}

// RegisterAsync will register many transactions
func (iterator *TransactionIterator) RegisterAsync(list TransactionList, fundingID string, concorrencyLevel int) {
	control := transactionRegisterAsyncControl{
		iterator:  iterator,
		list:      list,
		fundingID: fundingID,
	}
	control.start()
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

type transactionRegisterAsyncControl struct {
	list             TransactionList
	fundingID        string
	sent             chan bool
	completed        int
	concurrencyLevel int
	iterator         *TransactionIterator
}

func (control *transactionRegisterAsyncControl) register(processes int) {
	for i := 0; i < processes; i++ {
		go func() {
			transaction := domain.Transaction{}

			if err := control.list.Next(&transaction); err == nil {
				transaction.FundingID = control.fundingID
				if err := control.iterator.Register(&transaction); err != nil {
					panic(err)
				}
			}

			control.sent <- true
		}()
	}
}

func (control *transactionRegisterAsyncControl) getConcurrentyLevel() int {
	if count := control.list.Count(); count < control.concurrencyLevel {
		return count
	}
	if control.concurrencyLevel == 0 {
		return 100
	}
	return control.concurrencyLevel
}

func (control *transactionRegisterAsyncControl) isFinished() bool {
	return control.completed == control.list.Count()
}

func (control *transactionRegisterAsyncControl) start() {
	control.sent = make(chan bool, control.getConcurrentyLevel())
	control.register(control.getConcurrentyLevel())

	for {
		<-control.sent
		control.completed++

		if control.isFinished() {
			return
		} else if control.list.HasNext() {
			control.register(1)
		}
	}
}
