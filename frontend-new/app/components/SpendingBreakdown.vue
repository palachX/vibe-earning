<template>
  <div class="space-y-6">
    <!-- Верхняя строка: баланс → обязательства → свободные -->
    <div class="grid gap-4 md:grid-cols-3">
      <div class="rounded-xl border border-slate-800 bg-gradient-to-br from-slate-900/80 to-slate-900/40 p-4">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/15 text-emerald-400">
            <UIcon name="i-lucide-wallet" class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs uppercase tracking-wide text-slate-400">
              Текущий баланс
            </p>
            <p class="text-xl font-semibold text-emerald-400">
              {{ formatMoney(currentBalance) }}
            </p>
          </div>
        </div>
      </div>

      <div class="rounded-xl border border-slate-800 bg-gradient-to-br from-slate-900/80 to-slate-900/40 p-4">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-500/15 text-amber-400">
            <UIcon name="i-lucide-repeat" class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs uppercase tracking-wide text-slate-400">
              Будущие обязательства
            </p>
            <p class="text-xl font-semibold text-amber-400">
              − {{ formatMoney(totalObligations) }}
            </p>
            <p class="text-xs text-slate-500">
              {{ weeks }} нед. · {{ recurringBreakdown.length }} статей
            </p>
          </div>
        </div>
      </div>

      <div class="rounded-xl border border-sky-800/60 bg-gradient-to-br from-sky-950/40 to-slate-900/40 p-4 ring-1 ring-sky-500/20">
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-sky-500/15 text-sky-400">
            <UIcon name="i-lucide-piggy-bank" class="h-5 w-5" />
          </div>
          <div>
            <p class="text-xs uppercase tracking-wide text-slate-400">
              Свободные деньги
            </p>
            <p class="text-xl font-semibold text-sky-300">
              {{ formatMoney(freeMoney) }}
            </p>
            <p class="text-xs text-slate-500">
              после вычета recurring
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="grid gap-6 lg:grid-cols-5">
      <!-- Диаграмма распределения -->
      <div class="rounded-xl border border-slate-800 bg-slate-900/40 p-4 lg:col-span-2">
        <p class="mb-3 text-sm font-medium text-slate-200">
          Распределение обязательств
        </p>
        <ClientOnly>
          <div v-if="recurringBreakdown.length" class="mx-auto h-52 w-52">
            <Doughnut :data="chartData" :options="chartOptions" />
          </div>
          <div v-else class="flex h-52 items-center justify-center text-sm text-slate-500">
            Нет постоянных трат
          </div>
        </ClientOnly>
      </div>

      <!-- Список постоянных трат -->
      <div class="rounded-xl border border-slate-800 bg-slate-900/40 p-4 lg:col-span-3">
        <div class="mb-4 flex items-center justify-between">
          <p class="text-sm font-medium text-slate-200">
            Постоянные обязательства
          </p>
          <UButton
            size="xs"
            variant="ghost"
            to="/recurring"
            icon="i-lucide-arrow-right"
          >
            Управление
          </UButton>
        </div>

        <div
          v-if="!recurringBreakdown.length"
          class="rounded-lg border border-dashed border-slate-700 px-4 py-8 text-center text-sm text-slate-500"
        >
          Добавьте постоянные траты, чтобы увидеть разбивку
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="(item, idx) in recurringBreakdown"
            :key="item.id"
            class="group rounded-lg border border-slate-800/80 bg-slate-950/40 px-4 py-3 transition-colors hover:border-slate-700"
          >
            <div class="mb-2 flex items-center justify-between gap-3">
              <div class="flex min-w-0 items-center gap-2">
                <span
                  class="h-2.5 w-2.5 shrink-0 rounded-full"
                  :style="{ backgroundColor: chartColors[idx % chartColors.length] }"
                />
                <p class="truncate text-sm font-medium text-slate-100">
                  {{ item.name }}
                </p>
                <UBadge
                  size="xs"
                  variant="subtle"
                  color="neutral"
                >
                  {{ humanFrequency(item.frequency) }}
                </UBadge>
              </div>
              <div class="shrink-0 text-right">
                <p class="text-sm font-semibold text-amber-300">
                  {{ formatMoney(item.total) }}
                </p>
                <p class="text-xs text-slate-500">
                  {{ item.percent.toFixed(1) }}%
                </p>
              </div>
            </div>
            <div class="h-1.5 overflow-hidden rounded-full bg-slate-800">
              <div
                class="h-full rounded-full transition-all duration-500"
                :style="{
                  width: `${item.percent}%`,
                  backgroundColor: chartColors[idx % chartColors.length]
                }"
              />
            </div>
            <p class="mt-1.5 text-xs text-slate-500">
              {{ formatMoney(item.amount) }} за платёж
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Ближайшая неделя -->
    <div
      v-if="currentWeek"
      class="rounded-xl border border-slate-800 bg-slate-900/40 p-4"
    >
      <div class="mb-4 flex flex-wrap items-center justify-between gap-2">
        <div>
          <p class="text-sm font-medium text-slate-200">
            Ближайшая неделя
          </p>
          <p class="text-xs text-slate-500">
            с {{ weekStartLabel }}
          </p>
        </div>
        <UBadge variant="subtle" color="neutral">
          Opening → Closing
        </UBadge>
      </div>

      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
        <div class="rounded-lg border border-emerald-900/40 bg-emerald-950/20 px-4 py-3">
          <p class="text-xs uppercase text-emerald-500/80">
            Доходы
          </p>
          <p class="mt-1 text-lg font-semibold text-emerald-400">
            + {{ formatMoney(currentWeek.income_total) }}
          </p>
        </div>

        <div class="rounded-lg border border-rose-900/40 bg-rose-950/20 px-4 py-3">
          <p class="text-xs uppercase text-rose-500/80">
            Разовые расходы
          </p>
          <p class="mt-1 text-lg font-semibold text-rose-400">
            − {{ formatMoney(currentWeek.expense_total) }}
          </p>
        </div>

        <div class="rounded-lg border border-amber-900/40 bg-amber-950/20 px-4 py-3">
          <p class="text-xs uppercase text-amber-500/80">
            Постоянные
          </p>
          <p class="mt-1 text-lg font-semibold text-amber-400">
            − {{ formatMoney(currentWeek.recurring_total) }}
          </p>
        </div>

        <div class="rounded-lg border border-slate-700 bg-slate-950/40 px-4 py-3">
          <p class="text-xs uppercase text-slate-500">
            На начало
          </p>
          <p class="mt-1 text-lg font-semibold text-slate-200">
            {{ formatMoney(currentWeek.opening_balance) }}
          </p>
        </div>

        <div class="rounded-lg border border-sky-800/50 bg-sky-950/20 px-4 py-3">
          <p class="text-xs uppercase text-sky-500/80">
            На конец
          </p>
          <p class="mt-1 text-lg font-semibold text-sky-300">
            {{ formatMoney(currentWeek.closing_balance) }}
          </p>
        </div>
      </div>

      <!-- Визуальная полоса движения денег за неделю -->
      <div class="mt-4 rounded-lg border border-slate-800 bg-slate-950/30 p-3">
        <p class="mb-2 text-xs uppercase tracking-wide text-slate-500">
          Движение за неделю
        </p>
        <div class="flex flex-wrap items-center gap-2 text-sm">
          <span class="rounded-md bg-slate-800 px-2 py-1 text-slate-200">
            {{ formatMoney(currentWeek.opening_balance) }}
          </span>
          <UIcon name="i-lucide-plus" class="h-4 w-4 text-emerald-500" />
          <span class="rounded-md bg-emerald-950/50 px-2 py-1 text-emerald-400">
            {{ formatMoney(currentWeek.income_total) }}
          </span>
          <UIcon name="i-lucide-minus" class="h-4 w-4 text-rose-500" />
          <span class="rounded-md bg-rose-950/50 px-2 py-1 text-rose-400">
            {{ formatMoney(currentWeek.expense_total) }}
          </span>
          <UIcon name="i-lucide-minus" class="h-4 w-4 text-amber-500" />
          <span class="rounded-md bg-amber-950/50 px-2 py-1 text-amber-400">
            {{ formatMoney(currentWeek.recurring_total) }}
          </span>
          <UIcon name="i-lucide-equal" class="h-4 w-4 text-sky-500" />
          <span class="rounded-md bg-sky-950/50 px-2 py-1 font-semibold text-sky-300">
            {{ formatMoney(currentWeek.closing_balance) }}
          </span>
        </div>
      </div>

      <!-- Пояснение почему меняется closing -->
      <div class="mt-4 rounded-lg border border-slate-800 bg-slate-900/30 p-3">
        <div class="mb-3 flex flex-wrap items-center justify-between gap-3">
          <p class="text-xs font-semibold uppercase tracking-wide text-slate-300">
            Почему меняется баланс
          </p>
          <UBadge
            variant="subtle"
            color="neutral"
            :class="weekNetDelta >= 0 ? 'text-emerald-300' : 'text-rose-300'"
          >
            {{ weekNetDelta >= 0 ? '+' : '−' }}{{ formatMoney(Math.abs(weekNetDelta)) }}
          </UBadge>
        </div>

        <p class="text-sm text-slate-200">
          Closing = Opening + Доходы − Разовые расходы − Постоянные траты
        </p>
        <div class="mt-2 flex flex-wrap items-center gap-2 text-sm">
          <span class="rounded-md bg-slate-800 px-2 py-1 text-slate-200">
            {{ formatMoney(currentWeek.opening_balance) }}
          </span>
          <UIcon name="i-lucide-plus" class="h-4 w-4 text-emerald-500" />
          <span class="rounded-md bg-emerald-950/50 px-2 py-1 text-emerald-400">
            {{ formatMoney(currentWeek.income_total) }}
          </span>
          <span class="text-slate-500">−</span>
          <span class="rounded-md bg-rose-950/50 px-2 py-1 text-rose-400">
            {{ formatMoney(currentWeek.expense_total) }}
          </span>
          <span class="text-slate-500">−</span>
          <span class="rounded-md bg-amber-950/50 px-2 py-1 text-amber-400">
            {{ formatMoney(currentWeek.recurring_total) }}
          </span>
          <UIcon name="i-lucide-equal" class="h-4 w-4 text-sky-500" />
          <span class="rounded-md bg-sky-950/50 px-2 py-1 font-semibold text-sky-300">
            {{ formatMoney(currentWeek.closing_balance) }}
          </span>
        </div>
      </div>

      <!-- Какие recurring сработают в ближайшую неделю -->
      <div class="mt-4 rounded-lg border border-slate-800 bg-slate-900/30 p-3">
        <div class="mb-2 flex items-center justify-between gap-3">
          <p class="text-xs font-semibold uppercase tracking-wide text-slate-300">
            Какие постоянные траты сработают
          </p>
          <UBadge variant="subtle" color="neutral">
            {{ recurringInWeek.length }} платеж(ей)
          </UBadge>
        </div>

        <ul v-if="recurringInWeek.length" class="space-y-1.5">
          <li
            v-for="occ in recurringInWeek"
            :key="`${occ.id}-${occ.date}`"
            class="flex items-center justify-between gap-3 rounded-md border border-slate-800/60 bg-slate-950/25 px-3 py-2"
          >
            <div class="min-w-0">
              <p class="truncate text-sm text-slate-100">
                {{ occ.name }}
              </p>
              <p class="text-xs text-slate-500">
                {{ formatShortDate(occ.date) }}
              </p>
            </div>
            <p class="shrink-0 text-sm font-semibold text-amber-300">
              − {{ formatMoney(occ.amount) }}
            </p>
          </li>
        </ul>

        <div v-else class="text-sm text-slate-500">
          В ближайшую неделю recurring-платежей не запланировано.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend
} from 'chart.js'
import type { ForecastItem, RecurringExpense } from '~/stores/finance'
import {
  buildRecurringBreakdown,
  getRecurringOccurrencesInWeek,
  formatMoney,
  humanFrequency
} from '~/composables/useSpendingBreakdown'

ChartJS.register(ArcElement, Tooltip, Legend)

interface Props {
  currentBalance: string | null
  freeMoney: string | null
  recurring: RecurringExpense[]
  forecast: ForecastItem[]
  weeks?: number
}

const props = withDefaults(defineProps<Props>(), {
  weeks: 52
})

const chartColors = [
  '#f59e0b',
  '#f97316',
  '#eab308',
  '#fb7185',
  '#a78bfa',
  '#38bdf8',
  '#34d399',
  '#f472b6'
]

const recurringBreakdown = computed(() =>
  buildRecurringBreakdown(props.recurring, props.weeks)
)

const totalObligations = computed(() =>
  recurringBreakdown.value.reduce((acc, r) => acc + r.total, 0)
)

const currentWeek = computed(() => props.forecast[0] ?? null)

const weekNetDelta = computed(() => {
  if (!currentWeek.value) return 0
  return (
    Number(currentWeek.value.income_total || 0)
    - Number(currentWeek.value.expense_total || 0)
    - Number(currentWeek.value.recurring_total || 0)
  )
})

const weekStartLabel = computed(() => {
  if (!currentWeek.value) return '—'
  return new Date(currentWeek.value.week_start).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
})

const recurringInWeek = computed(() => {
  if (!currentWeek.value) return []
  return getRecurringOccurrencesInWeek(props.recurring, currentWeek.value.week_start)
})

const chartData = computed(() => ({
  labels: recurringBreakdown.value.map(r => r.name),
  datasets: [
    {
      data: recurringBreakdown.value.map(r => r.total),
      backgroundColor: recurringBreakdown.value.map((_, i) => chartColors[i % chartColors.length]),
      borderColor: '#0f172a',
      borderWidth: 2,
      hoverOffset: 6
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '62%',
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      callbacks: {
        label(context: { label?: string, parsed?: number, dataset: { data: number[] } }) {
          const total = context.dataset.data.reduce((a, b) => a + b, 0)
          const value = context.parsed || 0
          const pct = total > 0 ? ((value / total) * 100).toFixed(1) : '0'
          return `${context.label}: ${formatMoney(value)} (${pct}%)`
        }
      }
    }
  }
}

function formatShortDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
}
</script>
