## V2 — Orders and Order Items

Во второй версии проекта добавлена работа с заказами и позициями заказа.

### Реализовано

#### Orders
- Создание нового заказа.
- Получение заказа по ID.
- Получение заказов пользователя по email.
- Валидация входных данных.
- Обработка ошибок `400`, `404` и `500`.

**Методы storage:**
```go
CreateOrder(order models.Order) (int, error)
GetOrderByID(id int) (models.Order, error)
GetOrdersByUserEmail(email string) ([]models.Order, error)
```

**HTTP endpoints:**
- `POST /orders`
- `GET /orders/{id}`
- `GET /orders/user/{email}`

---

#### Order Items
- Добавление позиции в заказ.
- Получение всех позиций заказа.
- Автоматический пересчёт `total_price` после добавления товара.

**Методы storage:**
```go
AddOrderItem(orderItem models.OrderItem) error
GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error)
UpdateOrderTotalPrice(orderID int) error
```

**HTTP endpoints:**
- `POST /orders/{id}/items`
- `GET /orders/{id}/items`

---

#### Total Price
После добавления новой позиции в заказ автоматически пересчитывается общая стоимость заказа:
- Формула расчёта: `product price × quantity`
- Подсчёт выполняется SQL-запросом с `SUM` и `JOIN` таблиц `order_items` и `products`.

---

#### Tests
Для HTTP-хендлеров используются ручные моки storage:
- Успешный запрос.
- Некорректные входные данные (`400 Bad Request`).
- Отсутствие записи (`404 Not Found`).
- Внутренняя ошибка слоя storage (`500 Internal Server Error`).

---

### Проверка проекта

```bash
# Запуск тестов
go test -v ./...

# Запуск линтера
golangci-lint run ./...
```