package repositories

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
)

type transactionBody struct {
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	FundingID   string    `json:"funding_id"`
}

func (data *transactionBody) loads(transaction *domain.Transaction) {
	data.Description = transaction.Description
	data.Amount = transaction.Amount
	data.Date = transaction.Date
	data.FundingID = transaction.FundingID
}

// TransactionRestRepository is the rest repository of transaction
type TransactionRestRepository struct {
	Server string
	Client WebClient
}

func (repository *TransactionRestRepository) endpoint() string {
	return repository.Server + "transaction"
}

func (repository *TransactionRestRepository) prepareBody(transaction *domain.Transaction) io.Reader {
	data := transactionBody{}
	data.loads(transaction)
	body, _ := json.Marshal(data)
	return bytes.NewReader(body)
}

// Store a transaction into the endpoint
func (repository *TransactionRestRepository) Store(transaction *domain.Transaction) error {
	request, _ := http.NewRequest("POST", repository.endpoint(), repository.prepareBody(transaction))
	response, err := repository.Client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 201 {
		if b, err := ioutil.ReadAll(response.Body); err == nil {
			panic(strconv.Itoa(response.StatusCode) + " status - " + string(b))
		} else {
			panic(err)
		}
	}

	return nil
}

// WebClient interface
//
// This is a simple client that sends a request to anywhere
type WebClient interface {
	Do(*http.Request) (*http.Response, error)
}

// NewTransactionRepository create a new repository
func NewTransactionRepository(server string) *TransactionRestRepository {
	return &TransactionRestRepository{
		Server: server,
		Client: &http.Client{},
	}
}
