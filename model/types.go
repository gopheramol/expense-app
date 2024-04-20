// types.go

package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Expense represents an expense entity
type Expense struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Amount    float64            `bson:"amount"`
	CreatedAt string             `bson:"createdAt"`
}
