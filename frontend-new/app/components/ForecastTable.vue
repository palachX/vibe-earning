<template>
  <div class="space-y-3">
    <!-- Итоги за горизонт -->
    <div class="rounded-xl border border-slate-800 bg-slate-900/30 p-4">
      <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
        <div>
          <p class="text-sm font-medium text-slate-200">
            Итоги за {{ limitedItems.length }} недель
          </p>
          <p class="text-xs text-slate-500">
            Net change = Доходы − Разовые расходы − Постоянные траты
          </p>
        </div>
        <UBadge variant="subtle" color="neutral">
          Closing: {{ formatAmount(closingAfterHorizon) }}
        </UBadge>
      </div>

      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <div class="rounded-lg border border-emerald-900/30 bg-emerald-950/10 p-3">
          <p class="text-xs uppercase tracking-wide text-emerald-500/80">
            Доходы
          </p>
          <p class="mt-1 text-lg font-semibold text-emerald-400">
            + {{ formatAmount(sumIncome) }}
          </p>
        </div>

        <div class="rounded-lg border border-rose-900/30 bg-rose-950/10 p-3">
          <p class="text-xs uppercase tracking-wide text-rose-500/80">
            Разовые расходы
          </p>
          <p class="mt-1 text-lg font-semibold text-rose-400">
            − {{ formatAmount(sumOneTimeExpenses) }}
          </p>
        </div>

        <div class="rounded-lg border border-amber-900/30 bg-amber-950/10 p-3">
          <p class="text-xs uppercase tracking-wide text-amber-500/80">
            Постоянные траты
          </p>
          <p class="mt-1 text-lg font-semibold text-amber-400">
            − {{ formatAmount(sumRecurring) }}
          </p>
        </div>

        <div class="rounded-lg border border-slate-700 bg-slate-950/40 p-3">
          <p class="text-xs uppercase tracking-wide text-slate-500">
            Изменение баланса
          </p>
          <p
            class="mt-1 text-lg font-semibold"
            :class="netOverall >= 0 ? 'text-emerald-400' : 'text-rose-400'"
          >
            {{ netOverall >= 0 ? '+' : '−' }} {{ formatAmount(Math.abs(netOverall)) }}
          </p>
          <p class="mt-1 text-xs text-slate-500">
            Opening: {{ formatAmount(openingBeforeHorizon) }}
          </p>
        </div>
      </div>

      <div class="mt-4 grid gap-4 lg:grid-cols-2">
        <div class="rounded-lg border border-slate-800 bg-slate-900/30 p-3">
          <div class="mb-2 flex items-center justify-between gap-2">
            <p class="text-xs font-semibold uppercase tracking-wide text-amber-400">
              ТОП recurring
            </p>
            <UBadge variant="subtle" color="neutral">
              {{ recurringCountTotal }} платеж(ей)
            </UBadge>
          </div>
          <ul v-if="topRecurring.length" class="space-y-1.5">
            <li
              v-for="t in topRecurring"
              :key="t.name"
              class="flex items-center justify-between gap-3 rounded-md border border-slate-800/60 bg-slate-950/25 px-3 py-2"
            >
              <div class="min-w-0">
                <p class="truncate text-sm text-slate-100">
                  {{ t.name }}
                </p>
                <p class="text-xs text-slate-500">
                  {{ t.count }} платеж(ей) · {{ t.sampleDates.join(', ') || '—' }}
                </p>
              </div>
              <p class="shrink-0 text-sm font-semibold text-amber-300">
                − {{ formatAmount(t.total) }}
              </p>
            </li>
          </ul>
          <p v-else class="text-sm text-slate-500">
            Нет recurring-платежей в горизонте.
          </p>
        </div>

        <div class="rounded-lg border border-slate-800 bg-slate-900/30 p-3">
          <div class="mb-2 flex items-center justify-between gap-2">
            <p class="text-xs font-semibold uppercase tracking-wide text-rose-400">
              ТОП разовых расходов
            </p>
            <UBadge variant="subtle" color="neutral">
              {{ oneTimeExpenseCountTotal }} операций
            </UBadge>
          </div>
          <ul v-if="topOneTimeExpenses.length" class="space-y-1.5">
            <li
              v-for="t in topOneTimeExpenses"
              :key="t.name"
              class="flex items-center justify-between gap-3 rounded-md border border-slate-800/60 bg-slate-950/25 px-3 py-2"
            >
              <div class="min-w-0">
                <p class="truncate text-sm text-slate-100">
                  {{ t.name }}
                </p>
                <p class="text-xs text-slate-500">
                  {{ t.count }} операций · {{ t.sampleDates.join(', ') || '—' }}
                </p>
              </div>
              <p class="shrink-0 text-sm font-semibold text-rose-300">
                − {{ formatAmount(t.total) }}
              </p>
            </li>
          </ul>
          <p v-else class="text-sm text-slate-500">
            Нет разовых расходов в горизонте.
          </p>
        </div>
      </div>
    </div>

    <div class="overflow-x-auto rounded-xl border border-slate-800 bg-slate-900/40">
      <table class="min-w-full divide-y divide-slate-800 text-sm">
        <thead class="bg-slate-900/60">
          <tr>
            <th class="w-8 px-2 py-2" />
            <th class="px-3 py-2 text-left text-xs font-semibold uppercase tracking-wide text-slate-300">
              Неделя
            </th>
            <th class="px-3 py-2 text-right text-xs font-semibold uppercase tracking-wide text-emerald-500/90">
              Доходы
            </th>
            <th class="px-3 py-2 text-right text-xs font-semibold uppercase tracking-wide text-rose-500/90">
              Разовые
            </th>
            <th class="px-3 py-2 text-right text-xs font-semibold uppercase tracking-wide text-amber-500/90">
              Постоянные
            </th>
            <th class="px-3 py-2 text-right text-xs font-semibold uppercase tracking-wide text-slate-300">
              Изменение
            </th>
            <th class="px-3 py-2 text-right text-xs font-semibold uppercase tracking-wide text-slate-400">
              На начало
            </th>
            <th class="px-3 py-2 text-right text-xs font-semibold uppercase tracking-wide text-sky-400">
              На конец
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <template v-for="item in limitedItems" :key="item.week_start">
            <tr
              class="cursor-pointer transition-colors hover:bg-slate-800/40"
              :class="{ 'bg-slate-800/25': isExpanded(item.week_start) }"
              @click="toggle(item.week_start)"
            >
              <td class="px-2 py-2 text-center text-slate-500">
                <UIcon
                  :name="isExpanded(item.week_start) ? 'i-lucide-chevron-down' : 'i-lucide-chevron-right'"
                  class="h-4 w-4"
                />
              </td>
              <td class="px-3 py-2 text-slate-100">
                <p>{{ formatDate(item.week_start) }}</p>
                <p class="text-xs text-slate-500">
                  {{ weekRangeLabel(item.week_start) }}
                </p>
              </td>
              <td class="px-3 py-2 text-right text-emerald-400">
                <span v-if="hasIncome(item)">+ {{ formatAmount(item.income_total) }}</span>
                <span v-else class="text-slate-600">—</span>
              </td>
              <td class="px-3 py-2 text-right text-rose-400">
                <span v-if="hasExpense(item)">− {{ formatAmount(item.expense_total) }}</span>
                <span v-else class="text-slate-600">—</span>
              </td>
              <td class="px-3 py-2 text-right text-amber-400">
                <span v-if="hasRecurring(item)">− {{ formatAmount(item.recurring_total) }}</span>
                <span v-else class="text-slate-600">—</span>
              </td>
              <td class="px-3 py-2 text-right font-medium" :class="netClass(item)">
                {{ formatNet(item) }}
              </td>
              <td class="px-3 py-2 text-right text-slate-300">
                {{ formatAmount(item.opening_balance) }}
              </td>
              <td class="px-3 py-2 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <UIcon
                    v-if="closingDelta(item) !== 0"
                    :name="closingDelta(item) > 0 ? 'i-lucide-trending-up' : 'i-lucide-trending-down'"
                    class="h-3.5 w-3.5"
                    :class="closingDelta(item) > 0 ? 'text-emerald-500' : 'text-rose-500'"
                  />
                  <span class="font-semibold text-sky-300">
                    {{ formatAmount(item.closing_balance) }}
                  </span>
                </div>
              </td>
            </tr>

            <tr v-if="isExpanded(item.week_start)" class="bg-slate-950/50">
              <td colspan="8" class="px-4 py-4">
                <div class="space-y-4">
                  <div class="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-400">
                    <span class="text-slate-300">{{ formatAmount(item.opening_balance) }}</span>
                    <span class="mx-1">+</span>
                    <span class="text-emerald-400">{{ formatAmount(item.income_total) }}</span>
                    <span class="mx-1">−</span>
                    <span class="text-rose-400">{{ formatAmount(item.expense_total) }}</span>
                    <span class="mx-1">−</span>
                    <span class="text-amber-400">{{ formatAmount(item.recurring_total) }}</span>
                    <span class="mx-1">=</span>
                    <span class="font-medium text-sky-300">{{ formatAmount(item.closing_balance) }}</span>
                  </div>

                  <div class="grid gap-4 md:grid-cols-3">
                    <div class="rounded-lg border border-amber-900/30 bg-amber-950/10 p-3">
                      <p class="mb-2 text-xs font-semibold uppercase tracking-wide text-amber-500/90">
                        Постоянные траты
                      </p>
                      <ul v-if="recurringInWeek(item.week_start).length" class="space-y-1.5">
                        <li
                          v-for="occ in recurringInWeek(item.week_start)"
                          :key="`${occ.id}-${occ.date}`"
                          class="flex items-center justify-between gap-2 text-xs"
                        >
                          <span class="truncate text-slate-200">{{ occ.name }}</span>
                          <span class="shrink-0 text-amber-300">
                            {{ formatShortDate(occ.date) }} · {{ formatAmount(occ.amount) }}
                          </span>
                        </li>
                      </ul>
                      <p v-else class="text-xs text-slate-500">
                        Нет платежей на этой неделе
                      </p>
                    </div>

                    <div class="rounded-lg border border-rose-900/30 bg-rose-950/10 p-3">
                      <p class="mb-2 text-xs font-semibold uppercase tracking-wide text-rose-500/90">
                        Разовые расходы
                      </p>
                      <ul v-if="expensesInWeek(item.week_start).length" class="space-y-1.5">
                        <li
                          v-for="tx in expensesInWeek(item.week_start)"
                          :key="tx.id"
                          class="flex items-center justify-between gap-2 text-xs"
                        >
                          <span class="truncate text-slate-200">{{ tx.description }}</span>
                          <span class="shrink-0 text-rose-300">
                            {{ formatShortDate(tx.date) }} · {{ formatAmount(tx.amount) }}
                          </span>
                        </li>
                      </ul>
                      <p v-else class="text-xs text-slate-500">
                        Нет разовых расходов
                      </p>
                    </div>

                    <div class="rounded-lg border border-emerald-900/30 bg-emerald-950/10 p-3">
                      <p class="mb-2 text-xs font-semibold uppercase tracking-wide text-emerald-500/90">
                        Доходы
                      </p>
                      <ul v-if="incomesInWeek(item.week_start).length" class="space-y-1.5">
                        <li
                          v-for="tx in incomesInWeek(item.week_start)"
                          :key="tx.id"
                          class="flex items-center justify-between gap-2 text-xs"
                        >
                          <span class="truncate text-slate-200">{{ tx.description }}</span>
                          <span class="shrink-0 text-emerald-300">
                            {{ formatShortDate(tx.date) }} · {{ formatAmount(tx.amount) }}
                          </span>
                        </li>
                      </ul>
                      <p v-else class="text-xs text-slate-500">
                        Нет доходов
                      </p>
                    </div>
                  </div>

                  <p v-if="closingDelta(item) < 0" class="text-xs text-slate-500">
                    Баланс уменьшается на {{ formatAmount(Math.abs(closingDelta(item))) }},
                    потому что расходы ({{ formatAmount(totalOutflow(item)) }})
                    превышают доходы ({{ formatAmount(item.income_total) }}).
                  </p>
                  <p v-else-if="closingDelta(item) > 0" class="text-xs text-slate-500">
                    Баланс растёт на {{ formatAmount(closingDelta(item)) }} — доходы больше расходов.
                  </p>
                  <p v-else class="text-xs text-slate-500">
                    Баланс не меняется: доходы и расходы за неделю сбалансированы.
                  </p>
                </div>
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <p class="text-xs text-slate-500">
      Нажмите на строку, чтобы увидеть детали. Closing = Opening + Доходы − Разовые − Постоянные.
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Expense, ForecastItem, Income, RecurringExpense } from '~/stores/finance'
import {
  getExpensesInWeek,
  getIncomesInWeek,
  getRecurringOccurrencesInWeek,
  weekNetChange
} from '~/composables/useSpendingBreakdown'

interface Props {
  items: ForecastItem[]
  limit?: number
  recurring?: RecurringExpense[]
  expenses?: Expense[]
  incomes?: Income[]
}

const props = withDefaults(defineProps<Props>(), {
  recurring: () => [],
  expenses: () => [],
  incomes: () => []
})

const expanded = ref<Set<string>>(new Set())

const limitedItems = computed(() =>
  props.limit ? props.items.slice(0, props.limit) : props.items
)

const sumIncome = computed(() =>
  limitedItems.value.reduce((acc, i) => acc + Number(i.income_total || 0), 0)
)
const sumOneTimeExpenses = computed(() =>
  limitedItems.value.reduce((acc, i) => acc + Number(i.expense_total || 0), 0)
)
const sumRecurring = computed(() =>
  limitedItems.value.reduce((acc, i) => acc + Number(i.recurring_total || 0), 0)
)

const openingBeforeHorizon = computed(
  () => Number(limitedItems.value[0]?.opening_balance || 0)
)
const closingAfterHorizon = computed(
  () => Number(limitedItems.value[limitedItems.value.length - 1]?.closing_balance || 0)
)
const netOverall = computed(() => closingAfterHorizon.value - openingBeforeHorizon.value)

interface TopItem {
  name: string
  total: number
  count: number
  sampleDates: string[]
}

const recurringCountTotal = computed(() => {
  let count = 0
  for (const w of limitedItems.value) {
    count += getRecurringOccurrencesInWeek(props.recurring, w.week_start).length
  }
  return count
})

const oneTimeExpenseCountTotal = computed(() => {
  let count = 0
  for (const w of limitedItems.value) {
    count += getExpensesInWeek(props.expenses, w.week_start).length
  }
  return count
})

const topRecurring = computed<TopItem[]>(() => {
  const map = new Map<string, { total: number, count: number, sampleDates: Set<string> }>()

  for (const w of limitedItems.value) {
    const occs = getRecurringOccurrencesInWeek(props.recurring, w.week_start)
    for (const occ of occs) {
      const entry = map.get(occ.name) || { total: 0, count: 0, sampleDates: new Set<string>() }
      entry.total += occ.amount
      entry.count += 1
      if (entry.sampleDates.size < 3) entry.sampleDates.add(occ.date)
      map.set(occ.name, entry)
    }
  }

  return [...map.entries()]
    .map(([name, v]) => ({
      name,
      total: v.total,
      count: v.count,
      sampleDates: [...v.sampleDates].slice(0, 3).map(d => formatShortDate(d))
    }))
    .sort((a, b) => b.total - a.total)
    .slice(0, 3)
})

const topOneTimeExpenses = computed<TopItem[]>(() => {
  const map = new Map<string, { total: number, count: number, sampleDates: Set<string> }>()

  for (const w of limitedItems.value) {
    const txs = getExpensesInWeek(props.expenses, w.week_start)
    for (const tx of txs) {
      const entry = map.get(tx.description) || { total: 0, count: 0, sampleDates: new Set<string>() }
      entry.total += tx.amount
      entry.count += 1
      if (entry.sampleDates.size < 3) entry.sampleDates.add(tx.date)
      map.set(tx.description, entry)
    }
  }

  return [...map.entries()]
    .map(([name, v]) => ({
      name,
      total: v.total,
      count: v.count,
      sampleDates: [...v.sampleDates].slice(0, 3).map(d => formatShortDate(d))
    }))
    .sort((a, b) => b.total - a.total)
    .slice(0, 3)
})

function toggle(weekStart: string) {
  const next = new Set(expanded.value)
  if (next.has(weekStart)) {
    next.delete(weekStart)
  } else {
    next.add(weekStart)
  }
  expanded.value = next
}

function isExpanded(weekStart: string) {
  return expanded.value.has(weekStart)
}

function recurringInWeek(weekStart: string) {
  return getRecurringOccurrencesInWeek(props.recurring, weekStart)
}

function expensesInWeek(weekStart: string) {
  return getExpensesInWeek(props.expenses, weekStart)
}

function incomesInWeek(weekStart: string) {
  return getIncomesInWeek(props.incomes, weekStart)
}

function hasIncome(item: ForecastItem) {
  return Number(item.income_total || 0) > 0
}

function hasExpense(item: ForecastItem) {
  return Number(item.expense_total || 0) > 0
}

function hasRecurring(item: ForecastItem) {
  return Number(item.recurring_total || 0) > 0
}

function closingDelta(item: ForecastItem) {
  return Number(item.closing_balance) - Number(item.opening_balance)
}

function totalOutflow(item: ForecastItem) {
  return Number(item.expense_total || 0) + Number(item.recurring_total || 0)
}

function netClass(item: ForecastItem) {
  const net = weekNetChange(item)
  if (net > 0) return 'text-emerald-400'
  if (net < 0) return 'text-rose-400'
  return 'text-slate-500'
}

function formatNet(item: ForecastItem) {
  const net = weekNetChange(item)
  if (net === 0) return '0,00'
  const sign = net > 0 ? '+' : '−'
  return `${sign} ${formatAmount(Math.abs(net))}`
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short'
  })
}

function formatShortDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short'
  })
}

function weekRangeLabel(weekStart: string) {
  const start = new Date(weekStart)
  const end = new Date(start)
  end.setDate(end.getDate() + 6)
  const fmt = (d: Date) => d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
  return `${fmt(start)} — ${fmt(end)}`
}

function formatAmount(a: string | number | null | undefined) {
  return Number(a || 0).toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}
</script>
