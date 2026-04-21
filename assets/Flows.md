## OwnServiceMethods (customer/restaurant/courier)

* SaveWithPassword
* LoadByWalletAddress
* Get_Sevice_WalletAddress


## Order:

### Customer

- `CreateOrder` — `POST /orders` (status: CUSTOMER_CREATED)
- `CreateOrderWithItems` — `POST /orders/with-items` (status: CUSTOMER_CREATED + items)
- `AddItemIntoOrder` — `POST /orders/{id}/items`
- `RemoveItemFromOrder` — `DELETE /orders/{id}/items`
- `PayOrder` — `POST /orders/{id}/pay` (→ CUSTOMER_PAID через WalletClient)
- `CalculateOrderTotal` — `GET /orders/{id}/total`
- `GetCustomerWalletAddress` — `GET /customers/{id}/wallet`


###  Restaurant

- `ListRestaurants` — `GET /restaurants`
- `GetMenu` — `GET /menu?restaurant_id=...`
- `PUT /orders/{id}/status` (KITCHEN_ACCEPTED / KITCHEN_DENIED)
- `PUT /orders/{id}/status` (KITCHEN_PREPARING → DELIVERY_PENDING)

### Courier

- `ListCouriers` — `GET /couriers`
- `POST /orders/{id}/accept` (DELIVERY_PICKING)
- `PUT /orders/{id}/status` (DELIVERY_DELIVERING)
- `PUT /orders/{id}/status` (ORDER_COMPLETED)


ц## Order pipeline

| # | Текущий статус | Следующий (✅) | Альтернативные (❌) | Роль | Действие |
|---|---|---|---|---|---|
| 1 | `CUSTOMER_CREATED` | `CUSTOMER_PAID` | `CUSTOMER_CANCELLED` | Customer | PayOrder / Cancel |
| 2 | `CUSTOMER_PAID` | `KITCHEN_ACCEPTED` | `KITCHEN_DENIED` | Restaurant | AcceptOrder (total>0) / Deny (total≤0) |
| 3 | `KITCHEN_DENIED` | `KITCHEN_REFUNDED` | — | Restaurant | UpdateOrderStatus |
| 4 | `KITCHEN_ACCEPTED` | `KITCHEN_PREPARING` | — | Restaurant | UpdateOrderStatus |
| 5 | `KITCHEN_PREPARING` | `DELIVERY_PENDING` | — | Restaurant | UpdateOrderStatus |
| 6 | `DELIVERY_PENDING` | `DELIVERY_PICKING` | `DELIVERY_DENIED` | Courier | AcceptOrder / Deny |
| 7 | `DELIVERY_PICKING` | `DELIVERY_DELIVERING` | — | Courier | UpdateOrderStatus |
| 8 | `DELIVERY_DELIVERING` | `DELIVERY_COMPLETE` | — | Courier | Complete / Fail |
| 9 | `DELIVERY_COMPLETE` | — | — | Courier | Complete / Fail |
