## V3 — Checkout и транзакции

В третьей версии проекта добавлено оформление заказа через PostgreSQL-транзакцию.

### Реализовано

- endpoint `POST /checkout`;
- метод `PlaceOrder(userEmail string, items []models.OrderItem) (int, error)`;
- создание заказа внутри транзакции;
- проверка пользователя по email;
- проверка наличия товара;
- уменьшение `stock`;
- создание `order_items`;
- расчёт и обновление `total_price`;
- создание записи в таблице `transactions`;
- `COMMIT` при успешном оформлении;
- `ROLLBACK` при любой ошибке;
- ошибка `ErrInsufficientStock` при недостаточном количестве товара;
- отдельный пакет `handlers/checkout`;
- ручной mock и unit-тесты checkout handler.

### Endpoint

```text
POST /checkout
```


### Команды

Установка зависимостей:

```bash
go mod tidy
```

Запуск проекта:

```bash
docker compose up --build
```

Запуск тестов:

```bash
go test ./...
```

Подробный вывод тестов:

```bash
go test -v ./...
```

Линтер:

```bash
golangci-lint run ./...
```

Форматирование:

```bash
go fmt ./...
```

Остановка Docker:

```bash
docker compose down
```

Остановка с удалением данных PostgreSQL:

```bash
docker compose down -v
```