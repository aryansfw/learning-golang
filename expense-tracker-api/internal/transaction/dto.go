package transaction

import "time"

type TransactionFilters struct {
	Category *bool      `json:"category"`
	DateFrom *time.Time `json:"date_from"`
	DateTo   *time.Time `json:"date_to"`
}
