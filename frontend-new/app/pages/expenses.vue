<template>
  <UPage>
    <UPageHeader
      title="Расходы"
      description="Все траты и их влияние на свободные деньги."
    >
      <template #extra>
        <UButton
          color="red"
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
              Новый расход
            </p>
          </template>
          <ExpenseForm />
        </UCard>

        <UCard>
          <template #header>
            <p class="text-sm font-medium">
              Таблица расходов
            </p>
          </template>

          <div
            v-if="!expenses.length"
            class="text-sm text-gray-400"
          >
            Пока нет данных.
          </div>
          <div v-else>
            <UTable
              :data="expenses"
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
import ExpenseForm from '~/components/ExpenseForm.vue'
import { UButton } from '#components'

const store = useFinanceStore()
const { expenses } = storeToRefs(store)

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
  await store.deleteExpense(id)
}

onMounted(async () => {
  await store.loadAll()
})
</script>
