package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maximum-c/finance_dashboard/internal/handlers"
	"github.com/maximum-c/finance_dashboard/internal/service"
	"github.com/maximum-c/finance_dashboard/internal/storage"
	"log"
)

func main() {
	dbPath := "./data/finance.db"
	if err := storage.InitDB(dbPath); err != nil {
		log.Fatal("Failed to Initialize Database:", err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open Database:", err)
	}
	defer db.Close()

	transactionStorage := storage.NewTransactionStorage(db)
	TransactionService := service.NewTransactionService(transactionStorage)
	apiHandler := handlers.NewAPIHandler(TransactionService)
	csvHandler := handlers.NewCSVHandler(TransactionService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	api := r.Group("/api")
	{
		api.GET("/transactions/:accountID", apiHandler.GetTransactions)

		api.GET("/transactions/stats", apiHandler.GetTransactionStats)
	}

	r.POST("/upload/:accountID", csvHandler.UploadCSV)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
