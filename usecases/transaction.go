package usecases

import "github.com/carlosmaniero/budgetgo-ofximporter/domain"

// TransactionIterator contains all transactions usecases
type TransactionIterator struct {
	Repository TransactionRepository
}

// Register one transaction inside the repository
func (iterator *TransactionIterator) Register(transaction *domain.Transaction) error {
	return iterator.Repository.Store(transaction)
}

// TransactionRepository is the place where transactions are stored
type TransactionRepository interface {
	Store(*domain.Transaction) error
}
