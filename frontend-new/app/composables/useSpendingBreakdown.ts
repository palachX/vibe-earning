import type { Expense, Income, RecurringExpense } from '~/stores/finance'

export interface WeekOccurrence {
  id: number
  name: string
  amount: number
  date: string
}

export interface WeekTransaction {
  id: number
  description: string
  amount: number
  date: string
}

export interface RecurringBreakdownItem {
  id: number
  name: string
  frequency: string
  amount: number
  total: number
  percent: number
}

function dateOnly(d: Date): Date {
  return new Date(d.getFullYear(), d.getMonth(), d.getDate())
}

function toDateString(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function nextOccurrence(frequency: string, from: Date): Date {
  const d = new Date(from)
  if (frequency === 'weekly') {
    d.setDate(d.getDate() + 7)
  } else {
    d.setMonth(d.getMonth() + 1)
  }
  return d
}

function alignRecurringStart(r: RecurringExpense, from: Date): Date {
  let occ = dateOnly(new Date(r.start_date))
  const fromDate = dateOnly(from)
  if (occ >= fromDate) return occ

  if (r.frequency === 'weekly') {
    while (occ < fromDate) {
      occ = nextOccurrence('weekly', occ)
    }
  } else {
    while (occ < fromDate) {
      occ = nextOccurrence('monthly', occ)
    }
  }
  return occ
}

export function projectRecurringFuture(r: RecurringExpense, weeks: number): number {
  const today = dateOnly(new Date())
  const horizonEnd = new Date(today)
  horizonEnd.setDate(horizonEnd.getDate() + weeks * 7)

  if (r.end_date && dateOnly(new Date(r.end_date)) < today) {
    return 0
  }

  let occ = dateOnly(new Date(r.start_date))
  if (occ < today) {
    occ = alignRecurringStart(r, today)
  }
  if (occ <= today) {
    occ = nextOccurrence(r.frequency, occ)
  }

  const endDate = r.end_date ? dateOnly(new Date(r.end_date)) : null
  let total = 0

  while (occ <= horizonEnd) {
    if (!endDate || occ <= endDate) {
      total += Number(r.amount)
    }
    occ = nextOccurrence(r.frequency, occ)
  }

  return total
}

export function buildRecurringBreakdown(
  items: RecurringExpense[],
  weeks: number
): RecurringBreakdownItem[] {
  const raw = items
    .map(r => ({
      id: r.id,
      name: r.name,
      frequency: r.frequency,
      amount: Number(r.amount),
      total: projectRecurringFuture(r, weeks)
    }))
    .filter(r => r.total > 0)
    .sort((a, b) => b.total - a.total)

  const sum = raw.reduce((acc, r) => acc + r.total, 0)

  return raw.map(r => ({
    ...r,
    percent: sum > 0 ? (r.total / sum) * 100 : 0
  }))
}

export function formatMoney(value: number | string | null | undefined): string {
  const n = Number(value || 0)
  return n.toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

export function humanFrequency(f: string): string {
  if (f === 'weekly') return 'Еженедельно'
  if (f === 'monthly') return 'Ежемесячно'
  return f
}

function weekRange(weekStartStr: string): { from: Date, to: Date } {
  const from = dateOnly(new Date(weekStartStr))
  const to = new Date(from)
  to.setDate(to.getDate() + 7)
  return { from, to }
}

function isInWeek(dateStr: string, from: Date, to: Date): boolean {
  const d = dateOnly(new Date(dateStr))
  return d >= from && d < to
}

export function getRecurringOccurrencesInWeek(
  items: RecurringExpense[],
  weekStartStr: string
): WeekOccurrence[] {
  const { from, to } = weekRange(weekStartStr)
  const results: WeekOccurrence[] = []

  for (const r of items) {
    if (r.end_date && dateOnly(new Date(r.end_date)) < from) {
      continue
    }

    let occ = dateOnly(new Date(r.start_date))
    if (occ < from) {
      occ = alignRecurringStart(r, from)
    }

    const endDate = r.end_date ? dateOnly(new Date(r.end_date)) : null

    while (occ < to) {
      if (occ >= from && (!endDate || occ <= endDate)) {
        results.push({
          id: r.id,
          name: r.name,
          amount: Number(r.amount),
          date: toDateString(occ)
        })
      }
      occ = nextOccurrence(r.frequency, occ)
    }
  }

  return results.sort((a, b) => a.date.localeCompare(b.date) || a.name.localeCompare(b.name))
}

export function getExpensesInWeek(
  items: Expense[],
  weekStartStr: string
): WeekTransaction[] {
  const { from, to } = weekRange(weekStartStr)
  return items
    .filter(e => isInWeek(e.date, from, to))
    .map(e => ({
      id: e.id,
      description: e.description || 'Без описания',
      amount: Number(e.amount),
      date: e.date
    }))
    .sort((a, b) => a.date.localeCompare(b.date))
}

export function getIncomesInWeek(
  items: Income[],
  weekStartStr: string
): WeekTransaction[] {
  const { from, to } = weekRange(weekStartStr)
  return items
    .filter(i => isInWeek(i.date, from, to))
    .map(i => ({
      id: i.id,
      description: i.description || 'Без описания',
      amount: Number(i.amount),
      date: i.date
    }))
    .sort((a, b) => a.date.localeCompare(b.date))
}

export function weekNetChange(item: {
  income_total?: string
  expense_total?: string
  recurring_total?: string
}): number {
  return Number(item.income_total || 0)
    - Number(item.expense_total || 0)
    - Number(item.recurring_total || 0)
}
