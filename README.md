## V2 — Заказы и позиции заказа

Во второй версии проекта добавлена работа с заказами и позициями заказа[cite: 8].

---

### Заказы

Реализовано:

- Создание заказа[cite: 8];
- Получение заказа по ID[cite: 8];
- Получение заказов пользователя по email[cite: 8];
- Проверка входных данных;
- Обработка ошибок `400`, `404` и `500`.

**Методы storage:**

```go
CreateOrder(order models.Order) (int, error)
GetOrderByID(id int) (models.Order, error)
GetOrdersByUserEmail(email string) ([]models.Order, error)
```

**HTTP endpoints:**

- `POST /orders`[cite: 8]
- `GET /orders/{id}`[cite: 8]
- `GET /users/orders?email=example@example.com`[cite: 8]

При создании нового заказа:

- `total_price = 0`

Общая стоимость заказа не передаётся клиентом, а рассчитывается после добавления позиций.

---

### Позиции заказа

Реализовано:

- Добавление товара в заказ[cite: 8];
- Получение всех позиций заказа[cite: 8];
- Проверка `product_id` и `quantity`;
- Получение `order_id` из URL;
- Автоматический пересчёт общей стоимости заказа.

**Методы storage:**

```go
AddOrderItem(orderItem models.OrderItem) error
GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error)
UpdateOrderTotalPrice(orderID int) error
```

**HTTP endpoints:**

- `POST /orders/{id}/items`[cite: 8]
- `GET /orders/{id}/items`

Пример добавления позиции:

`POST /orders/1/items`

**Тело запроса:**

```json
{
  "ProductID": 1,
  "Quantity": 2
}
```

`OrderID` берётся из URL и записывается в `OrderItem` внутри handler.

---

### Пересчёт стоимости заказа

После добавления позиции вызывается:

```go
UpdateOrderTotalPrice(orderID int) error
```

Стоимость рассчитывается на основе:

$$product.price \times order\_items.quantity$$

Для расчёта используется SQL-запрос с:

- `JOIN`;
- `SUM`;
- `COALESCE`.

Если в заказе нет позиций, итоговая стоимость равна 0.

---

### Что изучается в V2

- Работа со связанными таблицами;
- Внешние ключи;
- SQL `JOIN`[cite: 8];
- Получение связанных данных[cite: 8];
- Работа с заказами и позициями заказа[cite: 8];
- Query-параметры в HTTP;
- Параметры URL;
- Разделение ошибок `400`, `404` и `500`;
- Тестирование handlers с ручными mock-объектами.

---

### Тесты и линтинг

Проверка всей версии:

```bash
go test ./...
```

Проверка линтером:

```bash
golangci-lint run ./...
```

Форматирование:

```bash
go fmt ./...
```

---

### Запуск

Проект запускается через Docker Compose:

```bash
docker compose up --build
```