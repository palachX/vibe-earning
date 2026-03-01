package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/example/vibe-earning/internal/models"
	"github.com/example/vibe-earning/internal/services"
)

type Handler struct {
	finance *services.FinanceService
}

func RegisterRoutes(rg *gin.RouterGroup, db *sql.DB) {
	h := &Handler{
		finance: services.NewFinanceService(db),
	}

	rg.POST("/incomes", h.createIncome)
	rg.GET("/incomes", h.listIncomes)
	rg.DELETE("/incomes/:id", h.deleteIncome)

	rg.POST("/expenses", h.createExpense)
	rg.GET("/expenses", h.listExpenses)
	rg.DELETE("/expenses/:id", h.deleteExpense)

	rg.POST("/recurring-expenses", h.createRecurringExpense)
	rg.GET("/recurring-expenses", h.listRecurring)
	rg.DELETE("/recurring-expenses/:id", h.deleteRecurring)

	rg.GET("/forecast", h.getForecast)
	rg.GET("/free-money", h.getFreeMoney)
	rg.GET("/current-balance", h.getCurrentBalance)
}

type incomeRequest struct {
	Date        string `json:"date" binding:"required"` // YYYY-MM-DD
	Amount      string `json:"amount" binding:"required"`
	Description string `json:"description"`
}

type expenseRequest struct {
	Date        string `json:"date" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	Description string `json:"description"`
}

type recurringRequest struct {
	Name      string  `json:"name" binding:"required"`
	Amount    string  `json:"amount" binding:"required"`
	Frequency string  `json:"frequency" binding:"required"` // "weekly" or "monthly"
	StartDate string  `json:"start_date" binding:"required"`
	EndDate   *string `json:"end_date"`
}

func parseDate(c *gin.Context, s string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected YYYY-MM-DD"})
		return time.Time{}, false
	}
	return t, true
}

func parseAmount(c *gin.Context, s string) (decimal.Decimal, bool) {
	dec, err := decimal.NewFromString(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return decimal.Zero, false
	}
	if dec.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return decimal.Zero, false
	}
	return dec, true
}

func (h *Handler) createIncome(c *gin.Context) {
	var req incomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, ok := parseDate(c, req.Date)
	if !ok {
		return
	}
	amount, ok := parseAmount(c, req.Amount)
	if !ok {
		return
	}

	inc, err := h.finance.CreateIncome(c.Request.Context(), date, amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, inc)
}

func (h *Handler) listIncomes(c *gin.Context) {
	items, err := h.finance.ListIncomes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []models.Income{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) deleteIncome(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.finance.DeleteIncome(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) createExpense(c *gin.Context) {
	var req expenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, ok := parseDate(c, req.Date)
	if !ok {
		return
	}
	amount, ok := parseAmount(c, req.Amount)
	if !ok {
		return
	}

	exp, err := h.finance.CreateExpense(c.Request.Context(), date, amount, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, exp)
}

func (h *Handler) listExpenses(c *gin.Context) {
	items, err := h.finance.ListExpenses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []models.Expense{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) deleteExpense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.finance.DeleteExpense(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) createRecurringExpense(c *gin.Context) {
	var req recurringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Frequency != "weekly" && req.Frequency != "monthly" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "frequency must be 'weekly' or 'monthly'"})
		return
	}

	startDate, ok := parseDate(c, req.StartDate)
	if !ok {
		return
	}

	var endDatePtr *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		t, ok := parseDate(c, *req.EndDate)
		if !ok {
			return
		}
		endDatePtr = &t
	}

	amount, ok := parseAmount(c, req.Amount)
	if !ok {
		return
	}

	rec, err := h.finance.CreateRecurringExpense(c.Request.Context(), req.Name, amount, req.Frequency, startDate, endDatePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rec)
}

func (h *Handler) listRecurring(c *gin.Context) {
	items, err := h.finance.ListRecurring(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if items == nil {
		items = []models.RecurringExpense{}
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) deleteRecurring(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.finance.DeleteRecurring(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) getForecast(c *gin.Context) {
	weeksStr := c.DefaultQuery("weeks", "52")
	weeks, err := strconv.Atoi(weeksStr)
	if err != nil || weeks <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weeks parameter"})
		return
	}

	fc, err := h.finance.ForecastBalance(c.Request.Context(), weeks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type respItem struct {
		WeekStart      string `json:"week_start"`
		OpeningBalance string `json:"opening_balance"`
		ClosingBalance string `json:"closing_balance"`
	}

	out := make([]respItem, 0, len(fc))
	for _, w := range fc {
		out = append(out, respItem{
			WeekStart:      w.WeekStart.Format("2006-01-02"),
			OpeningBalance: w.OpeningBalance.StringFixed(2),
			ClosingBalance: w.ClosingBalance.StringFixed(2),
		})
	}

	c.JSON(http.StatusOK, out)
}

func (h *Handler) getFreeMoney(c *gin.Context) {
	weeksStr := c.DefaultQuery("weeks", "52")
	weeks, err := strconv.Atoi(weeksStr)
	if err != nil || weeks <= 0 {
		weeks = 52
	}

	amount, err := h.finance.FreeMoney(c.Request.Context(), weeks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"free_money": amount.StringFixed(2),
		"weeks":      weeks,
	})
}

func (h *Handler) getCurrentBalance(c *gin.Context) {
	today := time.Now().UTC()
	amount, err := h.finance.GetCurrentBalance(c.Request.Context(), today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current_balance": amount.StringFixed(2),
		"date":            today.Format("2006-01-02"),
	})
}
