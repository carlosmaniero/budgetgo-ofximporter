package usecases

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecCreateTransaction(t *testing.T) {
	Convey("Scenario: Creating a transaction", t, func() {
		Convey("Given I've a valid transaction", func() {
			transaction := domain.Transaction{
				Description: "Golpher Shampoo",
				Amount:      -10,
				Date:        time.Now(),
				FundingID:   "pet-account",
			}

			Convey("When I create the transaction", func() {
				repository := mockTransactionRepository{count: 0}
				iterator := TransactionIterator{
					Repository: &repository,
				}
				err := iterator.Register(&transaction)

				Convey("Then the transaction is stored without errors", func() {
					So(err, ShouldBeNil)
					So(repository.count, ShouldEqual, 1)
				})
			})
		})
	})
	Convey("Scenario: Creating many transaction", t, func() {
		Convey("Given I've a valid transaction List", func() {
			transaction := domain.Transaction{
				Description: "Golpher Shampoo",
				Amount:      -10,
				Date:        time.Now(),
				FundingID:   "pet-account",
			}

			items := make([]*domain.Transaction, 0)

			for i := 0; i < 123; i++ {
				items = append(items, &transaction)
			}

			transactionList := SimpleTransactionList{
				items: items,
			}

			Convey("When I register many", func() {
				repository := mockTransactionRepository{count: 0}
				iterator := TransactionIterator{
					Repository: &repository,
				}
				iterator.RegisterAsync(&transactionList, "pet-account", 100)

				Convey("Then all transactions was got", func() {
					So(transactionList.Remaining(), ShouldEqual, 0)
				})

				Convey("And all transactions was registred", func() {
					So(repository.count, ShouldEqual, 123)
				})
			})
		})
	})
}

type SimpleTransactionList struct {
	items   []*domain.Transaction
	current int
	mux     sync.Mutex
}

// Next returns the next transaction
func (iterator *SimpleTransactionList) Next(transaction *domain.Transaction) error {
	iterator.mux.Lock()
	if !iterator.HasNext() {
		iterator.mux.Unlock()
		return errors.New("No more transactions")
	}

	transaction = iterator.items[iterator.current]
	iterator.current++
	iterator.mux.Unlock()

	return nil
}

// HasNext checks if has next
func (iterator *SimpleTransactionList) HasNext() bool {
	return iterator.current < len(iterator.items)
}

// Count returns the total of transactions
func (iterator *SimpleTransactionList) Count() int {
	return len(iterator.items)
}

// Remaining returns the total of transactions
func (iterator *SimpleTransactionList) Remaining() int {
	iterator.mux.Lock()
	remaining := iterator.Count() - iterator.current
	iterator.mux.Unlock()
	return remaining
}

type mockTransactionRepository struct {
	count int
}

func (repository *mockTransactionRepository) Store(*domain.Transaction) error {
	repository.count++
	return nil
}
