package service

import (
	"encoding/csv"
	"fmt"
	"github.com/maximum-c/finance_dashboard/internal/models"
	"github.com/maximum-c/finance_dashboard/internal/storage"
	"io"
	"strconv"
	"strings"
	"time"
)

type TransactionService struct {
	Storage *storage.TransactionStorage
}

func NewTransactionService(Storage *storage.TransactionStorage) *TransactionService {
	return &TransactionService{Storage: Storage}
}
func validateHeaders(headers []string) (map[string]int, error) {
	required := map[string]bool{
		"date":        false,
		"description": false,
		"amount":      false,
	}

	headerMap := make(map[string]int)

	for i, header := range headers {
		normalized := strings.ToLower(strings.TrimSpace(header))
		headerMap[normalized] = i
		if _, isRequired := required[normalized]; isRequired {
			required[normalized] = true
		}
	}

	for header, found := range required {
		if !found {
			return nil, fmt.Errorf("missing required header: %s", header)
		}
	}
	return headerMap, nil

}
func (s *TransactionService) ImportCSV(file io.Reader, accoundID int64) (int, error) {
	reader := csv.NewReader(file)

	csvHeaders, err := reader.Read()

	if err != nil {
		return 0, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	headerMap, err := validateHeaders(csvHeaders)
	if err != nil {
		return 0, fmt.Errorf("invalid headers: %w", err)
	}

	var transactions []models.Transaction
	lineNumber := 1
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, fmt.Errorf("failed to read line %d: %w", lineNumber, err)
		}

		transaction, err := parseTransaction(record, headerMap, accoundID)
		if err != nil {
			return 0, fmt.Errorf("failed to parse transaction on line %d: %w", lineNumber, err)
		}
		transactions = append(transactions, transaction)
		lineNumber++
	}
	err = s.Storage.AddTransactions(transactions)
	if err != nil {
		return 0, err
	}
	return len(transactions), nil
}
func parseTransaction(record []string, headerMap map[string]int, accountID int64) (models.Transaction, error) {
	//todo implment transaction parsing.
	dateStr := record[headerMap["date"]]
	description := record[headerMap["description"]]
	amountStr := record[headerMap["amount"]]

	date, err := time.Parse("2006-01-02 3:04 PM", dateStr)
	if err != nil {
		return models.Transaction{}, err
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("invalid amount: %s", amountStr)
	}
	var category string
	if categoryIndex, exists := headerMap["category"]; exists && categoryIndex < len(record) {
		category = record[categoryIndex]
	}
	return models.Transaction{
		Date:        date,
		Description: description,
		Amount:      amount,
		Category:    category,
		AccountID:   accountID,
		CreatedAt:   time.Now(),
	}, nil

}

func (s *TransactionService) FetchTransactionsWithFilter(filters models.TransactionFilter) ([]models.Transaction, error) {
	return s.Storage.GetTransactions(filters)
}

func (s *TransactionService) FetchTransactionStats(filters models.TransactionFilter) (map[string]float64, error) {
	return s.Storage.GetTransactionsStats(filters)
}
