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
    return 'All orders'
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
            <span className="pill pill-soft">{orders.length} total</span>
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
      </div>
    </div>
  )
}
