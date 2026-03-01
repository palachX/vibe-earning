<template>
  <UPage>
    <UPageHeader
      title="Доходы"
      description="Фиксация всех поступлений и их влияние на прогноз."
    >
      <template #extra>
        <UButton
          color="primary"
          size="sm"
          icon="i-lucide-refresh-ccw"
          @click="reload"
        >
          Обновить
        </UButton>
      </template>
    </UPageHeader>

    <UPageBody>
      <div class="space-y-6">
        <UCard>
          <template #header>
            <p class="text-sm font-medium">
              Новый доход
            </p>
          </template>
          <IncomeForm />
        </UCard>

        <UCard>
          <template #header>
            <p class="text-sm font-medium">
              Список доходов
            </p>
          </template>
          <div
            v-if="!incomes.length"
            class="text-sm text-gray-400"
          >
            Пока нет данных.
          </div>
          <div v-else>
            <UTable
              :data="incomes"
              :columns="columns"
              :row-class="'text-sm'"
            />
          </div>
        </UCard>
      </div>
    </UPageBody>
  </UPage>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

import { useFinanceStore } from '~/stores/finance'
import IncomeForm from '~/components/IncomeForm.vue'
import { UButton } from '#components'

const store = useFinanceStore()
const { incomes } = storeToRefs(store)

const columns = [
  { accessorKey: 'date', header: 'Дата', cell: ({ row }) => {
    return formatDate(row.getValue('date'))
  } },
  { accessorKey: 'description', header: 'Описание' },
  { accessorKey: 'amount', header: 'Сумма', meta: { class: 'text-right' }, cell: ({ row }) => {
    return formatAmount(row.getValue('amount'))
  } },
  { accessorKey: 'actions', header: '', meta: { class: 'w-10 text-right' }, cell: ({ row }) => {
    return h(UButton, {
      icon: 'i-lucide-trash-2',
      color: 'error',
      variant: 'ghost',
      size: 'xs',
      onClick: () => remove(row.original.id)
    })
  } }
]

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('ru-RU')
}

function formatAmount(a: string) {
  return Number(a).toLocaleString('ru-RU', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

async function reload() {
  await store.loadAll()
}

async function remove(id: number) {
  await store.deleteIncome(id)
}

onMounted(async () => {
  await store.loadAll()
})
</script>
