<template>
  <UPage>
    <UPageHeader
      title="Дашборд"
      description="Краткий обзор доходов, расходов, свободных денег и прогноза баланса."
    />
    <UPageBody>
      <div class="flex items-center justify-between mb-6">
        <UButton
          color="primary"
          :loading="loading"
          icon="i-lucide-refresh-ccw"
          @click="refresh"
        >
          Обновить
        </UButton>
      </div>
      
      <div class="space-y-6">
        <!-- Ключевые показатели: ПЕРВЫМИ, самое важное -->
        <UCard>
          <template #header>
            <p class="text-sm font-semibold">
              Ключевые показатели
            </p>
          </template>
          <div class="grid gap-6 md:grid-cols-3">
            <div class="p-4 rounded-lg bg-emerald-500/10 border border-emerald-500/20">
              <div class="flex items-center gap-2 mb-2">
                <UIcon name="i-lucide-wallet" class="text-emerald-400" />
                <p class="text-xs uppercase text-gray-400">
                  Сейчас на счетах
                </p>
              </div>
              <p class="text-2xl font-bold text-emerald-400">
                {{ formattedCurrentBalance }}
              </p>
              <p class="text-xs text-gray-500 mt-1">
                Баланс по всем счетам
              </p>
            </div>
            <div class="p-4 rounded-lg bg-sky-500/10 border border-sky-500/20">
              <div class="flex items-center gap-2 mb-2">
                <UIcon name="i-lucide-piggy-bank" class="text-sky-300" />
                <p class="text-xs uppercase text-gray-400">
                  Свободные деньги
                </p>
              </div>
              <p class="text-2xl font-bold text-sky-300">
                {{ freeMoneyFormatted }}
              </p>
              <p class="text-xs text-gray-500 mt-1">
                После всех платежей (52 нед)
              </p>
            </div>
            <div class="p-4 rounded-lg border" :class="netTotal >= 0 ? 'bg-emerald-500/10 border-emerald-500/20' : 'bg-rose-500/10 border-rose-500/20'">
              <div class="flex items-center gap-2 mb-2">
                <UIcon :name="netTotal >= 0 ? 'i-lucide-trending-up' : 'i-lucide-trending-down'" :class="netTotal >= 0 ? 'text-emerald-400' : 'text-rose-400'" />
                <p class="text-xs uppercase text-gray-400">
                  Чистый результат
                </p>
              </div>
              <p
                class="text-2xl font-bold"
                :class="netTotal >= 0 ? 'text-emerald-400' : 'text-rose-400'"
              >
                {{ netTotalFormatted }}
              </p>
              <p class="text-xs text-gray-500 mt-1">
                {{ totalIncomesFormatted }} − {{ totalExpensesFormatted }}
              </p>
            </div>
          </div>
        </UCard>

        <!-- Быстрые действия: компактные формы -->
        <div class="grid gap-6 lg:grid-cols-2">
          <UCard>
            <template #header>
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-plus-circle" class="text-emerald-400" />
                <p class="text-sm font-medium text-emerald-400">
                  Добавить доход
                </p>
              </div>
            </template>
            <IncomeForm />
          </UCard>

          <UCard>
            <template #header>
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-minus-circle" class="text-rose-400" />
                <p class="text-sm font-medium text-rose-400">
                  Добавить расход
                </p>
              </div>
            </template>
            <ExpenseForm />
          </UCard>
        </div>

        <!-- Прогнозы: сгруппированы вместе -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-calendar" class="text-gray-400" />
              <div>
                <p class="text-sm font-medium">
                  Прогноз: 8 недель
                </p>
                <p class="text-xs text-gray-400">
                  Доходы и траты по неделям
                </p>
              </div>
            </div>
          </template>
          <ForecastTable
            :items="forecast"
            :limit="8"
            :recurring="recurring"
            :expenses="expenses"
            :incomes="incomes"
          />
        </UCard>

        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <UIcon name="i-lucide-line-chart" class="text-gray-400" />
                <div>
                  <p class="text-sm font-medium">
                    Прогноз: 52 недели
                  </p>
                  <p class="text-xs text-gray-400">
                    Closing balance
                  </p>
                </div>
              </div>
              <div class="flex gap-3 text-xs text-gray-400">
                <div v-if="forecast.length" class="text-right">
                  <p class="uppercase text-[10px]">
                    Start
                  </p>
                  <p class="text-gray-200">
                    {{ firstWeekRange }}
                  </p>
                </div>
                <div v-if="forecast.length" class="text-right">
                  <p class="uppercase text-[10px]">
                    Диапазон
                  </p>
                  <p class="text-gray-200">
                    {{ closingMinMax }}
                  </p>
                </div>
              </div>
            </div>
          </template>
          <ForecastChart :items="forecast" />
        </UCard>

        <!-- Детализация: последние операции + графики -->
        <div class="grid gap-6 lg:grid-cols-2">
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <UIcon name="i-lucide-arrow-down-left" class="text-emerald-400" />
                  <p class="text-sm font-medium text-emerald-400">
                    Последние доходы
                  </p>
                </div>
                <UButton
                  size="xs"
                  variant="ghost"
                  to="/incomes"
                  icon="i-lucide-arrow-right"
                >
                  Все
                </UButton>
              </div>
            </template>
            <div
              v-if="!recentIncomes.length"
              class="text-sm text-gray-400 py-4"
            >
              Пока нет данных.
            </div>
            <div v-else>
              <UTable
                :data="recentIncomes"
                :columns="incomeColumns"
                :row-class="'text-sm'"
              />
            </div>
          </UCard>

          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <UIcon name="i-lucide-arrow-up-right" class="text-rose-400" />
                  <p class="text-sm font-medium text-rose-400">
                    Последние расходы
                  </p>
                </div>
                <UButton
                  size="xs"
                  variant="ghost"
                  to="/expenses"
                  icon="i-lucide-arrow-right"
                >
                  Все
                </UButton>
              </div>
            </template>
            <div
              v-if="!recentExpenses.length"
              class="text-sm text-gray-400 py-4"
            >
              Пока нет данных.
            </div>
            <div v-else>
              <UTable
                :data="recentExpenses"
                :columns="expenseColumns"
                :row-class="'text-sm'"
              />
            </div>
          </UCard>

          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <UIcon name="i-lucide-bar-chart-2" class="text-emerald-400" />
                  <p class="text-sm font-medium text-emerald-400">
                  Доходы по месяцам
                  </p>
                </div>
                <UButton
                  size="xs"
                  variant="ghost"
                  to="/incomes"
                  icon="i-lucide-arrow-right"
                >
                  Детали
                </UButton>
              </div>
            </template>
            <IncomeChart :items="incomes" />
          </UCard>

          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <UIcon name="i-lucide-bar-chart-2" class="text-rose-400" />
                  <p class="text-sm font-medium text-rose-400">
                    Расходы по месяцам
                  </p>
                </div>
                <UButton
                  size="xs"
                  variant="ghost"
                  to="/expenses"
                  icon="i-lucide-arrow-right"
                >
                  Детали
                </UButton>
              </div>
            </template>
            <ExpenseChart :items="expenses" />
          </UCard>
        </div>

        <!-- Анализ трат: внизу как дополнительная информация -->
        <UCard>
          <template #header>
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-pie-chart" class="text-gray-400" />
              <div>
                <p class="text-sm font-medium">
                  Анализ распределения денег
                </p>
                <p class="text-xs text-gray-400">
                  Как распределяются свободные деньги и расходы (52 недели)
                </p>
              </div>
            </div>
          </template>
          <SpendingBreakdown
            :current-balance="currentBalance"
            :free-money="freeMoney"
            :recurring="recurring"
            :forecast="forecast"
          />
        </UCard>
      </div>
    </UPageBody>
  </UPage>
</template>

<script setup lang="ts">
import { computed, onMounted, h } from 'vue'
import { useFinanceStore } from '~/stores/finance'
import ForecastChart from '~/components/ForecastChart.vue'
import ForecastTable from '~/components/ForecastTable.vue'
import IncomeForm from '~/components/IncomeForm.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import IncomeChart from '~/components/IncomeChart.vue'
import ExpenseChart from '~/components/ExpenseChart.vue'
import SpendingBreakdown from '~/components/SpendingBreakdown.vue'
import { UButton } from '#components'

const store = useFinanceStore()
const { forecast, loading, incomes, expenses, recurring, freeMoney, currentBalance } = storeToRefs(store)

// Последние транзакции (показываем последние 5)
const recentIncomes = computed(() => incomes.value.slice(0, 5))
const recentExpenses = computed(() => expenses.value.slice(0, 5))

// Колонки для таблицы доходов
const incomeColumns = [
  { accessorKey: 'date', header: 'Дата', cell: ({ row }: { row: any }) => {
    return formatDate(row.getValue('date'))
  } },
  { accessorKey: 'description', header: 'Описание' },
  { accessorKey: 'amount', header: 'Сумма', meta: { class: 'text-right' } as any, cell: ({ row }: { row: any }) => {
    return formatAmount(row.getValue('amount'))
  } },
  { accessorKey: 'actions', header: '', meta: { class: 'w-10 text-right' } as any, cell: ({ row }: { row: any }) => {
    return h(UButton, {
      icon: 'i-lucide-trash-2',
      color: 'error',
      variant: 'ghost',
      size: 'xs',
      onClick: () => removeIncome(row.original.id)
    })
  } }
]

// Колонки для таблицы расходов
const expenseColumns = [
  { accessorKey: 'date', header: 'Дата', cell: ({ row }: { row: any }) => {
    return formatDate(row.getValue('date'))
  } },
  { accessorKey: 'description', header: 'Описание' },
  { accessorKey: 'amount', header: 'Сумма', meta: { class: 'text-right' } as any, cell: ({ row }: { row: any }) => {
    return formatAmount(row.getValue('amount'))
  } },
  { accessorKey: 'actions', header: '', meta: { class: 'w-10 text-right' } as any, cell: ({ row }: { row: any }) => {
    return h(UButton, {
      icon: 'i-lucide-trash-2',
      color: 'error',
      variant: 'ghost',
      size: 'xs',
      onClick: () => removeExpense(row.original.id)
    })
  } }
]

const totalIncomes = computed(() => {
  return incomes.value.reduce((acc, i) => acc + Number(i.amount || 0), 0)
})

const totalExpenses = computed(() => {
  return expenses.value.reduce((acc, e) => acc + Number(e.amount || 0), 0)
})

const totalIncomesFormatted = computed(() => {
  return totalIncomes.value.toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
})

const totalExpensesFormatted = computed(() => {
  return totalExpenses.value.toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
})

const netTotal = computed(() => {
  return totalIncomes.value - totalExpenses.value
})

const netTotalFormatted = computed(() => {
  return netTotal.value.toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
})

const freeMoneyFormatted = computed(() => {
  if (!freeMoney.value) return '0,00'
  return Number(freeMoney.value).toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
})

const firstWeekRange = computed(() => {
  if (!forecast.value.length) return '—'
  const first = forecast.value[0]
  if (!first) return '—'
  return `${new Date(first.week_start).toLocaleDateString('ru-RU')}`
})

const closingMinMax = computed(() => {
  if (!forecast.value.length) return '—'
  const closes = forecast.value.map(f => Number(f.closing_balance))
  const min = Math.min(...closes)
  const max = Math.max(...closes)
  return `${min.toLocaleString('ru-RU', { minimumFractionDigits: 2, maximumFractionDigits: 2 })} — ${max.toLocaleString('ru-RU', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
})

const formattedCurrentBalance = computed(() => {
  if (!currentBalance.value) return '0,00'
  return Number(currentBalance.value).toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
})

// Функции форматирования
function formatDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU')
}

function formatAmount(a: string) {
  return Number(a).toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// Функции удаления
async function removeIncome(id: number) {
  await store.deleteIncome(id)
}

async function removeExpense(id: number) {
  await store.deleteExpense(id)
}

async function refresh() {
  await store.refreshDashboard()
}

onMounted(async () => {
  // Для дашборда сразу загружаем как детальные списки, так и прогноз/freeMoney
  await Promise.all([store.loadAll(), store.refreshDashboard()])
})
</script>
