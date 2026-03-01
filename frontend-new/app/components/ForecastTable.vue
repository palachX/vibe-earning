<template>
  <div class="overflow-hidden rounded-xl border border-slate-800 bg-slate-900/40">
    <table class="min-w-full divide-y divide-slate-800 text-sm">
      <thead class="bg-slate-900/60">
        <tr>
          <th class="px-4 py-2 text-left text-xs font-semibold uppercase tracking-wide text-slate-300">
            Неделя
          </th>
          <th class="px-4 py-2 text-right text-xs font-semibold uppercase tracking-wide text-slate-300">
            Opening
          </th>
          <th class="px-4 py-2 text-right text-xs font-semibold uppercase tracking-wide text-slate-300">
            Closing
          </th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-800">
        <tr v-for="item in limitedItems" :key="item.week_start" class="hover:bg-slate-800/40">
          <td class="px-4 py-2 text-slate-100">
            {{ formatDate(item.week_start) }}
          </td>
          <td class="px-4 py-2 text-right text-slate-200">
            {{ formatAmount(item.opening_balance) }}
          </td>
          <td class="px-4 py-2 text-right text-slate-200">
            {{ formatAmount(item.closing_balance) }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { ForecastItem } from '~/stores/finance'

interface Props {
  items: ForecastItem[]
  limit?: number
}

const props = defineProps<Props>()

const limitedItems = computed(() =>
  props.limit ? props.items.slice(0, props.limit) : props.items
)

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU')
}

function formatAmount(a: string) {
  return Number(a).toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}
</script>

