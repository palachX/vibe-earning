<template>
  <form class="space-y-3" @submit.prevent="onSubmit">
    <div class="grid grid-cols-1 gap-3 sm:grid-cols-5">
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Название</label>
        <input
          v-model="name"
          type="text"
          required
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-amber-500"
        >
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Сумма</label>
        <input
          v-model="amount"
          type="number"
          min="0"
          step="0.01"
          required
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-amber-500"
        >
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Частота</label>
        <select
          v-model="frequency"
          required
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-amber-500"
        >
          <option value="monthly">
            Ежемесячно
          </option>
          <option value="weekly">
            Еженедельно
          </option>
        </select>
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Дата начала</label>
        <input
          v-model="startDate"
          type="date"
          required
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-amber-500"
        >
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Дата окончания (опционально)</label>
        <input
          v-model="endDate"
          type="date"
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-amber-500"
        >
      </div>
    </div>
    <button
      type="submit"
      class="inline-flex items-center rounded-md bg-amber-400 px-4 py-2 text-sm font-medium text-amber-950 hover:bg-amber-300 disabled:opacity-60"
      :disabled="submitting"
    >
      {{ submitting ? 'Сохранение…' : 'Добавить постоянную трату' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useFinanceStore } from '~/stores/finance'

const store = useFinanceStore()

const today = new Date().toISOString().slice(0, 10)
const name = ref('')
const amount = ref('')
const frequency = ref<'weekly' | 'monthly'>('monthly')
const startDate = ref(today)
const endDate = ref<string | null>(null)
const submitting = ref(false)

async function onSubmit() {
  if (!amount.value || Number(amount.value) <= 0 || !name.value) return
  submitting.value = true
  try {
    await store.addRecurring({
      name: name.value,
      amount: Number(amount.value).toFixed(2),
      frequency: frequency.value,
      start_date: startDate.value,
      end_date: endDate.value || null
    })
    name.value = ''
    amount.value = ''
    endDate.value = null
  } finally {
    submitting.value = false
  }
}
</script>

