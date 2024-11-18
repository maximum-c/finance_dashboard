package storage

import (
	"database/sql"
	"strings"

	"github.com/maximum-c/finance_dashboard/internal/models"
)

type TransactionStorage struct {
	db *sql.DB
}

func InitDB(dbPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS transactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date DATETIME NOT NULL,
            description TEXT NOT NULL,
            amount REAL NOT NULL,
            category TEXT,
            account_id INTEGER NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE TABLE IF NOT EXISTS accounts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            type TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
	}
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func NewTransactionStorage(db *sql.DB) *TransactionStorage {
	return &TransactionStorage{db: db}
}

func (s *TransactionStorage) CreateTransaction(t *models.Transaction) error {
	query := `
	INSERT INTO transactions (date, description, amount, account_id, created_at)
	VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, t.Date, t.Description, t.Amount, t.AccountID, t.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id
	return nil
}

func (s *TransactionStorage) GetTransactions(filter models.TransactionFilter) ([]models.Transaction, error) {
	query := `
		SELECT id, date, description, amount, category, account_id, created_at
		FROM transactions WHERE 1=1`
	var conditions []string
	var args []interface{}

	if filter.StartDate != nil {
		conditions = append(conditions, "date >= ?")
		args = append(args, filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "date <= ?")
		args = append(args, filter.EndDate)
	}

	if filter.Category != nil {
		conditions = append(conditions, "category = ?")
		args = append(args, filter.Category)
	}

	if filter.AccountID != nil {
		conditions = append(conditions, "account_id = ?")
		args = append(args, filter.AccountID)
	}

	for len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += "ORDER BY date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Description,
			&t.Category,
			&t.AccountID,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
