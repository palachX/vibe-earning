export interface Income {
  id: number
  user_id: number
  date: string
  amount: string
  description?: string
}

export interface Expense {
  id: number
  user_id: number
  date: string
  amount: string
  description?: string
}

export interface RecurringExpense {
  id: number
  user_id: number
  name: string
  amount: string
  frequency: string
  start_date: string
  end_date?: string | null
}

export interface ForecastItem {
  week_start: string
  opening_balance: string
  closing_balance: string
}

export const useFinanceStore = defineStore('finance', () => {
  // state
  const incomes = ref<Income[]>([])
  const expenses = ref<Expense[]>([])
  const recurring = ref<RecurringExpense[]>([])
  const forecast = ref<ForecastItem[]>([])
  const freeMoney = ref<string | null>(null)
  const currentBalance = ref<string | null>(null)
  const loading = ref(false)

  // api + toast (можно вынести выше, чтобы не создавать каждый раз)
  const api = useApi()
  const toast = useToast()
  // actions

  async function loadAll() {
    loading.value = true
    try {
      const [i, e, r] = await Promise.all([
        api.get<Income[]>('/incomes'),
        api.get<Expense[]>('/expenses'),
        api.get<RecurringExpense[]>('/recurring-expenses')
      ])
      incomes.value = i
      expenses.value = e
      recurring.value = r
    } finally {
      loading.value = false
    }
  }

  async function addIncome(payload: { date: string, amount: string, description?: string }) {
    loading.value = true
    try {
      const created = await api.post<Income>('/incomes', {
        method: 'POST',
        body: payload
      })
      incomes.value.unshift(created)
      await refreshDashboard()
      toast.add({ title: 'Доход добавлен', color: 'success' })
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось добавить доход', color: 'error' })
      throw e
    } finally {
      loading.value = false
    }
  }

  async function addExpense(payload: { date: string, amount: string, description?: string }) {
    loading.value = true
    try {
      const created = await api.post<Expense>('/expenses', {
        method: 'POST',
        body: payload
      })
      expenses.value.unshift(created)
      await refreshDashboard()
      toast.add({ title: 'Расход добавлен', color: 'success' })
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось добавить расход', color: 'error' })
      throw e
    } finally {
      loading.value = false
    }
  }

  async function addRecurring(payload: {
    name: string
    amount: string
    frequency: string
    start_date: string
    end_date?: string | null
  }) {
    loading.value = true
    try {
      const created = await api.post<RecurringExpense>('/recurring-expenses', {
        method: 'POST',
        body: payload
      })
      recurring.value.unshift(created)
      await refreshDashboard()
      toast.add({ title: 'Постоянная трата добавлена', color: 'success' })
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось добавить постоянную трату', color: 'error' })
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteIncome(id: number) {
    try {
      await api.del(`/incomes/${id}`, { method: 'DELETE' })
      incomes.value = incomes.value.filter(i => i.id !== id)
      await refreshDashboard()
      toast.add({ title: 'Доход удалён', color: 'success' })
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось удалить доход', color: 'error' })
    }
  }

  async function deleteExpense(id: number) {
    try {
      await api.del(`/expenses/${id}`, { method: 'DELETE' })
      expenses.value = expenses.value.filter(i => i.id !== id)
      await refreshDashboard()
      toast.add({ title: 'Расход удалён', color: 'success' })
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось удалить расход', color: 'error' })
    }
  }

  async function deleteRecurring(id: number) {
    try {
      await api.del(`/recurring-expenses/${id}`, { method: 'DELETE' })
      recurring.value = recurring.value.filter(i => i.id !== id)
      await refreshDashboard()
      toast.add({ title: 'Постоянная трата удалена', color: 'success' })
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось удалить запись', color: 'error' })
    }
  }

  async function fetchForecast(weeks = 52) {
    try {
      forecast.value = await api.get<ForecastItem[]>(`/forecast?weeks=${weeks}`)
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось загрузить прогноз', color: 'error' })
      throw e
    }
  }

  async function fetchFreeMoney(weeks = 52) {
    try {
      const res = await api.get<{ free_money: string }>(`/free-money?weeks=${weeks}`)
      freeMoney.value = res.free_money
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось загрузить свободные деньги', color: 'error' })
      throw e
    }
  }

  async function fetchCurrentBalance() {
    try {
      const res = await api.get<{ current_balance: string }>(`/current-balance`)
      currentBalance.value = res.current_balance
    } catch (e: any) {
      toast.add({ title: e?.data?.error || 'Не удалось загрузить текущий баланс', color: 'error' })
      throw e
    }
  }

  async function refreshDashboard() {
    loading.value = true
    try {
      await Promise.all([
        fetchForecast(52),
        fetchFreeMoney(52),
        fetchCurrentBalance()
      ])
    } finally {
      loading.value = false
    }
  }

  return {
    incomes,
    expenses,
    recurring,
    forecast,
    freeMoney,
    currentBalance,
    loading,
    loadAll,
    addIncome,
    addExpense,
    addRecurring,
    deleteIncome,
    deleteExpense,
    deleteRecurring,
    fetchForecast,
    fetchFreeMoney,
    fetchCurrentBalance,
    refreshDashboard
  }
})
