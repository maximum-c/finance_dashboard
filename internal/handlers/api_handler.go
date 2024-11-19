package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/maximum-c/finance_dashboard/internal/models"
	"github.com/maximum-c/finance_dashboard/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

type APIHandler struct {
	TransactionService *service.TransactionService
}

func NewAPIHandler(service *service.TransactionService) *APIHandler {
	return &APIHandler{
		TransactionService: service,
	}
}
func parseFilters(filters *models.TransactionFilter, c *gin.Context) {

	if startDate := c.Query("startDate"); startDate != "" {
		parsedDate, err := time.Parse("2006-01-02 3:04 PM", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Invalid Start Date format",
				"description": err.Error(),
			})
			log.Printf("invalid start date: %v\n", startDate)
			return
		}
		filters.StartDate = &parsedDate
	}

	if endDate := c.Query("endDate"); endDate != "" {
		parsedDate, err := time.Parse("2006-01-02 3:04 PM", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":       "Invalid End Date format",
				"description": err.Error(),
			})
			log.Printf("invalid end date: %v\n", endDate)
			return
		}
		filters.EndDate = &parsedDate
	}

	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}

	if accountID := c.Query("accountID"); accountID != "" {
		parsedID, err := strconv.ParseInt(accountID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
			log.Printf("invalid account id: %v\n", accountID)
			return
		}
		filters.AccountID = &parsedID
	}
}

func (h *APIHandler) GetTransactions(c *gin.Context) {
	var filters models.TransactionFilter
	parseFilters(&filters, c)

	transactions, err := h.TransactionService.FetchTransactionsWithFilter(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to Fetch Transactions",
			"description": err.Error(),
		})
		log.Print("failed to fetch transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
}

func (h *APIHandler) GetTransactionStats(c *gin.Context) {
	var filters models.TransactionFilter
	parseFilters(&filters, c)

	stats, err := h.TransactionService.FetchTransactionStats(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "Failed to fetch transaction statistics",
			"description": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})

}
