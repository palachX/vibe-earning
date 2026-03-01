<template>
  <UPage>
    <UPageHeader
      title="Дашборд"
      description="Краткий обзор доходов, расходов, свободных денег и прогноза баланса."
    />
    <UPageBody>
      <UButton
        color="primary"
        :loading="loading"
        icon="i-lucide-refresh-ccw"
        @click="refresh"
      >
        Обновить
      </UButton>
      <!-- Карточки по доходам / расходам / свободным деньгам -->
      <div class="mb-6">
        <UCard>
          <template #header>
            <p class="text-sm font-medium">
              Сводка по движениям средств
            </p>
          </template>
          <div class="grid gap-4 md:grid-cols-4">
            <div>
              <p class="text-xs uppercase text-gray-400">
                Доходы (все время)
              </p>
              <p class="mt-1 text-xl font-semibold text-emerald-400">
                {{ totalIncomesFormatted }}
              </p>
              <p class="text-xs text-gray-500">
                {{ incomes.length }} операций
              </p>
            </div>
            <div>
              <p class="text-xs uppercase text-gray-400">
                Текущий баланс
              </p>
              <p class="mt-1 text-xl font-semibold text-emerald-400">
                {{ formattedCurrentBalance }}
              </p>
            </div>
            <div>
              <p class="text-xs uppercase text-gray-400">
                Расходы (все время)
              </p>
              <p class="mt-1 text-xl font-semibold text-rose-400">
                {{ totalExpensesFormatted }}
              </p>
              <p class="text-xs text-gray-500">
                {{ expenses.length }} операций
              </p>
            </div>
            <div>
              <p class="text-xs uppercase text-gray-400">
                Свободные деньги
              </p>
              <p class="mt-1 text-xl font-semibold text-sky-300">
                {{ freeMoneyFormatted }}
              </p>
              <p class="text-xs text-gray-500">
                с учётом будущих recurring (52 недели)
              </p>
            </div>
          </div>
        </UCard>
      </div>

      <!-- Формы добавления доходов и расходов -->
      <div class="grid gap-6 lg:grid-cols-2 mb-6">
        <UCard>
          <template #header>
            <p class="text-sm font-medium text-emerald-400">
              Добавить доход
            </p>
          </template>
          <IncomeForm />
        </UCard>

        <UCard>
          <template #header>
            <p class="text-sm font-medium text-rose-400">
              Добавить расход
            </p>
          </template>
          <ExpenseForm />
        </UCard>
      </div>
      <div class="grid gap-6 lg:grid-cols-2">
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium">
                  Прогноз на 52 недели
                </p>
                <p class="text-xs text-gray-400">
                  Ожидаемый closing balance по неделям
                </p>
              </div>
              <div class="flex gap-4 text-xs text-gray-400">
                <div v-if="forecast.length">
                  <p class="uppercase">
                    Текущая неделя
                  </p>
                  <p class="mt-0.5 text-gray-200">
                    {{ firstWeekRange }}
                  </p>
                </div>
                <div v-if="forecast.length">
                  <p class="uppercase">
                    Диапазон closing
                  </p>
                  <p class="mt-0.5 text-gray-200">
                    {{ closingMinMax }}
                  </p>
                </div>
              </div>
            </div>
          </template>
          <ForecastChart :items="forecast" />
        </UCard>

        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <p class="text-sm font-medium">
                Ближайшие 8 недель
              </p>
            </div>
          </template>
          <ForecastTable
            :items="forecast"
            :limit="8"
          />
        </UCard>
      </div>
      <!-- Последние транзакции -->
      <div class="grid gap-6 lg:grid-cols-2">
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <p class="text-sm font-medium text-emerald-400">
                Последние доходы
              </p>
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
            class="text-sm text-gray-400"
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
              <p class="text-sm font-medium text-rose-400">
                Последние расходы
              </p>
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
            class="text-sm text-gray-400"
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
      </div>

      <!-- Графики доходов и расходов -->
      <div class="grid gap-6 lg:grid-cols-2">
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <p class="text-sm font-medium text-emerald-400">
                График доходов по месяцам
              </p>
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
              <p class="text-sm font-medium text-rose-400">
                График расходов по месяцам
              </p>
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
import { UButton } from '#components'

const store = useFinanceStore()
const { forecast, loading, incomes, expenses, freeMoney, currentBalance } = storeToRefs(store)

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

const totalIncomesFormatted = computed(() => {
  const sum = incomes.value.reduce((acc, i) => acc + Number(i.amount || 0), 0)
  return sum.toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
})

const totalExpensesFormatted = computed(() => {
  const sum = expenses.value.reduce((acc, e) => acc + Number(e.amount || 0), 0)
  return sum.toLocaleString('ru-RU', {
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
