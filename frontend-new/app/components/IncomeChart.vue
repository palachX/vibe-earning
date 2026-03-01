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
import type { Income } from '~/stores/finance'

ChartJS.register(LineElement, PointElement, CategoryScale, LinearScale, Tooltip, Legend)

interface Props {
  items: Income[]
}

const props = defineProps<Props>()

// Группируем доходы по месяцам для графика
const monthlyData = computed(() => {
  const grouped: Record<string, number> = {}
  
  props.items.forEach(income => {
    const date = new Date(income.date)
    const monthKey = date.toLocaleDateString('ru-RU', { year: 'numeric', month: 'short' })
    
    if (!grouped[monthKey]) {
      grouped[monthKey] = 0
    }
    grouped[monthKey] += Number(income.amount)
  })
  
  // Сортируем по дате
  const sortedEntries = Object.entries(grouped).sort(([a], [b]) => {
    const dateA = new Date(a)
    const dateB = new Date(b)
    return dateA.getTime() - dateB.getTime()
  })
  
  return {
    labels: sortedEntries.map(([month]) => month),
    values: sortedEntries.map(([, value]) => value)
  }
})

const chartData = computed(() => ({
  labels: monthlyData.value.labels,
  datasets: [
    {
      label: 'Доходы',
      data: monthlyData.value.values,
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
          return `${context.dataset.label}: ${value} ₽`
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
      grid: { color: 'rgba(55,65,81,0.4)' },
      beginAtZero: true
    }
  }
}
</script>

<style scoped>
div {
  height: 260px;
}
</style>
