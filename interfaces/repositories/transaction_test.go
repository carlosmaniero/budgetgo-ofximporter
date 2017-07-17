package repositories

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecTransactionRepository(t *testing.T) {
	Convey("Scenario: consulting an endpoint", t, func() {
		Convey("Given I've a repository with localhost server", func() {
			client := mockClient{count: 0}
			repository := TransactionRestRepository{
				Server: "http://localhost/",
				Client: &client,
			}
			Convey("When I check get the endpoint", func() {
				endpoint := repository.endpoint()
				Convey("Then it returns the correct endpoint", func() {
					So(endpoint, ShouldEqual, "http://localhost/transaction")
				})
			})

			Convey("When I send a transaction", func() {
				repository.Store(&domain.Transaction{})
				Convey("Then the transaction is sent by the web client", func() {
					So(client.count, ShouldEqual, 1)
				})
			})
		})
	})
}

type mockClient struct {
	count int
}

func (client *mockClient) Do(*http.Request) (*http.Response, error) {
	client.count++
	response := http.Response{
		StatusCode: 201,
		Body:       ioutil.NopCloser(strings.NewReader("ok")),
	}
	return &response, nil
}
