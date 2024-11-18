package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maximum-c/finance_dashboard/internal/service"
)

type Handler struct {
	Service *service.FinanceService
}

func (h *Handler) UploadCSV(c *gin.Context) {

	accountID, err := strconv.ParseInt(c.Param("accountID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		log.Printf("Error parsing Account ID: %v\n", err)
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Failed to get file from form",
			"descripton": err.Error(),
		})
		return
	}
	if header.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Empty File",
		})
		return
	}

	const maxFileSize = 10 * 1024 * 1024 // 10 MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds the limit of 10 MB.",
		})
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File must be CSV",
		})
		return
	}
	count, err := h.Service.ImportCSV(file, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "successfully imported transactions",
		"filename": header.Filename,
		"count":    count,
	})
}
