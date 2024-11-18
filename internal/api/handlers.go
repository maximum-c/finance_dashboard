package api

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/maximum-c/finance_dashboard/internal/models"
	"github.com/maximum-c/finance_dashboard/internal/service"
)

type Handler struct {
	service *service.FinanceService
}

func (h *Handler) UploadCSV(c *gin.Context) {
	
	accountID, err := strconv.ParseInt(c.Param("accountID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
	}

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

c.JSON(http.StatusOK, gin.H{
		"message": "successfully imported transactions",
		"count":   len(transactions),
	})
}
func parseTransaction(record []string, headers []string) (models.Transaction, error) {
	//todo implment transaction parsing.
	var t models.Transaction
	tm, err := time.Parse("01/02/2006 3:04PM PST", record[0])
	if err != nil {
		return models.Transaction{}, err
	}
	t.Date = tm
	t.Description = record[1]
	amount, err := strconv.ParseFloat(record[2])
	if err != nil {
		return models.Transaction{}, err
	}
	t.Amount = amount
	t.AccountID = 
	return models.Transaction{}, nil
}
