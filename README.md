# V4 — Аналитика

Добавлены аналитические SQL-запросы с использованием `JOIN`, `GROUP BY` и `SUM`.

История заказов пользователя:

```text
GET /users/history?email=alice@example.com
```

Возвращает информацию о заказах пользователя, включая:

- ID заказа;
- товар;
- количество;
- итоговую стоимость заказа;
- статус транзакции;
- дату создания.

Популярные товары:

```text
GET /products/popular
```

## API

### Пользователи

```text
POST /users
GET  /users
GET  /users/{email}
GET  /users/orders?email=alice@example.com
GET  /users/history?email=alice@example.com
```

### Товары

```text
POST   /products
GET    /products
GET    /products/{id}
PUT    /products/{id}
DELETE /products/{id}
GET    /products/popular
```

### Заказы

```text
POST /orders
GET  /orders/{id}

POST /orders/{id}/items
GET  /orders/{id}/items
```

### Checkout

```text
POST /checkout
```

### Проверка состояния приложения

```text
GET /health
```

## Переменные окружения

Создайте файл `.env` на основе `.env.example`.

Пример:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
DB_HOST=postgres
DB_PORT=5432
```

## Запуск через Docker

Запуск проекта:

```bash
docker compose up --build
```

Остановка проекта:

```bash
docker compose down
```

Остановка проекта с удалением volume базы данных:

```bash
docker compose down -v
```

## Тесты

Запуск всех тестов:

```bash
go test ./...
```

## Линтер

```bash
golangci-lint run ./...
```

## Форматирование

```bash
go fmt ./...
```