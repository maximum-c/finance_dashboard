package models

import (
	"time"
)

type Transaction struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	AccountID   int64     `json:"account_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Account struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
