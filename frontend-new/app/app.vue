<template>
  <UApp>
    <div class="min-h-screen bg-gray-950 text-gray-50">
      <UHeader class="border-b border-gray-800/60 bg-gray-900/80 backdrop-blur">
        <template #left>
          <NuxtLink to="/" class="text-lg font-semibold tracking-tight">
            Vibe Earning
          </NuxtLink>
        </template>
        <UNavigationMenu :items="links" />
      </UHeader>

      <UMain>
        <UContainer class="py-6">
          <NuxtPage />
        </UContainer>
      </UMain>

    </div>
  </UApp>
</template>

<script setup lang="ts">
const links = [
  { label: 'Дашборд', to: '/' },
  { label: 'Доходы', to: '/incomes' },
  { label: 'Расходы', to: '/expenses' },
  { label: 'Постоянные траты', to: '/recurring' }
]
</script>
