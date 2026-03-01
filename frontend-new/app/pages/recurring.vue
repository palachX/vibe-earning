<template>
  <UPage>
    <UPageHeader
      title="Постоянные траты"
      description="Подписки, аренда и другие регулярные обязательства."
    >
      <template #extra>
        <UButton
          color="amber"
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
              Новая постоянная трата
            </p>
          </template>
          <RecurringForm />
        </UCard>

        <UCard>
          <template #header>
            <p class="text-sm font-medium">
              Список постоянных трат
            </p>
          </template>

          <div
            v-if="!recurring.length"
            class="text-sm text-gray-400"
          >
            Пока нет данных.
          </div>
          <div v-else>
            <UTable
              :data="recurring"
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
import RecurringForm from '~/components/RecurringForm.vue'
import { UButton } from '#components'

const store = useFinanceStore()
const { recurring } = storeToRefs(store)

const columns = [
  { accessorKey: 'name', header: 'Название' },
  { accessorKey: 'frequency', header: 'Частота', cell: ({ row }) => {
    return humanFrequency(row.getValue('frequency'))
  } },
  { accessorKey: 'start_date', header: 'Старт', cell: ({ row }) => {
    return formatDate(row.getValue('start_date'))
  } },
  { accessorKey: 'end_date', header: 'Окончание', cell: ({ row }) => {
    const val = row.getValue('end_date')
    return val ? formatDate(val) : '—'
  } },
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

function humanFrequency(f: string) {
  if (f === 'weekly') return 'Еженедельно'
  if (f === 'monthly') return 'Ежемесячно'
  return f
}

async function reload() {
  await store.loadAll()
}

async function remove(id: number) {
  await store.deleteRecurring(id)
}

onMounted(async () => {
  await store.loadAll()
})
</script>
