package domain

import "time"

// Transaction is the transaction representation
type Transaction struct {
	ID        string
	Name      string
	Amount    float64
	Date      time.Time
	FundingID string
}
