import { useEffect, useMemo, useState } from 'react'
import { useParams, useNavigate, useLocation } from 'react-router-dom'

export default function Dashboard() {
  const { role } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [sessionMessage, setSessionMessage] = useState('')
  const [timeLeft, setTimeLeft] = useState(null)
  const [orders, setOrders] = useState([])
  const [ordersLoading, setOrdersLoading] = useState(false)
  const [ordersError, setOrdersError] = useState('')
  const [statusFilter, setStatusFilter] = useState('')
  const [createOrderModal, setCreateOrderModal] = useState(false)
  const [selectedCourier, setSelectedCourier] = useState('')
  const [couriers, setCouriers] = useState([])
  const [couriersLoading, setCouriersLoading] = useState(false)
  const [couriersError, setCouriersError] = useState('')
  const [restaurants, setRestaurants] = useState([])
  const [restaurantsLoading, setRestaurantsLoading] = useState(false)
  const [restaurantsError, setRestaurantsError] = useState('')
  const [creatingOrder, setCreatingOrder] = useState(false)
  const [createOrderError, setCreateOrderError] = useState('')
  const [selectedRestaurant, setSelectedRestaurant] = useState('')
  const [restaurantMenu, setRestaurantMenu] = useState([])
  const [restaurantMenuLoading, setRestaurantMenuLoading] = useState(false)
  const [restaurantMenuError, setRestaurantMenuError] = useState('')
  const [orderItems, setOrderItems] = useState({})
  const [menuItems, setMenuItems] = useState([])
  const [menuItemsLoading, setMenuItemsLoading] = useState(false)
  const [menuItemsError, setMenuItemsError] = useState('')
  const [menuForm, setMenuForm] = useState({
    name: '',
    price: '',
    quantity: '',
    description: ''
  })
  const [menuSaving, setMenuSaving] = useState(false)
  const [menuSaveError, setMenuSaveError] = useState('')
  const [menuSaveSuccess, setMenuSaveSuccess] = useState('')
  const [addItemModal, setAddItemModal] = useState(false)
  const [addItemOrder, setAddItemOrder] = useState(null)
  const [addItemRestaurantId, setAddItemRestaurantId] = useState('')
  const [addItemMenu, setAddItemMenu] = useState([])
  const [addItemMenuLoading, setAddItemMenuLoading] = useState(false)
  const [addItemMenuError, setAddItemMenuError] = useState('')
  const [addItemQuantities, setAddItemQuantities] = useState({})
  const [addItemSaving, setAddItemSaving] = useState(false)
  const [addItemError, setAddItemError] = useState('')
  const [addItemSuccess, setAddItemSuccess] = useState('')

  const apiBaseByRole = useMemo(
    () => ({
      customer: import.meta.env.VITE_CUSTOMER_API_URL || 'http://localhost:8091',
      courier: import.meta.env.VITE_COURIER_API_URL || 'http://localhost:8090',
      restaurant: import.meta.env.VITE_RESTAURANT_API_URL || 'http://localhost:8092',
    }),
    [],
  )

  const apiBase = apiBaseByRole[role] || apiBaseByRole.customer

  useEffect(() => {
    hydrateUser()
  }, [])

  const hydrateUser = () => {
    const stateUser = location.state?.user
    const stored = localStorage.getItem('currentUser')
    const savedUser = stateUser || (stored ? JSON.parse(stored) : null)

    if (!savedUser) {
      navigate(`/${role}/auth`, { state: { role } })
      return
    }

    if (savedUser.role !== role) {
      navigate('/')
      return
    }

    const expired = savedUser.expiration && savedUser.expiration * 1000 <= Date.now()
    if (expired) {
      handleExpiredSession()
      return
    }

    setUser(savedUser)
    setLoading(false)
  }

  useEffect(() => {
    const userId = user?.id || user?.Id
    if (!userId) return

    const controller = new AbortController()
    const fetchOrders = async () => {
      setOrdersLoading(true)
      setOrdersError('')

      const params = new URLSearchParams()
      if (role === 'customer') params.set('customer_id', userId)
      if (role === 'courier') params.set('courier_id', userId)
      if (role === 'restaurant') params.set('restaurant_id', userId)
      if (statusFilter) params.set('status', statusFilter)

      const query = params.toString()
      const endpoint = query ? `${apiBase}/orders?${query}` : `${apiBase}/orders`

      try {
        const response = await fetch(endpoint, { signal: controller.signal })
        const data = await response.json()

        if (!response.ok) {
          throw new Error(data?.error || 'Не удалось получить заказы')
        }

        setOrders(Array.isArray(data) ? data : [])
      } catch (error) {
        if (error.name === 'AbortError') return
        setOrdersError(error.message)
      } finally {
        setOrdersLoading(false)
      }
    }

    fetchOrders()

    return () => controller.abort()
  }, [apiBase, role, user, statusFilter])

  useEffect(() => {
    if (role !== 'customer') return

    const controller = new AbortController()
    const fetchCouriers = async () => {
      setCouriersLoading(true)
      setCouriersError('')
      try {
        const response = await fetch(`${apiBase}/couriers`, { signal: controller.signal })
        const data = await response.json()

        if (!response.ok) {
          throw new Error(data?.error || 'Не удалось получить курьеров')
        }

        setCouriers(Array.isArray(data) ? data : [])
      } catch (error) {
        if (error.name === 'AbortError') return
        setCouriersError(error.message)
        setCouriers([])
      } finally {
        setCouriersLoading(false)
      }
    }

    fetchCouriers()
    return () => controller.abort()
  }, [apiBase, role])

  useEffect(() => {
    if (role !== 'customer') return

    const controller = new AbortController()
    const fetchRestaurants = async () => {
      setRestaurantsLoading(true)
      setRestaurantsError('')
      try {
        const response = await fetch(`${apiBase}/restaurants`, { signal: controller.signal })
        const data = await response.json()

        if (!response.ok) {
          throw new Error(data?.error || 'Не удалось получить рестораны')
        }

        setRestaurants(Array.isArray(data) ? data : [])
      } catch (error) {
        if (error.name === 'AbortError') return
        setRestaurantsError(error.message)
        setRestaurants([])
      } finally {
        setRestaurantsLoading(false)
      }
    }

    fetchRestaurants()
    return () => controller.abort()
  }, [apiBase, role])

  const fetchMenuItems = async (signal) => {
    if (role !== 'restaurant') return
    const restaurantId = user?.id || user?.Id
    if (!restaurantId) return

    setMenuItemsLoading(true)
    setMenuItemsError('')

    try {
      const response = await fetch(
        `${apiBase}/menu/show?restaurant_id=${restaurantId}`,
        { signal }
      )
      const data = await response.json()

      if (!response.ok) {
        throw new Error(data?.error || 'Не удалось получить меню')
      }

      setMenuItems(Array.isArray(data) ? data : [])
    } catch (error) {
      if (error.name === 'AbortError') return
      setMenuItemsError(error.message)
      setMenuItems([])
    } finally {
      setMenuItemsLoading(false)
    }
  }

  useEffect(() => {
    if (role !== 'restaurant') return

    const controller = new AbortController()
    fetchMenuItems(controller.signal)
    return () => controller.abort()
  }, [apiBase, role, user])

  const handleExpiredSession = () => {
    setLoading(false)
    setSessionMessage('Session expired, please sign in again.')
    handleSignOut({ expired: true })
  }

  const handleSignOut = ({ expired = false } = {}) => {
    localStorage.removeItem('currentUser')
    navigate(`/${role}/auth`, { state: { role, expired } })
  }

  const handleCreateOrder = async () => {
    if (!selectedCourier) {
      setCreateOrderError('Please select a courier')
      return
    }
    if (!selectedRestaurant) {
      setCreateOrderError('Please provide a restaurant id')
      return
    }

    const itemsPayload = restaurantMenu
      .map((item) => {
        const id = item.order_item_id || item.orderItemID || item.id
        const quantity = Number(orderItems[id] || 0)
        return { restaurant_item_id: id, quantity }
      })
      .filter((item) => item.quantity > 0)

    if (itemsPayload.length === 0) {
      setCreateOrderError('Please select at least one menu item')
      return
    }
    setCreatingOrder(true)
    setCreateOrderError('')
    try {
      const response = await fetch(`${apiBase}/orders`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          customer_id: user.id || user.Id,
          courier_id: selectedCourier,
          restaurant_id: selectedRestaurant,
          status: 'created',
          items: itemsPayload
        })
      })
      const data = await response.json()
      if (!response.ok) {
        throw new Error(data.error || 'Failed to create order')
      }
      setCreateOrderModal(false)
      setSelectedCourier('')
      setSelectedRestaurant('')
      setRestaurantMenu([])
      setOrderItems({})
      // Refresh orders
      setStatusFilter('')
      // Orders will refresh automatically due to useEffect
    } catch (error) {
      setCreateOrderError(error.message)
    } finally {
      setCreatingOrder(false)
    }
  }

  const openAddItemModal = (order) => {
    setAddItemOrder(order)
    setAddItemRestaurantId('')
    setAddItemMenu([])
    setAddItemQuantities({})
    setAddItemMenuError('')
    setAddItemError('')
    setAddItemSuccess('')
    setAddItemModal(true)
  }

  const fetchMenuForAddItem = async () => {
    if (!addItemRestaurantId) {
      setAddItemMenuError('Restaurant id is required')
      return
    }
    setAddItemMenuLoading(true)
    setAddItemMenuError('')
    try {
      const response = await fetch(
        `${apiBase}/menu?restaurant_id=${addItemRestaurantId}`
      )
      const data = await response.json()

      if (!response.ok) {
        throw new Error(data?.error || 'Не удалось получить меню ресторана')
      }

      setAddItemMenu(Array.isArray(data) ? data : [])
      setAddItemQuantities({})
    } catch (error) {
      setAddItemMenuError(error.message)
      setAddItemMenu([])
    } finally {
      setAddItemMenuLoading(false)
    }
  }

  const handleAddItem = async () => {
    if (!addItemOrder?.id) {
      setAddItemError('Order id is missing')
      return
    }
    if (!addItemRestaurantId) {
      setAddItemError('Restaurant id is required')
      return
    }

    const selected = addItemMenu
      .map((item) => {
        const id = item.order_item_id || item.orderItemID || item.id
        const quantity = Number(addItemQuantities[id] || 0)
        return { restaurant_item_id: id, quantity }
      })
      .filter((item) => item.quantity > 0)

    if (selected.length === 0) {
      setAddItemError('Please select one menu item')
      return
    }
    if (selected.length > 1) {
      setAddItemError('Please select only one item at a time')
      return
    }

    setAddItemSaving(true)
    setAddItemError('')
    setAddItemSuccess('')

    try {
      const response = await fetch(`${apiBase}/orders/${addItemOrder.id}/items`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          restaurant_id: addItemRestaurantId,
          restaurant_item_id: selected[0].restaurant_item_id,
          quantity: selected[0].quantity
        })
      })
      const data = await response.json()
      if (!response.ok) {
        throw new Error(data?.error || 'Failed to add item')
      }

      setAddItemSuccess('Item added to order')
      setAddItemMenu([])
      setAddItemQuantities({})
    } catch (error) {
      setAddItemError(error.message)
    } finally {
      setAddItemSaving(false)
    }
  }

  const fetchRestaurantMenuForOrder = async () => {
    if (!selectedRestaurant) {
      setRestaurantMenuError('Restaurant id is required')
      return
    }
    setRestaurantMenuLoading(true)
    setRestaurantMenuError('')
    try {
      const response = await fetch(
        `${apiBase}/menu?restaurant_id=${selectedRestaurant.trim()}`
      )
      const data = await response.json()

      if (!response.ok) {
        throw new Error(data?.error || 'Не удалось получить меню ресторана')
      }

      const menu = Array.isArray(data) ? data : []
      if (menu.length === 0) {
        setRestaurantMenuError('В этом ресторане пока нет доступных блюд')
      }
      setRestaurantMenu(menu)
      setOrderItems({})
    } catch (error) {
      setRestaurantMenuError(error.message)
      setRestaurantMenu([])
    } finally {
      setRestaurantMenuLoading(false)
    }
  }

  useEffect(() => {
    if (!user?.expiration) return

    const updateCountdown = () => {
      const remaining = user.expiration * 1000 - Date.now()
      setTimeLeft(Math.max(remaining, 0))
      if (remaining <= 0) {
        handleExpiredSession()
      }
    }

    updateCountdown()
    const timer = setInterval(updateCountdown, 1000)
    return () => clearInterval(timer)
  }, [user])

  const countdownLabel = useMemo(() => {
    if (timeLeft === null) return ''
    const totalSeconds = Math.floor(timeLeft / 1000)
    const minutes = Math.floor(totalSeconds / 60)
    const seconds = totalSeconds % 60
    return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  }, [timeLeft])

  const orderScopeLabel = useMemo(() => {
    if (role === 'customer') return 'Your orders'
    if (role === 'courier') return 'Your deliveries'
    return 'Restaurant orders'
  }, [role])

  const formatDate = (value) => {
    if (!value) return '—'
    const parsed = new Date(value)
    return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString()
  }

  const formatStatus = (value) => {
    if (!value) return '—'
    const normalized = value.toString().trim()
    return normalized ? normalized.replace(/^./, (c) => c.toUpperCase()) : '—'
  }

  const formatPrice = (value) => {
    const num = Number(value)
    if (!Number.isFinite(num)) return '—'
    return `$${num.toFixed(2)}`
  }

  const handleMenuUpload = async (event) => {
    event.preventDefault()
    setMenuSaveError('')
    setMenuSaveSuccess('')

    const restaurantId = user?.id || user?.Id
    if (!restaurantId) {
      setMenuSaveError('Restaurant ID is missing')
      return
    }

    if (!menuForm.name.trim()) {
      setMenuSaveError('Name is required')
      return
    }

    if (!menuForm.description.trim()) {
      setMenuSaveError('Description is required')
      return
    }

    const priceValue = Number(menuForm.price)
    if (!Number.isFinite(priceValue) || priceValue <= 0) {
      setMenuSaveError('Price must be greater than 0')
      return
    }

    const quantityValue = Number(menuForm.quantity || 0)
    if (!Number.isInteger(quantityValue) || quantityValue < 0) {
      setMenuSaveError('Quantity must be 0 or more')
      return
    }

    setMenuSaving(true)
    try {
      const response = await fetch(`${apiBase}/menu/upload`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          restaurant_id: restaurantId,
          name: menuForm.name.trim(),
          price: priceValue,
          quantity: quantityValue,
          description: menuForm.description.trim()
        })
      })
      const data = await response.json()
      if (!response.ok) {
        throw new Error(data?.error || data?.error_message || 'Failed to upload menu item')
      }
      setMenuSaveSuccess('Menu item uploaded')
      setMenuForm({ name: '', price: '', quantity: '', description: '' })
      await fetchMenuItems()
    } catch (error) {
      setMenuSaveError(error.message)
    } finally {
      setMenuSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="dashboard-shell">
        <div className="dashboard-loading">Loading...</div>
      </div>
    )
  }

  return (
    <div className="dashboard-shell">
      <div className="dashboard-header">
        <div>
          <p className="eyebrow">{role.toUpperCase()} ACCOUNT</p>
          <h1 className="headline">Hello, {user?.name || 'guest'}</h1>
        </div>
        <button onClick={handleSignOut} className="dashboard-ghost">
          Sign out
        </button>
      </div>

      {sessionMessage && <div className="dashboard-alert">{sessionMessage}</div>}

      <div className="dashboard-grid">
        <section className="dashboard-card">
          <header className="card-header">
            <div>
              <p className="eyebrow">Profile</p>
              <h2 className="card-title">Personal information</h2>
            </div>
          </header>
          <dl className="profile-list">
            {role === 'restaurant' && (
              <div>
                <dt>Restaurant ID</dt>
                <dd className="mono">{user?.id || user?.Id || '—'}</dd>
              </div>
            )}
            <div>
              <dt>Name</dt>
              <dd>{user?.name || '—'}</dd>
            </div>
            <div>
              <dt>Wallet</dt>
              <dd className="mono">{user?.wallet_address || '—'}</dd>
            </div>
            <div>
              <dt>Delivery address</dt>
              <dd>{user?.address || '—'}</dd>
            </div>
            {role === 'restaurant' && (
              <div>
                <dt>Active</dt>
                <dd className="pill">{user?.is_active ? 'Active' : 'Inactive'}</dd>
              </div>
            )}
            <div>
              <dt>Role</dt>
              <dd className="pill">{role}</dd>
            </div>
          </dl>
        </section>

        <section className="dashboard-card">
          <header className="card-header">
            <div>
              <p className="eyebrow">Session</p>
              <h2 className="card-title">Security</h2>
            </div>
          </header>
          <div className="session-row">
            <div>
              <p className="label">Token expires in</p>
              <p className="countdown">{countdownLabel || '—'}</p>
            </div>
            <span className="dot" aria-hidden />
            <div>
              <p className="label">Expiration time</p>
              <p className="value">{user?.expiration ? new Date(user.expiration * 1000).toLocaleString() : '—'}</p>
            </div>
          </div>
          <div className="session-note">You will be signed out automatically when the token expires.</div>
        </section>

        <section className="dashboard-card orders-card">
          <header className="card-header">
            <div>
              <p className="eyebrow">Orders</p>
              <h2 className="card-title">{orderScopeLabel}</h2>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
              {role === 'customer' && (
                <button
                  onClick={() => setCreateOrderModal(true)}
                  className="dashboard-ghost"
                  style={{ background: 'var(--accent)', color: 'white', border: 'none' }}
                >
                  Create Order
                </button>
              )}
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="filter-select"
                style={{
                  padding: '0.4rem 0.8rem',
                  borderRadius: '6px',
                  border: '1px solid var(--border)',
                  background: 'var(--card-bg)',
                  color: 'var(--text-main)',
                  fontSize: '0.875rem',
                  cursor: 'pointer',
                }}
              >
                <option value="">All Statuses</option>
                <option value="created">Created</option>
                <option value="pending">Pending</option>
                <option value="delivering">Delivering</option>
                <option value="delivered">Delivered</option>
                <option value="cancelled">Cancelled</option>
              </select>
              <span className="pill pill-soft">{orders.length} total</span>
            </div>
          </header>

          {ordersLoading && <div className="orders-state">Загружаем заказы...</div>}
          {!ordersLoading && ordersError && <div className="dashboard-alert">{ordersError}</div>}
          {!ordersLoading && !ordersError && orders.length === 0 && (
            <div className="orders-state">Пока нет заказов</div>
          )}

          {!ordersLoading && !ordersError && orders.length > 0 && (
            <div className="orders-list">
              {orders.map((order) => (
                <div key={order.id} className="order-row">
                  <div>
                    <p className="label">Order</p>
                    <p className="order-id">#{String(order.id).slice(0, 8)}</p>
                    <p className="order-hint">Customer: {order.customer_id || '—'}</p>
                    <p className="order-hint">Courier: {order.courier_id || '—'}</p>
                  </div>
                  <div className="order-status">
                    <p className="label">Status</p>
                    <span className="pill pill-ghost">{formatStatus(order.status)}</span>
                    {role === 'customer' && (
                      <button
                        onClick={() => openAddItemModal(order)}
                        className="dashboard-ghost"
                        style={{ marginTop: '0.75rem' }}
                      >
                        Add item
                      </button>
                    )}
                  </div>
                  <div className="order-dates">
                    <p className="label">Created</p>
                    <p className="value">{formatDate(order.created_at)}</p>
                    <p className="order-hint">Updated: {formatDate(order.updated_at)}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        {role === 'restaurant' && (
          <section className="dashboard-card menu-card">
            <header className="card-header">
              <div>
                <p className="eyebrow">Menu</p>
                <h2 className="card-title">Your menu items</h2>
              </div>
              <button
                onClick={() => fetchMenuItems()}
                className="dashboard-ghost"
                disabled={menuItemsLoading}
              >
                Refresh
              </button>
            </header>

            {menuItemsLoading && <div className="orders-state">Загружаем меню...</div>}
            {!menuItemsLoading && menuItemsError && (
              <div className="dashboard-alert">{menuItemsError}</div>
            )}
            {!menuItemsLoading && !menuItemsError && menuItems.length === 0 && (
              <div className="orders-state">There are no items on the menu yet</div>
            )}

            {!menuItemsLoading && !menuItemsError && menuItems.length > 0 && (
              <div className="menu-list">
                {menuItems.map((item) => (
                  <div
                    key={item.order_item_id || item.orderItemID || item.id}
                    className="menu-row"
                  >
                    <div>
                      <p className="label">Название</p>
                      <p className="menu-name">{item.name || '—'}</p>
                      <p className="menu-desc">{item.description || '—'}</p>
                    </div>
                    <div className="menu-meta">
                      <p className="label">Цена</p>
                      <p className="menu-value">{formatPrice(item.price)}</p>
                      <p className="menu-hint">Количество: {item.quantity ?? '—'}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <div className="menu-divider" />

            <form className="menu-form" onSubmit={handleMenuUpload}>
              <h3 className="menu-form-title">Add new menu item</h3>
              <div className="menu-form-grid">
                <input
                  type="text"
                  value={menuForm.name}
                  onChange={(e) => setMenuForm((prev) => ({ ...prev, name: e.target.value }))}
                  placeholder="Item name"
                  className="menu-input"
                  required
                />
                <input
                  type="number"
                  min="0"
                  step="0.01"
                  value={menuForm.price}
                  onChange={(e) => setMenuForm((prev) => ({ ...prev, price: e.target.value }))}
                  placeholder="Price"
                  className="menu-input"
                  required
                />
                <input
                  type="number"
                  min="0"
                  step="1"
                  value={menuForm.quantity}
                  onChange={(e) => setMenuForm((prev) => ({ ...prev, quantity: e.target.value }))}
                  placeholder="Quantity"
                  className="menu-input"
                  required
                />
              </div>
              <textarea
                value={menuForm.description}
                onChange={(e) => setMenuForm((prev) => ({ ...prev, description: e.target.value }))}
                placeholder="Description"
                className="menu-textarea"
                rows={3}
                required
              />
              {menuSaveError && <div className="dashboard-alert">{menuSaveError}</div>}
              {menuSaveSuccess && <div className="menu-success">{menuSaveSuccess}</div>}
              <button type="submit" className="menu-submit" disabled={menuSaving}>
                {menuSaving ? 'Uploading...' : 'Upload menu item'}
              </button>
            </form>
          </section>
        )}

        {createOrderModal && (
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              background: 'rgba(0,0,0,0.5)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              zIndex: 1000
            }}
            onClick={() => setCreateOrderModal(false)}
          >
            <div
              style={{
                background: 'var(--card-bg)',
                padding: '2rem',
                borderRadius: '8px',
                maxWidth: '400px',
                width: '90%',
                boxShadow: '0 4px 12px rgba(0,0,0,0.15)'
              }}
              onClick={(e) => e.stopPropagation()}
            >
              <h3 style={{ marginTop: 0 }}>Create New Order</h3>
              <div style={{ marginBottom: '1rem' }}>
                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Select Courier:</label>
                <select
                  value={selectedCourier}
                  onChange={(e) => setSelectedCourier(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '0.5rem',
                    borderRadius: '4px',
                    border: '1px solid var(--border)',
                    background: 'var(--bg)',
                    color: 'var(--text-main)'
                  }}
                  disabled={couriersLoading || couriers.length === 0}
                >
                  <option value="">Choose a courier...</option>
                  {couriersLoading && <option disabled>Loading...</option>}
                  {!couriersLoading && couriersError && <option disabled>{couriersError}</option>}
                  {!couriersLoading && !couriersError && couriers.length === 0 && <option disabled>No couriers found</option>}
                  {!couriersLoading && !couriersError && couriers.map((courier) => (
                    <option key={courier.id} value={courier.id}>
                      {courier.name || courier.wallet_address || courier.id}
                    </option>
                  ))}
                </select>
              </div>
              <div style={{ marginBottom: '1rem' }}>
                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Select Restaurant:</label>
                <select
                  value={selectedRestaurant}
                  onChange={(e) => setSelectedRestaurant(e.target.value)}
                  style={{
                    width: '100%',
                    padding: '0.5rem',
                    borderRadius: '4px',
                    border: '1px solid var(--border)',
                    background: 'var(--bg)',
                    color: 'var(--text-main)'
                  }}
                  disabled={restaurantsLoading || restaurants.length === 0}
                >
                  <option value="">Choose a restaurant...</option>
                  {restaurantsLoading && <option disabled>Loading...</option>}
                  {!restaurantsLoading && restaurantsError && <option disabled>{restaurantsError}</option>}
                  {!restaurantsLoading && !restaurantsError && restaurants.length === 0 && <option disabled>No restaurants found</option>}
                  {!restaurantsLoading && !restaurantsError && restaurants.map((restaurant) => (
                    <option key={restaurant.id} value={restaurant.id}>
                      {restaurant.name || restaurant.id}
                    </option>
                  ))}
                </select>
                <button
                  onClick={fetchRestaurantMenuForOrder}
                  className="dashboard-ghost"
                  disabled={restaurantMenuLoading || !selectedRestaurant}
                  style={{ marginTop: '0.5rem' }}
                >
                  {restaurantMenuLoading ? 'Loading menu...' : 'Load menu'}
                </button>
                {restaurantMenuError && (
                  <div style={{ color: 'red', marginTop: '0.5rem' }}>{restaurantMenuError}</div>
                )}
              </div>
              {restaurantMenu.length > 0 && (
                <div style={{ marginBottom: '1rem' }}>
                  <p style={{ marginBottom: '0.5rem' }}>Menu items</p>
                  <div style={{ display: 'grid', gap: '0.75rem' }}>
                    {restaurantMenu.map((item) => {
                      const id = item.order_item_id || item.orderItemID || item.id
                      return (
                        <div
                          key={id}
                          style={{
                            border: '1px solid var(--border)',
                            borderRadius: '6px',
                            padding: '0.75rem'
                          }}
                        >
                          <div style={{ display: 'flex', justifyContent: 'space-between', gap: '1rem' }}>
                            <div>
                              <div style={{ fontWeight: 600 }}>{item.name || '—'}</div>
                              <div style={{ fontSize: '0.85rem', opacity: 0.8 }}>{item.description || '—'}</div>
                            </div>
                            <div style={{ textAlign: 'right' }}>
                              <div style={{ fontWeight: 600 }}>{formatPrice(item.price)}</div>
                              <input
                                type="number"
                                min="0"
                                step="1"
                                value={orderItems[id] || ''}
                                onChange={(e) =>
                                  setOrderItems((prev) => ({
                                    ...prev,
                                    [id]: e.target.value
                                  }))
                                }
                                placeholder="Qty"
                                style={{
                                  width: '72px',
                                  marginTop: '0.4rem',
                                  padding: '0.35rem',
                                  borderRadius: '4px',
                                  border: '1px solid var(--border)',
                                  background: 'var(--bg)',
                                  color: 'var(--text-main)'
                                }}
                              />
                            </div>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                </div>
              )}
              {createOrderError && (
                <div style={{ color: 'red', marginBottom: '1rem' }}>{createOrderError}</div>
              )}
              <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
                <button
                  onClick={() => setCreateOrderModal(false)}
                  className="dashboard-ghost"
                  disabled={creatingOrder}
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateOrder}
                  disabled={creatingOrder || !selectedCourier}
                  style={{
                    background: 'var(--accent)',
                    color: 'white',
                    border: 'none',
                    padding: '0.5rem 1rem',
                    borderRadius: '4px',
                    cursor: creatingOrder ? 'not-allowed' : 'pointer'
                  }}
                >
                  {creatingOrder ? 'Creating...' : 'Create Order'}
                </button>
              </div>
            </div>
          </div>
        )}

        {addItemModal && (
          <div
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              background: 'rgba(0,0,0,0.5)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              zIndex: 1000
            }}
            onClick={() => setAddItemModal(false)}
          >
            <div
              style={{
                background: 'var(--card-bg)',
                padding: '2rem',
                borderRadius: '8px',
                maxWidth: '420px',
                width: '90%',
                boxShadow: '0 4px 12px rgba(0,0,0,0.15)'
              }}
              onClick={(e) => e.stopPropagation()}
            >
              <h3 style={{ marginTop: 0 }}>Add item to order</h3>
              <p style={{ marginBottom: '1rem', opacity: 0.75 }}>
                Order #{String(addItemOrder?.id || '').slice(0, 8)}
              </p>
              <div style={{ marginBottom: '1rem' }}>
                <label style={{ display: 'block', marginBottom: '0.5rem' }}>Restaurant id:</label>
                <input
                  type="text"
                  value={addItemRestaurantId}
                  onChange={(e) => setAddItemRestaurantId(e.target.value)}
                  placeholder="Enter restaurant UUID"
                  style={{
                    width: '100%',
                    padding: '0.5rem',
                    borderRadius: '4px',
                    border: '1px solid var(--border)',
                    background: 'var(--bg)',
                    color: 'var(--text-main)'
                  }}
                />
                <button
                  onClick={fetchMenuForAddItem}
                  className="dashboard-ghost"
                  disabled={addItemMenuLoading || !addItemRestaurantId}
                  style={{ marginTop: '0.5rem' }}
                >
                  {addItemMenuLoading ? 'Loading menu...' : 'Load menu'}
                </button>
                {addItemMenuError && (
                  <div style={{ color: 'red', marginTop: '0.5rem' }}>{addItemMenuError}</div>
                )}
              </div>
              {addItemMenu.length > 0 && (
                <div style={{ marginBottom: '1rem' }}>
                  <p style={{ marginBottom: '0.5rem' }}>Menu items</p>
                  <div style={{ display: 'grid', gap: '0.75rem' }}>
                    {addItemMenu.map((item) => {
                      const id = item.order_item_id || item.orderItemID || item.id
                      return (
                        <div
                          key={id}
                          style={{
                            border: '1px solid var(--border)',
                            borderRadius: '6px',
                            padding: '0.75rem'
                          }}
                        >
                          <div style={{ display: 'flex', justifyContent: 'space-between', gap: '1rem' }}>
                            <div>
                              <div style={{ fontWeight: 600 }}>{item.name || '—'}</div>
                              <div style={{ fontSize: '0.85rem', opacity: 0.8 }}>{item.description || '—'}</div>
                            </div>
                            <div style={{ textAlign: 'right' }}>
                              <div style={{ fontWeight: 600 }}>{formatPrice(item.price)}</div>
                              <input
                                type="number"
                                min="0"
                                step="1"
                                value={addItemQuantities[id] || ''}
                                onChange={(e) =>
                                  setAddItemQuantities((prev) => ({
                                    ...prev,
                                    [id]: e.target.value
                                  }))
                                }
                                placeholder="Qty"
                                style={{
                                  width: '72px',
                                  marginTop: '0.4rem',
                                  padding: '0.35rem',
                                  borderRadius: '4px',
                                  border: '1px solid var(--border)',
                                  background: 'var(--bg)',
                                  color: 'var(--text-main)'
                                }}
                              />
                            </div>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                </div>
              )}
              {addItemError && (
                <div style={{ color: 'red', marginBottom: '0.75rem' }}>{addItemError}</div>
              )}
              {addItemSuccess && (
                <div style={{ color: 'green', marginBottom: '0.75rem' }}>{addItemSuccess}</div>
              )}
              <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
                <button
                  onClick={() => setAddItemModal(false)}
                  className="dashboard-ghost"
                  disabled={addItemSaving}
                >
                  Cancel
                </button>
                <button
                  onClick={handleAddItem}
                  disabled={addItemSaving || !addItemRestaurantId}
                  style={{
                    background: 'var(--accent)',
                    color: 'white',
                    border: 'none',
                    padding: '0.5rem 1rem',
                    borderRadius: '4px',
                    cursor: addItemSaving ? 'not-allowed' : 'pointer'
                  }}
                >
                  {addItemSaving ? 'Adding...' : 'Add item'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
