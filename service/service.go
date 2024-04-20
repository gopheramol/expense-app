// service.go

package service

import (
	"strconv"
	"time"

	"github.com/gopheramol/expense-app/db"

	"github.com/gopheramol/expense-app/model"
)

// ExpenseService provides methods for managing expenses
type ExpenseService interface {
	CreateExpense(name string, amountStr string) error
	GetExpenses() ([]model.Expense, error)
}

type expenseService struct {
	store db.ExpenseStore
}

// NewExpenseService creates a new instance of ExpenseService
func NewExpenseService(store db.ExpenseStore) ExpenseService {
	return &expenseService{store: store}
}

func (s *expenseService) CreateExpense(name string, amountStr string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return err
	}

	t := time.Now().Format("2006-01-02 3:04: PM")

	expense := &model.Expense{Name: name, Amount: amount, CreatedAt: t}
	return s.store.CreateExpense(expense)
}

func (s *expenseService) GetExpenses() ([]model.Expense, error) {
	return s.store.GetExpenses()
}
