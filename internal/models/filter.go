package models

import "time"

type TransactionFilter struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Category  *string    `json:"category,omitempty"`
	AccountID *int64     `json:"account_id,omitempty"`
}
