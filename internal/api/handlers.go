package api

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maximum-c/finance_dashboard/internal/models"
	"github.com/maximum-c/finance_dashboard/internal/service"
)

type Handler struct {
	service *service.FinanceService
}

func (h *Handler) UploadCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to get file from form",
			"descripton": err.Error(),
		})
		return
	}
	defer file.Close()

	if header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Empty File",
		})
		return
	}

	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File must be CSV",
		})
		return
	}

	reader := csv.NewReader(file)

	csvHeaders, err := reader.Read()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Failed to read CSV headers",
			"description": err.Error(),
		})
		return
	}

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
}
func parseTransaction(record []string, headers []string) (models.Transaction, error) {
	//todo implment transaction parsing.
	return models.Transaction{}, nil
}
