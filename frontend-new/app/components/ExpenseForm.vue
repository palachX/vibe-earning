<template>
  <form class="space-y-3" @submit.prevent="onSubmit">
    <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Дата</label>
        <input
          v-model="date"
          type="date"
          required
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-rose-500"
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
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-rose-500"
        >
      </div>
      <div>
        <label class="mb-1 block text-xs font-medium text-slate-300">Описание</label>
        <input
          v-model="description"
          type="text"
          class="w-full rounded-md border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-rose-500"
        >
      </div>
    </div>
    <button
      type="submit"
      class="inline-flex items-center rounded-md bg-rose-500 px-4 py-2 text-sm font-medium text-rose-950 hover:bg-rose-400 disabled:opacity-60"
      :disabled="submitting"
    >
      {{ submitting ? 'Сохранение…' : 'Добавить расход' }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useFinanceStore } from '~/stores/finance'

const store = useFinanceStore()

const today = new Date().toISOString().slice(0, 10)
const date = ref(today)
const amount = ref('')
const description = ref('')
const submitting = ref(false)

async function onSubmit() {
  if (!amount.value || Number(amount.value) <= 0) return
  submitting.value = true
  try {
    await store.addExpense({
      date: date.value,
      amount: Number(amount.value).toFixed(2),
      description: description.value || undefined
    })
    amount.value = ''
    description.value = ''
  } finally {
    submitting.value = false
  }
}
</script>

