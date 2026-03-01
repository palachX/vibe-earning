<template>
  <ClientOnly>
    <div class="rounded-xl border border-slate-800 bg-slate-900/40 p-4">
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend
} from 'chart.js'
import type { ForecastItem } from '~/stores/finance'

ChartJS.register(LineElement, PointElement, CategoryScale, LinearScale, Tooltip, Legend)

interface Props {
  items: ForecastItem[]
}

const props = defineProps<Props>()

const chartData = computed(() => ({
  labels: props.items.map(i =>
    new Date(i.week_start).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
  ),
  datasets: [
    {
      label: 'Closing balance',
      data: props.items.map(i => Number(i.closing_balance)),
      borderColor: '#22c55e',
      backgroundColor: 'rgba(34,197,94,0.2)',
      tension: 0.3,
      fill: true
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      labels: {
        color: '#e5e7eb'
      }
    },
    tooltip: {
      callbacks: {
        label(context: any) {
          const value = Number(context.parsed.y || 0).toLocaleString('ru-RU', {
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
          })
          return `${context.dataset.label}: ${value}`
        }
      }
    }
  },
  scales: {
    x: {
      ticks: { color: '#9ca3af' },
      grid: { color: 'rgba(55,65,81,0.4)' }
    },
    y: {
      ticks: { color: '#9ca3af' },
      grid: { color: 'rgba(55,65,81,0.4)' }
    }
  }
}
</script>

<style scoped>
div {
  height: 260px;
}
</style>

