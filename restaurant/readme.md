Functionality:
1. Register
2. Login
3. Accept order
4. Deny order
5. Open
6. Close
7. Edit menu

| Order Step | Action | Status Transition      | Request             | Response       |
|------------|--------|------------------------|---------------------|----------------|
| 1 | See new order   | CUSTOMER_PAID (already)| Get new order       |                |
| 2 | Deny order      | KITCHEN_DENIED         | Update order status | Precess refund |
| 3 | Accept order    | KITCHEN_ACCEPTED       | Update order status |                |
| 4 | Complete order  | KITCHEN_COMPLETED      | Update order status |                |

