package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gopheramol/expense-app/model"
	"github.com/gopheramol/expense-app/service"
)

// ExpenseHandler handles HTTP requests related to expenses
type ExpenseHandler interface {
	HandleCreateExpense(c *gin.Context)
	HandleGetExpenses(c *gin.Context)
}

type expenseHandler struct {
	service service.ExpenseService
}

// NewExpenseHandler creates a new instance of ExpenseHandler
func NewExpenseHandler(service service.ExpenseService) ExpenseHandler {
	return &expenseHandler{service: service}
}

func (h *expenseHandler) HandleCreateExpense(c *gin.Context) {
	name := c.PostForm("name")
	amountStr := c.PostForm("amount")
	if err := h.service.CreateExpense(name, amountStr); err != nil {
		c.String(http.StatusInternalServerError, "Error creating expense")
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *expenseHandler) HandleGetExpenses(c *gin.Context) {
	expenses, err := h.service.GetExpenses()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving expenses")
		return
	}
	total := calculateTotalAmount(expenses)
	c.HTML(http.StatusOK, "index.html", gin.H{"expenses": expenses, "total": total})
}

func calculateTotalAmount(expenses []model.Expense) float64 {
	total := 0.0
	for _, expense := range expenses {
		total += expense.Amount
	}
	return total
}
