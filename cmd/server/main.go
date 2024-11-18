package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	log.Fatal(r.Run(":8080"))
}
