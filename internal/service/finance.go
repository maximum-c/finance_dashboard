package service

import (
	"database/sql"
	//"time"

	"github.com/maximum-c/finance_dashboard/internal/models"
)

type FinanceService struct {
	db *sql.DB
}

func NewFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{db: db}
}

func (s *FinanceService) ImportTransactions(transactions []models.Transaction) error {
	// todo implement
	return nil
}

func (s *FinanceService) GetTransactions(filter models.TransactionFilter) ([]models.Transaction, error) {
	return nil, nil
}

func (s *FinanceService) GetAccountBalace(accountID int64) (float64, error) {
	return 0, nil
}
