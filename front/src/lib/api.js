const bases = {
  customer: import.meta.env.VITE_CUSTOMER_API_URL || 'http://localhost:8091',
  courier: import.meta.env.VITE_COURIER_API_URL || 'http://localhost:8090',
  restaurant: import.meta.env.VITE_RESTAURANT_API_URL || 'http://localhost:8092',
}

function appendQuery(path, params = {}) {
  const entries = Object.entries(params).filter(([, value]) => value !== undefined && value !== null && value !== '')
  if (!entries.length) return path
  const query = new URLSearchParams(entries).toString()
  return `${path}${path.includes('?') ? '&' : '?'}${query}`
}

async function request(service, path, options = {}) {
  const response = await fetch(`${bases[service]}${path}`, {
    method: options.method || 'GET',
    headers: {
      Accept: 'application/json',
      ...(options.body ? { 'Content-Type': 'application/json' } : {}),
      ...options.headers,
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  })

  const text = await response.text()
  let data = null

  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
  }

  if (!response.ok) {
    const message =
      data?.error_message ||
      data?.errorMessage ||
      data?.error ||
      data?.message ||
      `${response.status} ${response.statusText}`
    throw new Error(message)
  }

  return data
}

function field(item, ...names) {
  for (const name of names) {
    if (item?.[name] !== undefined && item?.[name] !== null) return item[name]
  }
  return undefined
}

export function normalizeOrder(order) {
  return {
    id: field(order, 'id', 'ID', 'Id'),
    customerId: field(order, 'customer_id', 'customerId', 'CustomerID'),
    courierId: field(order, 'courier_id', 'courierId', 'CourierID'),
    createdAt: field(order, 'created_at', 'createdAt', 'CreatedAt'),
    updatedAt: field(order, 'updated_at', 'updatedAt', 'UpdatedAt'),
    status: field(order, 'status', 'Status') || 'CUSTOMER_CREATED',
  }
}

export function normalizeRestaurant(item) {
  return {
    id: field(item, 'id', 'ID', 'Id'),
    name: field(item, 'name', 'Name') || 'Restaurant',
    address: field(item, 'address', 'Address') || '',
    status: field(item, 'status', 'Status') ?? true,
  }
}

export function normalizeCourier(item) {
  return {
    id: field(item, 'id', 'ID', 'Id'),
    name: field(item, 'name', 'Name') || 'Courier',
    transportType: field(item, 'transport_type', 'transportType', 'TransportType') || 'courier',
    isActive: field(item, 'is_active', 'isActive', 'IsActive') ?? true,
  }
}

export function normalizeMenuItem(item) {
  return {
    id: field(item, 'id', 'ID', 'Id', 'order_item_id', 'OrderItemID'),
    restaurantId: field(item, 'restaurant_id', 'restaurantId', 'RestaurantID'),
    name: field(item, 'name', 'Name') || 'Menu item',
    price: Number(field(item, 'price', 'Price') || 0),
    quantity: Number(field(item, 'quantity', 'Quantity') || 0),
    description: field(item, 'description', 'Description') || '',
  }
}

function normalizeList(data, normalizer) {
  return Array.isArray(data) ? data.map(normalizer) : []
}

function authService(role) {
  if (!['customer', 'courier', 'restaurant'].includes(role)) {
    throw new Error(`Unknown role: ${role}`)
  }
  return role
}

export const api = {
  bases,
  auth: {
    login(role, payload) {
      return request(authService(role), '/login', {
        method: 'POST',
        body: {
          wallet_address: payload.walletAddress,
          password: payload.password,
        },
      })
    },
    register(role, payload) {
      const body = {
        name: payload.name,
        wallet_address: payload.walletAddress,
        password: payload.password,
      }

      if (role === 'courier') {
        body.transport_type = payload.transportType || 'bicycle'
      } else {
        body.address = payload.address
        body.is_active = true
      }

      return request(authService(role), '/register', {
        method: 'POST',
        body,
      })
    },
  },
  customer: {
    restaurants() {
      return request('customer', '/restaurants').then((data) => normalizeList(data, normalizeRestaurant))
    },
    couriers() {
      return request('customer', '/couriers').then((data) => normalizeList(data, normalizeCourier))
    },
    menu(restaurantId) {
      return request('customer', appendQuery('/menu', { restaurant_id: restaurantId })).then((data) =>
        normalizeList(data, normalizeMenuItem),
      )
    },
    orders(filters = {}) {
      return request(
        'customer',
        appendQuery('/orders', {
          customer_id: filters.customerId,
          courier_id: filters.courierId,
          status: filters.status,
        }),
      ).then((data) => normalizeList(data, normalizeOrder))
    },
    order(orderId) {
      return request('customer', `/orders/${orderId}`).then(normalizeOrder)
    },
    status(orderId) {
      return request('customer', `/orders/${orderId}/status`).then((data) => data?.status || data?.Status || '')
    },
    updateStatus(orderId, status) {
      return request('customer', `/orders/${orderId}/status`, {
        method: 'PUT',
        body: { status },
      })
    },
    createOrder(customerId, courierId) {
      return request('customer', '/orders', {
        method: 'POST',
        body: {
          customer_id: customerId,
          courier_id: courierId,
        },
      }).then(normalizeOrder)
    },
    addItem(orderId, item) {
      return request('customer', `/orders/${orderId}/items`, {
        method: 'POST',
        body: {
          restaurant_item_id: item.id,
          price: item.price,
          quantity: item.quantity,
        },
      })
    },
    pay(orderId) {
      return request('customer', `/orders/${orderId}/pay`, { method: 'POST' })
    },
    total(orderId) {
      return request('customer', `/orders/${orderId}/total`).then((data) => Number(data?.total || data?.Total || 0))
    },
  },
  restaurant: {
    menu(restaurantId) {
      return request('restaurant', appendQuery('/menu/show', { restaurant_id: restaurantId })).then((data) =>
        normalizeList(data, normalizeMenuItem),
      )
    },
    uploadMenuItem(item) {
      return request('restaurant', '/menu/upload', {
        method: 'POST',
        body: {
          restaurant_id: item.restaurantId,
          name: item.name,
          price: Number(item.price),
          quantity: Number(item.quantity),
          description: item.description,
        },
      })
    },
    orders(restaurantId, status = '') {
      return request('restaurant', appendQuery('/orders', { restaurant_id: restaurantId, status })).then((data) =>
        normalizeList(data, normalizeOrder),
      )
    },
  },
}
