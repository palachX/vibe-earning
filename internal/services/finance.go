package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/example/vibe-earning/internal/models"
)

const defaultUserID int64 = 1

type FinanceService struct {
	db *sql.DB
}

func NewFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{db: db}
}

// --- CRUD for incomes / expenses / recurring ---

// ListIncomes returns all incomes for the default user ordered by date desc.
func (s *FinanceService) ListIncomes(ctx context.Context) ([]models.Income, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, date, amount, description
		FROM incomes
		WHERE user_id = $1
		ORDER BY date DESC, id DESC
	`, defaultUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Income
	for rows.Next() {
		var inc models.Income
		var amt decimal.Decimal
		if err := rows.Scan(&inc.ID, &inc.UserID, &inc.Date, &amt, &inc.Description); err != nil {
			return nil, err
		}
		inc.Amount = amt
		res = append(res, inc)
	}
	return res, rows.Err()
}

func (s *FinanceService) DeleteIncome(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM incomes WHERE id = $1 AND user_id = $2
	`, id, defaultUserID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// ListExpenses returns all expenses for the default user ordered by date desc.
func (s *FinanceService) ListExpenses(ctx context.Context) ([]models.Expense, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, date, amount, description
		FROM expenses
		WHERE user_id = $1
		ORDER BY date DESC, id DESC
	`, defaultUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Expense
	for rows.Next() {
		var exp models.Expense
		var amt decimal.Decimal
		if err := rows.Scan(&exp.ID, &exp.UserID, &exp.Date, &amt, &exp.Description); err != nil {
			return nil, err
		}
		exp.Amount = amt
		res = append(res, exp)
	}
	return res, rows.Err()
}

func (s *FinanceService) DeleteExpense(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM expenses WHERE id = $1 AND user_id = $2
	`, id, defaultUserID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// ListRecurring returns all recurring expenses for the default user.
func (s *FinanceService) ListRecurring(ctx context.Context) ([]models.RecurringExpense, error) {
	recs, err := s.loadRecurring(ctx)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (s *FinanceService) DeleteRecurring(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx, `
		DELETE FROM recurring_expenses WHERE id = $1 AND user_id = $2
	`, id, defaultUserID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// CreateIncome inserts a new income row using a transaction.
func (s *FinanceService) CreateIncome(ctx context.Context, date time.Time, amount decimal.Decimal, description string) (*models.Income, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO incomes (user_id, date, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, date, amount, description
	`, defaultUserID, date, amount, description)

	var inc models.Income
	var amt decimal.Decimal
	if err = row.Scan(&inc.ID, &inc.UserID, &inc.Date, &amt, &inc.Description); err != nil {
		return nil, err
	}
	inc.Amount = amt

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &inc, nil
}

// CreateExpense inserts a new expense row using a transaction.
func (s *FinanceService) CreateExpense(ctx context.Context, date time.Time, amount decimal.Decimal, description string) (*models.Expense, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO expenses (user_id, date, amount, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, date, amount, description
	`, defaultUserID, date, amount, description)

	var exp models.Expense
	var amt decimal.Decimal
	if err = row.Scan(&exp.ID, &exp.UserID, &exp.Date, &amt, &exp.Description); err != nil {
		return nil, err
	}
	exp.Amount = amt

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &exp, nil
}

// CreateRecurringExpense inserts a new recurring expense.
func (s *FinanceService) CreateRecurringExpense(ctx context.Context, name string, amount decimal.Decimal, frequency string, startDate time.Time, endDate *time.Time) (*models.RecurringExpense, error) {
	if frequency == "" {
		return nil, errors.New("frequency is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, `
		INSERT INTO recurring_expenses (user_id, name, amount, frequency, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, amount, frequency, start_date, end_date
	`, defaultUserID, name, amount, frequency, startDate, endDate)

	var rec models.RecurringExpense
	var amt decimal.Decimal
	if err = row.Scan(&rec.ID, &rec.UserID, &rec.Name, &amt, &rec.Frequency, &rec.StartDate, &rec.EndDate); err != nil {
		return nil, err
	}
	rec.Amount = amt

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &rec, nil
}

// GetCurrentBalance returns total incomes - expenses - recurring up to today (inclusive).
func (s *FinanceService) GetCurrentBalance(ctx context.Context, today time.Time) (decimal.Decimal, error) {
	startOfDay := dateOnly(today)

	var incomeTotal decimal.Decimal
	if err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM incomes
		WHERE user_id = $1 AND date <= $2
	`, defaultUserID, startOfDay).Scan(&incomeTotal); err != nil {
		return decimal.Zero, err
	}

	var expenseTotal decimal.Decimal
	if err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = $1 AND date <= $2
	`, defaultUserID, startOfDay).Scan(&expenseTotal); err != nil {
		return decimal.Zero, err
	}

	// Sum realized recurring expenses up to today.
	recurringSum, err := s.sumRecurringUpTo(ctx, startOfDay)
	if err != nil {
		return decimal.Zero, err
	}

	return incomeTotal.Sub(expenseTotal).Sub(recurringSum), nil
}

// ForecastBalance returns weekly forecast starting from the next Monday for given number of weeks.
func (s *FinanceService) ForecastBalance(ctx context.Context, weeks int) ([]models.WeekForecast, error) {
	if weeks <= 0 {
		return nil, errors.New("weeks must be > 0")
	}

	today := dateOnly(time.Now().UTC())
	start := nextMonday(today)
	end := start.AddDate(0, 0, weeks*7)

	// current balance as of start-1 day
	balanceBeforeStart, err := s.GetCurrentBalance(ctx, start.AddDate(0, 0, -1))
	if err != nil {
		return nil, err
	}

	// preload incomes and expenses in range
	incomes, err := s.loadIncomesInRange(ctx, start, end)
	if err != nil {
		return nil, err
	}
	expenses, err := s.loadExpensesInRange(ctx, start, end)
	if err != nil {
		return nil, err
	}
	recurrings, err := s.loadRecurring(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]models.WeekForecast, 0, weeks)
	opening := balanceBeforeStart

	for i := 0; i < weeks; i++ {
		weekStart := start.AddDate(0, 0, i*7)
		weekEnd := weekStart.AddDate(0, 0, 7)

		incTotal := sumTransactions(incomes, weekStart, weekEnd)
		expTotal := sumTransactions(expenses, weekStart, weekEnd)
		recTotal := sumRecurringForWeek(recurrings, weekStart, weekEnd)

		closing := opening.Add(incTotal).Sub(expTotal).Sub(recTotal)

		results = append(results, models.WeekForecast{
			WeekStart:      weekStart,
			OpeningBalance: opening,
			ClosingBalance: closing,
			IncomeTotal:    incTotal,
			ExpenseTotal:   expTotal,
			RecurringTotal: recTotal,
		})

		opening = closing
	}

	return results, nil
}

// FreeMoney returns current balance minus future mandatory recurring obligations
// within the given number of weeks starting from today (default 52).
func (s *FinanceService) FreeMoney(ctx context.Context, weeks int) (decimal.Decimal, error) {
	if weeks <= 0 {
		weeks = 52
	}
	today := dateOnly(time.Now().UTC())
	balance, err := s.GetCurrentBalance(ctx, today)
	if err != nil {
		return decimal.Zero, err
	}

	horizonEnd := today.AddDate(0, 0, weeks*7)

	recurrings, err := s.loadRecurring(ctx)
	if err != nil {
		return decimal.Zero, err
	}

	var futureRecurring decimal.Decimal
	for _, r := range recurrings {
		// Пропускаем recurring, которые уже закончились
		if r.EndDate != nil && r.EndDate.Before(today) {
			continue
		}

		occ := r.StartDate
		if occ.Before(today) {
			occ = alignRecurringStart(r, today)
		}

		// Если после выравнивания дата всё ещё в прошлом или сегодня,
		// берём следующее вхождение
		if occ.Before(today) || occ.Equal(today) {
			occ = nextOccurrence(r, occ)
		}

		// Считаем все вхождения в горизонте (строго после today)
		for !occ.After(horizonEnd) {
			if r.EndDate == nil || !occ.After(*r.EndDate) {
				futureRecurring = futureRecurring.Add(r.Amount)
			}
			occ = nextOccurrence(r, occ)
		}
	}

	return balance.Sub(futureRecurring), nil
}

// --- helpers ---

func dateOnly(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func nextMonday(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysUntilMonday := (8 - weekday) % 7
	if daysUntilMonday == 0 {
		return t
	}
	return t.AddDate(0, 0, daysUntilMonday)
}

type datedTx struct {
	Date   time.Time
	Amount decimal.Decimal
}

func (s *FinanceService) loadIncomesInRange(ctx context.Context, from, to time.Time) ([]datedTx, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT date, amount
		FROM incomes
		WHERE user_id = $1 AND date >= $2 AND date < $3
	`, defaultUserID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []datedTx
	for rows.Next() {
		var d time.Time
		var amt decimal.Decimal
		if err := rows.Scan(&d, &amt); err != nil {
			return nil, err
		}
		res = append(res, datedTx{Date: d, Amount: amt})
	}
	return res, rows.Err()
}

func (s *FinanceService) loadExpensesInRange(ctx context.Context, from, to time.Time) ([]datedTx, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT date, amount
		FROM expenses
		WHERE user_id = $1 AND date >= $2 AND date < $3
	`, defaultUserID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []datedTx
	for rows.Next() {
		var d time.Time
		var amt decimal.Decimal
		if err := rows.Scan(&d, &amt); err != nil {
			return nil, err
		}
		res = append(res, datedTx{Date: d, Amount: amt})
	}
	return res, rows.Err()
}

func (s *FinanceService) loadRecurring(ctx context.Context) ([]models.RecurringExpense, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, name, amount, frequency, start_date, end_date
		FROM recurring_expenses
		WHERE user_id = $1
	`, defaultUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.RecurringExpense
	for rows.Next() {
		var r models.RecurringExpense
		var amt decimal.Decimal
		if err := rows.Scan(&r.ID, &r.UserID, &r.Name, &amt, &r.Frequency, &r.StartDate, &r.EndDate); err != nil {
			return nil, err
		}
		r.Amount = amt
		res = append(res, r)
	}
	return res, rows.Err()
}

func sumTransactions(txs []datedTx, from, to time.Time) decimal.Decimal {
	total := decimal.Zero
	for _, tx := range txs {
		if (tx.Date.Equal(from) || tx.Date.After(from)) && tx.Date.Before(to) {
			total = total.Add(tx.Amount)
		}
	}
	return total
}

func sumRecurringForWeek(recs []models.RecurringExpense, from, to time.Time) decimal.Decimal {
	total := decimal.Zero
	for _, r := range recs {
		// Пропускаем recurring, которые закончились до начала недели
		if r.EndDate != nil && r.EndDate.Before(from) {
			continue
		}

		occ := r.StartDate
		if occ.Before(from) {
			occ = alignRecurringStart(r, from)
		}

		for occ.Before(to) {
			// Проверяем, что вхождение в диапазоне [from, to) и не после end_date
			if (occ.Equal(from) || occ.After(from)) && (r.EndDate == nil || !occ.After(*r.EndDate)) {
				total = total.Add(r.Amount)
			}
			occ = nextOccurrence(r, occ)
		}
	}
	return total
}

func sumRecurringUpToDate(r models.RecurringExpense, upTo time.Time) decimal.Decimal {
	total := decimal.Zero
	occ := r.StartDate
	for !occ.After(upTo) {
		if r.EndDate != nil && occ.After(*r.EndDate) {
			break
		}
		total = total.Add(r.Amount)
		occ = nextOccurrence(r, occ)
	}
	return total
}

func (s *FinanceService) sumRecurringUpTo(ctx context.Context, upTo time.Time) (decimal.Decimal, error) {
	recs, err := s.loadRecurring(ctx)
	if err != nil {
		return decimal.Zero, err
	}
	total := decimal.Zero
	for _, r := range recs {
		total = total.Add(sumRecurringUpToDate(r, upTo))
	}
	return total, nil
}

func nextOccurrence(r models.RecurringExpense, from time.Time) time.Time {
	switch r.Frequency {
	case "weekly":
		return from.AddDate(0, 0, 7)
	case "monthly":
		return from.AddDate(0, 1, 0)
	default:
		// treat unknown as monthly to avoid infinite loops
		return from.AddDate(0, 1, 0)
	}
}

func alignRecurringStart(r models.RecurringExpense, from time.Time) time.Time {
	occ := r.StartDate
	if !occ.Before(from) {
		return occ
	}
	switch r.Frequency {
	case "weekly":
		// jump in 7-day steps
		for occ.Before(from) {
			occ = occ.AddDate(0, 0, 7)
		}
	case "monthly":
		for occ.Before(from) {
			occ = occ.AddDate(0, 1, 0)
		}
	default:
		for occ.Before(from) {
			occ = occ.AddDate(0, 1, 0)
		}
	}
	return occ
}

