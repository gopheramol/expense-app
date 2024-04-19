package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Expense struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Amount    float64            `bson:"amount"`
	CreatedAt string             `bson:"createdAt"`
}

var (
	expensesCollection *mongo.Collection
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Create or get the expenses collection
	expensesCollection = client.Database("expenseTracker").Collection("expenses")
	err = createCollectionIfNotExists(expensesCollection)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// Set the template directory
	templatesDir := filepath.Join(".", "templates")
	router.LoadHTMLGlob(filepath.Join(templatesDir, "*.html"))

	// Serve static files
	staticDir := "./static"
	router.Static("/static", staticDir)

	// Define a route for the home page
	router.GET("/", func(c *gin.Context) {
		expenses, err := getExpenses()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error retrieving expenses")
			return
		}

		total := calculateTotalAmount(expenses)
		c.HTML(http.StatusOK, "index.html", gin.H{"expenses": expenses, "total": total})
	})

	// Define a route for creating expenses
	router.POST("/expenses", func(c *gin.Context) {
		name := c.PostForm("name")
		amountStr := c.PostForm("amount")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "Bad Request")
			return
		}

		t := time.Now().Format("2006-01-02 3:04: PM")

		expense := Expense{Name: name, Amount: amount, CreatedAt: t}
		err = createExpense(&expense)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating expense")
			return
		}

		c.Redirect(http.StatusSeeOther, "/")
	})

	// Start the server
	router.Run(":8080")
}

func createExpense(expense *Expense) error {
	_, err := expensesCollection.InsertOne(context.Background(), expense)
	return err
}

func getExpenses() ([]Expense, error) {
	// Define a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the options for sorting
	findOptions := options.Find().SetSort(bson.D{{"createdAt", -1}}) // Sort by createdAt field in descending order

	// Find all documents in the collection, sorted by createdAt field in descending order
	cursor, err := expensesCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Initialize a slice to hold the expenses
	var expenses []Expense

	// Iterate over the cursor and decode each document
	for cursor.Next(ctx) {
		var expense Expense
		if err := cursor.Decode(&expense); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func calculateTotalAmount(expenses []Expense) float64 {
	total := 0.0
	for _, expense := range expenses {
		total += expense.Amount
	}
	return total
}

func createCollectionIfNotExists(collection *mongo.Collection) error {
	// Check if the collection already exists
	names, err := collection.Database().ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == collection.Name() {
			// Collection already exists, no need to create it
			return nil
		}
	}

	// Collection does not exist, create it
	err = collection.Database().CreateCollection(context.Background(), collection.Name())
	if err != nil {
		return err
	}

	return nil
}
