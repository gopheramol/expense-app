// store.go

package db

import (
	"context"
	"log"
	"time"

	"github.com/gopheramol/expense-app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ExpenseStore provides methods for interacting with the expenses collection
type ExpenseStore interface {
	CreateExpense(expense *model.Expense) error
	GetExpenses() ([]model.Expense, error)
}

type mongoStore struct {
	collection *mongo.Collection
}

// NewMongoStore creates a new instance of ExpenseStore backed by MongoDB
func NewMongoStore(client *mongo.Client, dbName, collectionName string) (ExpenseStore, error) {
	// Create or get the expenses collection
	collection := client.Database(dbName).Collection(collectionName)
	err := createCollectionIfNotExists(collection)
	if err != nil {
		log.Fatal(err)
	}
	return &mongoStore{collection: collection}, nil
}

func (m *mongoStore) CreateExpense(expense *model.Expense) error {
	_, err := m.collection.InsertOne(context.Background(), expense)
	return err
}

func (m *mongoStore) GetExpenses() ([]model.Expense, error) {
	// Define a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the options for sorting
	findOptions := options.Find().SetSort(bson.D{{"createdAt", -1}}) // Sort by createdAt field in descending order

	// Find all documents in the collection, sorted by createdAt field in descending order
	cursor, err := m.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Initialize a slice to hold the expenses
	var expenses []model.Expense

	// Iterate over the cursor and decode each document
	for cursor.Next(ctx) {
		var expense model.Expense
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
