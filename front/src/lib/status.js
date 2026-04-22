export const orderStatuses = [
  'CUSTOMER_CREATED',
  'CUSTOMER_PAID',
  'CUSTOMER_CANCELLED',
  'KITCHEN_ACCEPTED',
  'KITCHEN_DENIED',
  'KITCHEN_PREPARING',
  'COURIER_REFUNDED',
  'DELIVERY_PENDING',
  'DELIVERY_PICKING',
  'DELIVERY_DENIED',
  'DELIVERY_REFUNDED',
  'DELIVERY_DELIVERING',
  'ORDER_COMPLETED',
  'ORDER_FAILED',
]

export const statusLabels = {
  CUSTOMER_CREATED: 'Created',
  CUSTOMER_PAID: 'Paid',
  CUSTOMER_CANCELLED: 'Cancelled',
  KITCHEN_ACCEPTED: 'Kitchen accepted',
  KITCHEN_DENIED: 'Kitchen denied',
  KITCHEN_PREPARING: 'Preparing',
  COURIER_REFUNDED: 'Courier refunded',
  DELIVERY_PENDING: 'Ready for courier',
  DELIVERY_PICKING: 'Picking up',
  DELIVERY_DENIED: 'Delivery denied',
  DELIVERY_REFUNDED: 'Delivery refunded',
  DELIVERY_DELIVERING: 'En route',
  ORDER_COMPLETED: 'Delivered',
  ORDER_FAILED: 'Failed',
}

export const customerFlow = [
  'CUSTOMER_CREATED',
  'CUSTOMER_PAID',
  'KITCHEN_PREPARING',
  'DELIVERY_PICKING',
  'DELIVERY_DELIVERING',
  'ORDER_COMPLETED',
]

export const courierFlow = ['DELIVERY_PENDING', 'DELIVERY_PICKING', 'DELIVERY_DELIVERING', 'ORDER_COMPLETED']

export const kitchenFlow = ['CUSTOMER_PAID', 'KITCHEN_ACCEPTED', 'KITCHEN_PREPARING', 'DELIVERY_PENDING']

export function formatStatus(status) {
  return statusLabels[status] || status || 'Unknown'
}

export function isTerminal(status) {
  return ['ORDER_COMPLETED', 'ORDER_FAILED', 'CUSTOMER_CANCELLED', 'KITCHEN_DENIED', 'DELIVERY_DENIED'].includes(status)
}
