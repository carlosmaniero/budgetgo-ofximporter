package usecases

import (
	"testing"
	"time"

	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecCreateTransaction(t *testing.T) {
	Convey("Scenario: Creating a transaction", t, func() {
		Convey("Given I've a valid transaction", func() {
			transaction := domain.Transaction{
				Name:      "Golpher Shampoo",
				Amount:    -10,
				Date:      time.Now(),
				FundingID: "pet-account",
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
}

type mockTransactionRepository struct {
	count int
}

func (repository *mockTransactionRepository) Store(*domain.Transaction) error {
	repository.count++
	return nil
}
