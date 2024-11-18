package service

import (
	"database/sql"
	"fmt"
	//"time"
	"encoding/csv"
	"io"

	"github.com/maximum-c/finance_dashboard/internal/models"
)

type FinanceService struct {
	db *sql.DB
}

func NewFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{db: db}
}

func (s *FinanceService) ImportCSV(file io.Reader, accoundID int64) (int64, error) {
	reader := csv.NewReader(file)

	csvHeaders, err := reader.Read()

	if err != nil {
		return 0, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	headerMap, err := validateHeaders(csvHeaders)

	expectedHeaders := map[string]bool{
		"date":        false,
		"description": false,
		"amount":      false,
	}

	for _, col := range csvHeaders {
		if _, exists := expectedHeaders[strings.ToLower(col)]; exists {
			expectedHeaders[strings.ToLower(col)] = true
		}
	}

	for col, found := range expectedHeaders {
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Missing Required Column",
				"description": fmt.Sprintf("Column '%s' not found", col),
			})
		}
	}

	var transactions []models.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Failed to read CSV record",
				"description": err.Error(),
			})
			return
		}

		transaction, err := parseTransaction(record, csvHeaders)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Failed to parse transaction",
				"description": err.Error(),
			})
			return
		}
		transactions = append(transactions, transaction)
	}
	if err := h.service.ImportTransactions(transactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to import transactions",
			"description": err.Error(),
		})
		return
	}
}

func validateHeaders(headers []string) map[string]int

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
