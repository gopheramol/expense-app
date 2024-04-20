// main.go

package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/gopheramol/expense-app/db"
	"github.com/gopheramol/expense-app/handler"

	"github.com/gin-gonic/gin"
	"github.com/gopheramol/expense-app/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://db:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	expenseStore, err := db.NewMongoStore(client, "expenseTracker", "expenses")
	if err != nil {
		log.Fatal(err)
	}

	expenseService := service.NewExpenseService(expenseStore)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	router := gin.Default()

	templatesDir := filepath.Join(".", "templates")
	router.LoadHTMLGlob(filepath.Join(templatesDir, "*.html"))

	router.GET("/", expenseHandler.HandleGetExpenses)
	router.POST("/expenses", expenseHandler.HandleCreateExpense)

	router.Run(":8080")
}
