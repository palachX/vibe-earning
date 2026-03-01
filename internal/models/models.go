package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Income struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Date        time.Time       `json:"date"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description,omitempty"`
}

type Expense struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	Date        time.Time       `json:"date"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description,omitempty"`
}

type RecurringExpense struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Name      string          `json:"name"`
	Amount    decimal.Decimal `json:"amount"`
	Frequency string          `json:"frequency"`
	StartDate time.Time       `json:"start_date"`
	EndDate   *time.Time      `json:"end_date,omitempty"`
}

type WeekForecast struct {
	WeekStart       time.Time       `json:"week_start"`
	OpeningBalance  decimal.Decimal `json:"opening_balance"`
	ClosingBalance  decimal.Decimal `json:"closing_balance"`
	IncomeTotal     decimal.Decimal `json:"income_total"`
	ExpenseTotal    decimal.Decimal `json:"expense_total"`
	RecurringTotal  decimal.Decimal `json:"recurring_total"`
}

