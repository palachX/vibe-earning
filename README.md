## Vibe Earning — MVP системы учёта финансов

![Демонстрация](/assets/VibeEarning.gif)

> ⚠️ Приложение написано с использованием нейросетей для генерации кода.

Простой, но полностью рабочий продукт для персонального учёта финансов:

- **Backend**: Go 1.22 + Gin + PostgreSQL + `shopspring/decimal`
- **Frontend**: Nuxt 4 (Vue 3, Composition API, TypeScript) + Pinia + Nuxt UI + Chart.js
- **Инфраструктура**: Docker + docker-compose

Все расчёты (баланс, прогноз, свободные деньги) выполняются **на backend**, frontend только отправляет данные и отображает результат.

## Использование Taskfile

Для удобного управления проектом используется [Task](https://taskfile.dev/). Основные команды:

```bash
# Установка Task (если еще не установлен)
# Windows: scoop install task
# macOS: brew install go-task
# Linux: sudo snap install go-task --classic

# Запуск всех сервисов (postgres + backend + frontend)
task up

# Остановка всех сервисов
task down

# Пересборка и запуск
task restart

# Просмотр логов
task logs

# Только backend (с postgres)
task backend

# Только frontend
task frontend

# Миграции БД
task migrate

# Очистка
task clean
```

---


---

## Структура репозитория

- `cmd/server/main.go` — входная точка Go-сервера
- `internal/db` — подключение к PostgreSQL (`database/sql` + `pgx`)
- `internal/models` — модели (incomes, expenses, recurring, forecast)
- `internal/services` — доменная логика (создание записей, прогноз, free-money и т.д.)
- `internal/handlers` — HTTP-обработчики Gin (REST API)
- `migrations/001_init.sql` — SQL-миграция для PostgreSQL
- `docker-compose.yml` — запуск PostgreSQL, backend и frontend
- `Dockerfile` — билд backend-приложения
- `frontend/` — Nuxt 4 SPA (dashboard, формы, таблицы, график)

---

## Backend

### Технологии

- Go 1.22+
- Gin (HTTP router)
- PostgreSQL
- `database/sql` + `pgx` драйвер
- `github.com/shopspring/decimal` — точная работа с деньгами (без `float`)

### Модель данных (PostgreSQL)

Создаётся миграцией `migrations/001_init.sql` (автоматически при старте контейнера БД):

```sql
CREATE TABLE incomes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    date DATE NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    description TEXT
);

CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    date DATE NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    description TEXT
);

CREATE TABLE recurring_expenses (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name TEXT NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    frequency TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);
```

Во всех таблицах используется фиксированный `user_id = 1` (авторизация не реализована).

### Запуск backend через Docker

```bash
cd /path/to/vibe-earning
docker compose up --build
```

Что делает `docker-compose.yml`:

- стартует контейнер **PostgreSQL** с параметрами:
  - `POSTGRES_USER=app`
  - `POSTGRES_PASSWORD=app`
  - `POSTGRES_DB=appdb`
- применяет миграцию `migrations/001_init.sql`
- билдит контейнер **app** из `Dockerfile` и запускает Go-сервер

По умолчанию:

- backend доступен по адресу: `http://localhost:8080`
- строка подключения к БД: `postgres://app:app@db:5432/appdb?sslmode=disable`

Можно переопределить:

- `PORT` — порт HTTP-сервера
- `DATABASE_URL` — строка подключения к PostgreSQL

### Основная логика

- Используется `database/sql` с драйвером `pgx`
- Все операции записи (INSERT/DELETE) выполняются **через транзакции**
- Для денег используется `decimal.Decimal`, тип в БД — `NUMERIC(14,2)`
- Весь код работы с бизнес-логикой — в `internal/services/finance.go`

Основные операции:

- **Доходы / расходы / recurring**:
  - `CreateIncome`, `CreateExpense`, `CreateRecurringExpense`
  - `ListIncomes`, `ListExpenses`, `ListRecurring`
  - `DeleteIncome`, `DeleteExpense`, `DeleteRecurring`
- **Баланс**:
  - `GetCurrentBalance(today)` — сумма всех доходов − расходов − всех сработавших recurring до указанной даты (включительно)
- **Прогноз**:
  - `ForecastBalance(weeks)`:
    - стартовая точка — ближайший понедельник
    - неделя = понедельник–воскресенье
    - `opening_balance` первой недели = текущий баланс на день до `week_start`
    - для каждой недели: `closing = opening + incomes - expenses - recurring`
- **Свободные деньги**:
  - `FreeMoney(weeks)`:
    - берётся **текущий баланс на сегодня**
    - рассчитываются все будущие recurring в горизонте `weeks * 7` дней
    - `free_money = current_balance - future_recurring`

### REST API

Все суммы и даты передаются как **строки**:

- `date`: `"YYYY-MM-DD"`
- `amount`: `"1234.56"`

#### Доходы

- **POST `/incomes`** — добавить доход

  Тело:

  ```json
  {
    "date": "2026-02-01",
    "amount": "2500.00",
    "description": "Salary"
  }
  ```

- **GET `/incomes`** — список доходов
- **DELETE `/incomes/:id`** — удалить доход

#### Расходы

- **POST `/expenses`** — добавить расход

  ```json
  {
    "date": "2026-02-03",
    "amount": "120.00",
    "description": "Groceries"
  }
  ```

- **GET `/expenses`** — список расходов
- **DELETE `/expenses/:id`** — удалить расход

#### Постоянные траты (recurring)

- **POST `/recurring-expenses`**

  ```json
  {
    "name": "Rent",
    "amount": "800.00",
    "frequency": "monthly",
    "start_date": "2026-02-01"
  }
  ```

  - `frequency`: `"weekly"` или `"monthly"`
  - опционально: `"end_date": "YYYY-MM-DD"`

- **GET `/recurring-expenses`** — список постоянных трат
- **DELETE `/recurring-expenses/:id`** — удалить запись

#### Прогноз

- **GET `/forecast?weeks=52`**

  Ответ:

  ```json
  [
    {
      "week_start": "2026-02-02",
      "opening_balance": "1000.00",
      "closing_balance": "870.00"
    }
  ]
  ```

  Логика:

  - неделя: понедельник–воскресенье
  - итерация `weeks` недель вперёд
  - для каждой недели:
    - суммируются все доходы и расходы внутри недели
    - добавляются все recurring, которые попадают в эту неделю
    - считается `closing_balance`

#### Свободные деньги

- **GET `/free-money?weeks=52`**

  Ответ:

  ```json
  {
    "free_money": "1234.56",
    "weeks": 52
  }
  ```

  - текущий баланс (на сегодня) минус все будущие обязательные recurring в горизонте `weeks` недель

#### Текущий баланс

- **GET `/current-balance`**

  Ответ:

  ```json
  {
    "current_balance": "1500.00",
    "date": "2026-02-16"
  }
  ```

  - фактический баланс на сегодня (учтены доходы, расходы и уже сработавшие recurring)

---

## Frontend (Nuxt 4 + Nuxt UI)

### Технологии

- Nuxt 4 (SPA, `ssr: false`)
- Vue 3 + Composition API + TypeScript
- Pinia — состояние
- Nuxt UI — готовые UI-компоненты (страницы, таблицы, карточки, кнопки)
- TailwindCSS — базовая стилизация
- Chart.js + vue-chartjs — график прогноза

### Структура

```text
frontend/
  app.vue
  nuxt.config.ts
  tailwind.config.ts
  postcss.config.cjs
  assets/css/tailwind.css
  composables/
    useApi.ts
    useToast.ts
  components/
    ToastContainer.vue
    FreeMoneyCard.vue
    CurrentBalanceCard.vue
    ForecastTable.vue
    ForecastChart.vue
    IncomeForm.vue
    ExpenseForm.vue
    RecurringForm.vue
  stores/
    finance.ts
  pages/
    index.vue
    incomes.vue
    expenses.vue
    recurring.vue
```

### Конфигурация

`frontend/nuxt.config.ts`:

- `ssr: false` — чистый SPA-режим
- модули: `@pinia/nuxt`, `@nuxt/ui`
- runtime config:

  ```ts
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
    }
  }
  ```

Для работы с backend в другом домене/порте достаточно установить:

```bash
export NUXT_PUBLIC_API_BASE_URL="http://localhost:8080"
```

### Запуск frontend (pnpm)

```bash
cd frontend
pnpm install
pnpm run dev
```

По умолчанию:

- frontend: `http://localhost:3000`
- backend: `http://localhost:8080`

### Компосаблы

- `useApi.ts`:

  - обёртка над `ofetch` с `baseURL` из runtime config
  - методы: `get`, `post`, `del` — используются в Pinia-сторе

- `useToast.ts`:

  - глобальный стейт тостов через `useState`
  - методы: `success(message)`, `error(message)`
  - уведомления автоматически исчезают через несколько секунд

### Store (Pinia) — `stores/finance.ts`

**State:**

```ts
incomes: []
expenses: []
recurring: []
forecast: []
freeMoney: string | null
currentBalance: string | null
loading: boolean
```

**Actions:**

- `loadAll()` — грузит `incomes`, `expenses`, `recurring` (GET)
- `addIncome()`, `addExpense()`, `addRecurring()`:
  - отправляют POST запросы на backend
  - добавляют запись в список
  - вызывают `refreshDashboard()`
- `deleteIncome()`, `deleteExpense()`, `deleteRecurring()`:
  - DELETE на backend
  - обновляют локальный список
  - вызывают `refreshDashboard()`
- `fetchForecast(weeks)` — `GET /forecast`
- `fetchFreeMoney(weeks)` — `GET /free-money`
- `fetchCurrentBalance()` — `GET /current-balance`
- `refreshDashboard()` — параллельно вызывает `fetchForecast(52)`, `fetchFreeMoney(52)` и `fetchCurrentBalance()`

После любого добавления/удаления автоматически обновляются **прогноз** и **свободные деньги**.

### Основные страницы

#### Главная (`pages/index.vue`)

- Отображает:
  - `FreeMoneyCard` — крупная карточка со свободными деньгами (горизонт 52 недели)
  - `ForecastChart` — линейный график `closing_balance` по неделям
  - `ForecastTable` (первые 8 недель)
- Кнопка «Обновить дашборд» → `store.refreshDashboard()`
- При монтировании сразу запрашивает прогноз и свободные деньги

#### Доходы (`pages/incomes.vue`)

- Форма `IncomeForm`:
  - поля: `date`, `amount`, `description`
  - валидация: обязательные поля, `amount > 0`
- Таблица доходов:
  - дата, описание, сумма
  - кнопка «Удалить» рядом с каждой строкой
- При открытии страницы:
  - вызывает `store.loadAll()` (загружает все списки)

#### Расходы (`pages/expenses.vue`)

- Аналогично доходам, но с компонентом `ExpenseForm`
- Таблица расходов (дата, описание, сумма, удаление)

#### Постоянные траты (`pages/recurring.vue`)

- `RecurringForm`:
  - `name`, `amount`, `frequency (weekly/monthly)`, `start_date`
  - валидация: обязательные поля, `amount > 0`
- Таблица recurring:
  - название, частота, дата начала, сумма, кнопка «Удалить»

### UI / UX

- TailwindCSS:
  - тёмная тема (фон `bg-slate-950`, текст `text-slate-100`)
  - минималистичные карточки и таблицы
- Лоадер:
  - `loading` в сторе можно использовать для дизейбла кнопок и показа состояния «Обновление…»
- Toast уведомления:
  - успех/ошибка для всех основных операций (создание/удаление, загрузка прогноза и free-money)
  - отображаются в верхней части экрана через `ToastContainer`

---

## Примеры запросов (ручное тестирование backend)

Добавить доход:

```bash
curl -X POST http://localhost:8080/incomes \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-02-01","amount":"2500.00","description":"Salary"}'
```

Добавить расход:

```bash
curl -X POST http://localhost:8080/expenses \
  -H "Content-Type: application/json" \
  -d '{"date":"2026-02-03","amount":"120.00","description":"Groceries"}'
```

Добавить постоянную трату:

```bash
curl -X POST http://localhost:8080/recurring-expenses \
  -H "Content-Type: application/json" \
  -d '{"name":"Rent","amount":"800.00","frequency":"monthly","start_date":"2026-02-01"}'
```

Получить прогноз на 52 недели:

```bash
curl "http://localhost:8080/forecast?weeks=52"
```

Получить свободные деньги:

```bash
curl "http://localhost:8080/free-money?weeks=52"
```

---

## Итог

- Backend на Go:
  - Gin + PostgreSQL + `decimal`, корректная работа с деньгами, прогноз и free-money считаются **только на сервере**
- Frontend на Nuxt 3:
  - SPA-дэшборд с формами, таблицами, графиком и уведомлениями
- Всё поднимается стандартными командами:

```bash
# backend + postgres
docker compose up --build

# frontend
cd frontend
npm install
npm run dev
```

