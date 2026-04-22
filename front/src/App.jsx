import { AnimatePresence, motion } from 'framer-motion'
import {
  ArrowRight,
  Bike,
  ChefHat,
  CheckCircle2,
  CircleDollarSign,
  ClipboardList,
  Clock3,
  History,
  Loader2,
  LogOut,
  MapPin,
  Menu as MenuIcon,
  PackageCheck,
  PanelLeft,
  Plus,
  RefreshCw,
  Route as RouteIcon,
  Search,
  ShieldCheck,
  ShoppingBag,
  Store,
  Truck,
  User,
  Utensils,
  Wallet,
  X,
} from 'lucide-react'
import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  Link,
  Navigate,
  NavLink,
  Route,
  Routes,
  useLocation,
  useNavigate,
  useParams,
} from 'react-router-dom'
import { api } from './lib/api.js'
import { clearSession, readSession, writeSession } from './lib/session.js'
import { courierFlow, customerFlow, formatStatus, isTerminal, kitchenFlow, orderStatuses } from './lib/status.js'

const roles = {
  customer: {
    label: 'Customer',
    defaultPath: 'explore',
    icon: ShoppingBag,
    tone: 'tomato',
    nav: 'market',
    routes: [
      { path: 'explore', label: 'Explore', icon: Search },
      { path: 'orders', label: 'Orders', icon: ClipboardList },
      { path: 'profile', label: 'Profile', icon: User },
    ],
  },
  courier: {
    label: 'Courier',
    defaultPath: 'queue',
    icon: Bike,
    tone: 'basil',
    nav: 'rail',
    routes: [
      { path: 'queue', label: 'Queue', icon: PackageCheck },
      { path: 'status', label: 'Status', icon: RouteIcon },
      { path: 'history', label: 'History', icon: History },
      { path: 'profile', label: 'Profile', icon: User },
    ],
  },
  restaurant: {
    label: 'Restaurant',
    defaultPath: 'menu',
    icon: Store,
    tone: 'saffron',
    nav: 'ops',
    routes: [
      { path: 'menu', label: 'Menu', icon: Utensils },
      { path: 'orders', label: 'Orders', icon: ClipboardList },
      { path: 'kitchen', label: 'Kitchen', icon: ChefHat },
      { path: 'profile', label: 'Profile', icon: User },
    ],
  },
}

const roleKeys = Object.keys(roles)

function roleEntryPath(role) {
  const config = roles[role]
  if (!config) return '/'
  return readSession(role)?.id ? `/${role}/${config.defaultPath}` : `/${role}/auth`
}

function formatMoney(value) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 2,
  }).format(Number(value || 0))
}

function formatDate(value) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '—'
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

function shortId(value) {
  if (!value) return '—'
  return `${String(value).slice(0, 8)}…${String(value).slice(-4)}`
}

function useQuery(loader, dependencies, initialData = []) {
  const [state, setState] = useState({ data: initialData, loading: true, error: '' })

  const reload = useCallback(async () => {
    setState((current) => ({ ...current, loading: true, error: '' }))
    try {
      const data = await loader()
      setState({ data, loading: false, error: '' })
      return data
    } catch (error) {
      setState({ data: initialData, loading: false, error: error.message || 'Request failed' })
      return null
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, dependencies)

  useEffect(() => {
    reload()
  }, [reload])

  return { ...state, reload }
}

function AnimatedPage({ children, className = '' }) {
  return (
    <motion.div
      className={`page-motion ${className}`}
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -10 }}
      transition={{ duration: 0.22, ease: 'easeOut' }}
    >
      {children}
    </motion.div>
  )
}

function Button({ children, variant = 'primary', busy = false, className = '', ...props }) {
  return (
    <motion.button
      className={`button button-${variant} ${className}`}
      whileTap={{ scale: props.disabled ? 1 : 0.98 }}
      whileHover={props.disabled ? undefined : { y: -1 }}
      disabled={busy || props.disabled}
      {...props}
    >
      {busy ? <Loader2 className="spin" size={16} /> : null}
      {children}
    </motion.button>
  )
}

function Alert({ tone = 'danger', children }) {
  if (!children) return null
  return <div className={`alert alert-${tone}`}>{children}</div>
}

function EmptyState({ icon: Icon = ClipboardList, title, action }) {
  return (
    <div className="empty-state">
      <Icon size={26} />
      <strong>{title}</strong>
      {action}
    </div>
  )
}

function SkeletonRows({ rows = 4 }) {
  return (
    <div className="skeleton-stack">
      {Array.from({ length: rows }).map((_, index) => (
        <div className="skeleton-row" key={index} />
      ))}
    </div>
  )
}

function StatusPill({ status }) {
  return <span className={`status-pill status-${status}`}>{formatStatus(status)}</span>
}

function PageHeader({ eyebrow, title, icon: Icon, action, children }) {
  return (
    <header className="page-header">
      <div className="header-title">
        {Icon ? (
          <span className="header-icon">
            <Icon size={20} />
          </span>
        ) : null}
        <div>
          {eyebrow ? <span className="eyebrow">{eyebrow}</span> : null}
          <h1>{title}</h1>
        </div>
      </div>
      {children ? <div className="header-meta">{children}</div> : null}
      {action ? <div className="header-action">{action}</div> : null}
    </header>
  )
}

function Field({ label, children }) {
  return (
    <label className="field">
      <span>{label}</span>
      {children}
    </label>
  )
}

function StatusStepper({ status, flow }) {
  const activeIndex = Math.max(0, flow.indexOf(status))
  return (
    <div className="status-stepper">
      {flow.map((item, index) => {
        const done = index <= activeIndex
        return (
          <div className={`status-step ${done ? 'is-done' : ''}`} key={item}>
            <motion.span
              className="status-dot"
              animate={{ scale: done ? 1 : 0.82, opacity: done ? 1 : 0.45 }}
              transition={{ duration: 0.18 }}
            >
              {done ? <CheckCircle2 size={14} /> : <Clock3 size={14} />}
            </motion.span>
            <span>{formatStatus(item)}</span>
          </div>
        )
      })}
    </div>
  )
}

function OrderTable({ orders, emptyTitle, detailBase, onRefresh, loading }) {
  if (loading) return <SkeletonRows rows={5} />
  if (!orders.length) {
    return (
      <EmptyState
        icon={ClipboardList}
        title={emptyTitle}
        action={
          onRefresh ? (
            <Button variant="ghost" onClick={onRefresh}>
              <RefreshCw size={15} />
              Refresh
            </Button>
          ) : null
        }
      />
    )
  }

  return (
    <div className="table-shell">
      <table>
        <thead>
          <tr>
            <th>Order</th>
            <th>Customer</th>
            <th>Courier</th>
            <th>Status</th>
            <th>Updated</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>{shortId(order.id)}</td>
              <td>{shortId(order.customerId)}</td>
              <td>{shortId(order.courierId)}</td>
              <td>
                <StatusPill status={order.status} />
              </td>
              <td>{formatDate(order.updatedAt || order.createdAt)}</td>
              <td className="table-actions">
                {detailBase ? (
                  <Link className="icon-link" to={`${detailBase}/${order.id}`} aria-label="Open order">
                    <ArrowRight size={16} />
                  </Link>
                ) : null}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function PortalSelector() {
  return (
    <AnimatedPage className="portal-page">
      <div className="portal-brand">
        <div className="brand-mark">
          <Utensils size={26} />
        </div>
        <div>
          <span className="eyebrow">YAFDS</span>
          <h1>Food delivery control room</h1>
        </div>
      </div>

      <div className="role-grid">
        {roleKeys.map((role) => {
          const config = roles[role]
          const Icon = config.icon
          const entryPath = roleEntryPath(role)
          return (
            <motion.div className={`role-card role-card-${config.tone}`} key={role} whileHover={{ y: -4 }}>
              <div className="role-card-head">
                <span className="role-icon">
                  <Icon size={24} />
                </span>
                <div>
                  <span className="eyebrow">{config.label}</span>
                  <h2>{config.label} portal</h2>
                </div>
              </div>
              <div className="role-actions">
                <Link className="button button-primary" to={`/${role}/auth`}>
                  Sign in
                  <ArrowRight size={16} />
                </Link>
                <Link className="button button-ghost" to={entryPath}>
                  {entryPath.endsWith('/auth') ? 'Login required' : 'Open'}
                </Link>
              </div>
            </motion.div>
          )
        })}
      </div>
    </AnimatedPage>
  )
}

function AuthPage() {
  const { role } = useParams()
  const navigate = useNavigate()
  const config = roles[role]
  const [mode, setMode] = useState('login')
  const [form, setForm] = useState({
    name: '',
    walletAddress: '',
    address: '',
    transportType: 'bicycle',
    password: '',
  })
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  if (!config) return <Navigate to="/" replace />

  const Icon = config.icon

  function updateForm(event) {
    const { name, value } = event.target
    setForm((current) => ({ ...current, [name]: value }))
  }

  async function submit(event) {
    event.preventDefault()
    setBusy(true)
    setError('')

    try {
      if (mode === 'register') {
        await api.auth.register(role, form)
      }

      const loginResponse = await api.auth.login(role, form)
      writeSession(role, loginResponse)
      navigate(`/${role}/${config.defaultPath}`, { replace: true })
    } catch (requestError) {
      setError(requestError.message || 'Authentication failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <AnimatedPage className={`auth-page tone-${config.tone}`}>
      <Link className="back-link" to="/">
        <ArrowRight size={15} />
        Portals
      </Link>
      <div className="auth-layout">
        <section className="auth-panel">
          <div className="role-card-head">
            <span className="role-icon">
              <Icon size={24} />
            </span>
            <div>
              <span className="eyebrow">{config.label}</span>
              <h1>{mode === 'login' ? 'Sign in' : 'Create account'}</h1>
            </div>
          </div>

          <div className="segmented">
            <button className={mode === 'login' ? 'active' : ''} onClick={() => setMode('login')} type="button">
              Login
            </button>
            <button className={mode === 'register' ? 'active' : ''} onClick={() => setMode('register')} type="button">
              Register
            </button>
          </div>

          <form className="form-stack" onSubmit={submit}>
            {mode === 'register' ? (
              <Field label="Name">
                <input name="name" value={form.name} onChange={updateForm} required />
              </Field>
            ) : null}
            <Field label="Wallet address">
              <input name="walletAddress" value={form.walletAddress} onChange={updateForm} required />
            </Field>
            {mode === 'register' && role !== 'courier' ? (
              <Field label="Address">
                <input name="address" value={form.address} onChange={updateForm} required />
              </Field>
            ) : null}
            {mode === 'register' && role === 'courier' ? (
              <Field label="Transport">
                <select name="transportType" value={form.transportType} onChange={updateForm}>
                  <option value="bicycle">Bicycle</option>
                  <option value="scooter">Scooter</option>
                  <option value="car">Car</option>
                  <option value="walk">Walk</option>
                </select>
              </Field>
            ) : null}
            <Field label="Password">
              <input name="password" type="password" value={form.password} onChange={updateForm} required />
            </Field>
            <Alert>{error}</Alert>
            <Button busy={busy} type="submit">
              {mode === 'login' ? 'Sign in' : 'Register'}
            </Button>
          </form>
        </section>

        <section className="auth-status">
          <div className="endpoint-list">
            <span className="eyebrow">API</span>
            <strong>{api.bases[role]}</strong>
          </div>
          <div className="auth-meter">
            <span />
            <span />
            <span />
          </div>
        </section>
      </div>
    </AnimatedPage>
  )
}

function RoleApp() {
  const { role } = useParams()
  const config = roles[role]
  const [session, setSession] = useState(() => readSession(role))

  useEffect(() => {
    setSession(readSession(role))
  }, [role])

  if (!config) return <Navigate to="/" replace />
  if (!session?.id) return <Navigate to={`/${role}/auth`} replace />

  return <AppShell role={role} config={config} session={session} setSession={setSession} />
}

function AppShell({ role, config, session, setSession }) {
  const navigate = useNavigate()
  const location = useLocation()
  const [navOpen, setNavOpen] = useState(false)

  function logout() {
    clearSession(role)
    setSession(null)
    navigate(`/${role}/auth`, { replace: true })
  }

  const nav = (
    <nav className={`app-nav app-nav-${config.nav}`} aria-label={`${config.label} navigation`}>
      {config.routes.map((item) => {
        const Icon = item.icon
        return (
          <NavLink
            className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            key={item.path}
            to={`/${role}/${item.path}`}
            onClick={() => setNavOpen(false)}
          >
            <Icon size={18} />
            <span>{item.label}</span>
          </NavLink>
        )
      })}
    </nav>
  )

  return (
    <div className={`app-shell app-shell-${role} tone-${config.tone}`}>
      {config.nav !== 'market' ? (
        <aside className={`sidebar sidebar-${config.nav}`}>
          <Link className="shell-brand" to="/">
            <Utensils size={20} />
            <span>YAFDS</span>
          </Link>
          {nav}
        </aside>
      ) : null}

      <main className="workspace">
        <header className="topbar">
          <button className="mobile-menu" type="button" onClick={() => setNavOpen(true)} aria-label="Open menu">
            <MenuIcon size={19} />
          </button>
          <div className="topbar-left">
            <span className="workspace-role">{config.label}</span>
            <strong>{session.name || shortId(session.id)}</strong>
          </div>
          <div className="role-switcher">
            {roleKeys.map((item) => (
              <Link className={item === role ? 'active' : ''} key={item} to={roleEntryPath(item)}>
                {roles[item].label}
              </Link>
            ))}
          </div>
          <Button variant="ghost" onClick={logout}>
            <LogOut size={15} />
            Logout
          </Button>
        </header>

        {config.nav === 'market' ? nav : null}

        <AnimatePresence>
          {navOpen ? (
            <motion.div className="drawer-backdrop" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
              <motion.aside
                className="drawer"
                initial={{ x: -280 }}
                animate={{ x: 0 }}
                exit={{ x: -280 }}
                transition={{ duration: 0.2 }}
              >
                <button className="drawer-close" type="button" onClick={() => setNavOpen(false)} aria-label="Close menu">
                  <X size={18} />
                </button>
                {nav}
              </motion.aside>
            </motion.div>
          ) : null}
        </AnimatePresence>

        <section className="screen-frame">
          <AnimatePresence mode="wait">
            <Routes location={location} key={location.pathname}>
              <Route index element={<Navigate to={config.defaultPath} replace />} />
              {role === 'customer' ? (
                <>
                  <Route path="explore" element={<CustomerExplore session={session} />} />
                  <Route path="orders" element={<CustomerOrders session={session} />} />
                  <Route path="orders/:orderId" element={<CustomerOrderDetail session={session} />} />
                  <Route path="profile" element={<ProfileScreen role={role} session={session} />} />
                </>
              ) : null}
              {role === 'courier' ? (
                <>
                  <Route path="queue" element={<CourierQueue session={session} />} />
                  <Route path="status" element={<CourierStatus session={session} />} />
                  <Route path="history" element={<CourierHistory session={session} />} />
                  <Route path="profile" element={<ProfileScreen role={role} session={session} />} />
                </>
              ) : null}
              {role === 'restaurant' ? (
                <>
                  <Route path="menu" element={<RestaurantMenu session={session} />} />
                  <Route path="orders" element={<RestaurantOrders session={session} />} />
                  <Route path="kitchen" element={<KitchenBoard session={session} />} />
                  <Route path="profile" element={<ProfileScreen role={role} session={session} />} />
                </>
              ) : null}
              <Route path="*" element={<Navigate to={config.defaultPath} replace />} />
            </Routes>
          </AnimatePresence>
        </section>
      </main>
    </div>
  )
}

function CustomerExplore({ session }) {
  const navigate = useNavigate()
  const restaurantsQuery = useQuery(() => api.customer.restaurants(), [], [])
  const couriersQuery = useQuery(() => api.customer.couriers(), [], [])
  const [restaurantId, setRestaurantId] = useState('')
  const [courierId, setCourierId] = useState('')
  const [cart, setCart] = useState([])
  const [creating, setCreating] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (!restaurantId && restaurantsQuery.data.length) setRestaurantId(restaurantsQuery.data[0].id)
  }, [restaurantId, restaurantsQuery.data])

  useEffect(() => {
    if (!courierId && couriersQuery.data.length) setCourierId(couriersQuery.data[0].id)
  }, [courierId, couriersQuery.data])

  const menuQuery = useQuery(
    () => (restaurantId ? api.customer.menu(restaurantId) : Promise.resolve([])),
    [restaurantId],
    [],
  )

  const selectedRestaurant = restaurantsQuery.data.find((item) => item.id === restaurantId)
  const selectedCourier = couriersQuery.data.find((item) => item.id === courierId)
  const cartTotal = cart.reduce((sum, item) => sum + item.price * item.quantity, 0)

  function addToCart(item) {
    setCart((current) => {
      const existing = current.find((cartItem) => cartItem.id === item.id)
      if (existing) {
        return current.map((cartItem) =>
          cartItem.id === item.id ? { ...cartItem, quantity: cartItem.quantity + 1 } : cartItem,
        )
      }
      return [...current, { ...item, quantity: 1 }]
    })
  }

  function setQuantity(itemId, quantity) {
    setCart((current) =>
      current
        .map((item) => (item.id === itemId ? { ...item, quantity: Math.max(0, quantity) } : item))
        .filter((item) => item.quantity > 0),
    )
  }

  async function createOrder() {
    setCreating(true)
    setError('')
    setMessage('')

    try {
      const order = await api.customer.createOrder(session.id, courierId)
      for (const item of cart) {
        await api.customer.addItem(order.id, item)
      }
      setCart([])
      setMessage(`Created order ${shortId(order.id)}`)
      navigate(`/customer/orders/${order.id}`)
    } catch (requestError) {
      setError(requestError.message || 'Order creation failed')
    } finally {
      setCreating(false)
    }
  }

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Customer"
        title="Explore"
        icon={ShoppingBag}
        action={
          <Button variant="ghost" onClick={() => restaurantsQuery.reload()}>
            <RefreshCw size={15} />
            Refresh
          </Button>
        }
      />

      <div className="customer-grid">
        <section className="panel catalog-panel">
          <div className="panel-head">
            <h2>Restaurants</h2>
            <span>{restaurantsQuery.data.length}</span>
          </div>
          <Alert>{restaurantsQuery.error}</Alert>
          {restaurantsQuery.loading ? <SkeletonRows rows={4} /> : null}
          {!restaurantsQuery.loading && !restaurantsQuery.data.length ? (
            <EmptyState icon={Store} title="No restaurants returned" />
          ) : null}
          <div className="restaurant-list">
            {restaurantsQuery.data.map((restaurant) => (
              <button
                className={`selection-row ${restaurant.id === restaurantId ? 'active' : ''}`}
                key={restaurant.id}
                onClick={() => {
                  setRestaurantId(restaurant.id)
                  setCart([])
                }}
                type="button"
              >
                <Store size={17} />
                <span>
                  <strong>{restaurant.name}</strong>
                  <small>{restaurant.address || shortId(restaurant.id)}</small>
                </span>
                <StatusPill status={restaurant.status ? 'KITCHEN_ACCEPTED' : 'KITCHEN_DENIED'} />
              </button>
            ))}
          </div>
        </section>

        <section className="panel menu-panel">
          <div className="panel-head">
            <h2>{selectedRestaurant?.name || 'Menu'}</h2>
            <span>{menuQuery.data.length}</span>
          </div>
          <Alert>{menuQuery.error}</Alert>
          {menuQuery.loading ? <SkeletonRows rows={5} /> : null}
          {!menuQuery.loading && !menuQuery.data.length ? <EmptyState icon={Utensils} title="No menu items returned" /> : null}
          <div className="menu-grid">
            {menuQuery.data.map((item) => (
              <motion.article className="menu-item" key={item.id} whileHover={{ y: -2 }}>
                <div>
                  <h3>{item.name}</h3>
                  <p>{item.description || shortId(item.id)}</p>
                </div>
                <div className="menu-item-foot">
                  <strong>{formatMoney(item.price)}</strong>
                  <Button variant="secondary" onClick={() => addToCart(item)}>
                    <Plus size={15} />
                    Add
                  </Button>
                </div>
              </motion.article>
            ))}
          </div>
        </section>

        <aside className="panel order-draft">
          <div className="panel-head">
            <h2>Order draft</h2>
            <strong>{formatMoney(cartTotal)}</strong>
          </div>
          <Field label="Courier">
            <select value={courierId} onChange={(event) => setCourierId(event.target.value)}>
              {couriersQuery.data.map((courier) => (
                <option key={courier.id} value={courier.id}>
                  {courier.name} · {courier.transportType}
                </option>
              ))}
            </select>
          </Field>
          <Alert>{couriersQuery.error}</Alert>
          <div className="cart-lines">
            {!cart.length ? <EmptyState icon={ShoppingBag} title="Draft is empty" /> : null}
            {cart.map((item) => (
              <div className="cart-line" key={item.id}>
                <span>
                  <strong>{item.name}</strong>
                  <small>{formatMoney(item.price)}</small>
                </span>
                <div className="quantity-control">
                  <button type="button" onClick={() => setQuantity(item.id, item.quantity - 1)}>
                    -
                  </button>
                  <span>{item.quantity}</span>
                  <button type="button" onClick={() => setQuantity(item.id, item.quantity + 1)}>
                    +
                  </button>
                </div>
              </div>
            ))}
          </div>
          {selectedCourier ? (
            <div className="chosen-courier">
              <Truck size={16} />
              {selectedCourier.name}
            </div>
          ) : null}
          <Alert>{error}</Alert>
          <Alert tone="success">{message}</Alert>
          <Button busy={creating} disabled={!cart.length || !courierId} onClick={createOrder}>
            Create order
          </Button>
        </aside>
      </div>
    </AnimatedPage>
  )
}

function CustomerOrders({ session }) {
  const [status, setStatus] = useState('')
  const ordersQuery = useQuery(() => api.customer.orders({ customerId: session.id, status }), [session.id, status], [])

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Customer"
        title="Past orders"
        icon={ClipboardList}
        action={
          <div className="toolbar">
            <select value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="">All statuses</option>
              {orderStatuses.map((item) => (
                <option key={item} value={item}>
                  {formatStatus(item)}
                </option>
              ))}
            </select>
            <Button variant="ghost" onClick={ordersQuery.reload}>
              <RefreshCw size={15} />
              Refresh
            </Button>
          </div>
        }
      />
      <Alert>{ordersQuery.error}</Alert>
      <OrderTable
        detailBase="/customer/orders"
        emptyTitle="No customer orders returned"
        loading={ordersQuery.loading}
        onRefresh={ordersQuery.reload}
        orders={ordersQuery.data}
      />
    </AnimatedPage>
  )
}

function CustomerOrderDetail() {
  const { orderId } = useParams()
  const [busyAction, setBusyAction] = useState('')
  const detailQuery = useQuery(
    async () => {
      const [order, status, total] = await Promise.all([
        api.customer.order(orderId),
        api.customer.status(orderId),
        api.customer.total(orderId),
      ])
      return { order: { ...order, status: status || order.status }, total }
    },
    [orderId],
    null,
  )

  const detail = detailQuery.data

  async function payOrder() {
    setBusyAction('pay')
    try {
      await api.customer.pay(orderId)
      await detailQuery.reload()
    } finally {
      setBusyAction('')
    }
  }

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Order"
        title={shortId(orderId)}
        icon={ClipboardList}
        action={
          <Button variant="ghost" onClick={detailQuery.reload}>
            <RefreshCw size={15} />
            Refresh
          </Button>
        }
      />
      <Alert>{detailQuery.error}</Alert>
      {detailQuery.loading ? <SkeletonRows rows={4} /> : null}
      {detail ? (
        <div className="detail-grid">
          <section className="panel">
            <div className="panel-head">
              <h2>Status</h2>
              <StatusPill status={detail.order.status} />
            </div>
            <StatusStepper status={detail.order.status} flow={customerFlow} />
          </section>
          <section className="panel metrics-panel">
            <div className="metric">
              <CircleDollarSign size={20} />
              <span>Total</span>
              <strong>{formatMoney(detail.total)}</strong>
            </div>
            <div className="metric">
              <Clock3 size={20} />
              <span>Updated</span>
              <strong>{formatDate(detail.order.updatedAt || detail.order.createdAt)}</strong>
            </div>
          </section>
          <section className="panel">
            <div className="profile-list">
              <InfoRow icon={User} label="Customer" value={shortId(detail.order.customerId)} />
              <InfoRow icon={Truck} label="Courier" value={shortId(detail.order.courierId)} />
              <InfoRow icon={ClipboardList} label="Order" value={detail.order.id} />
            </div>
            <Button busy={busyAction === 'pay'} disabled={detail.order.status !== 'CUSTOMER_CREATED'} onClick={payOrder}>
              Pay order
            </Button>
          </section>
        </div>
      ) : null}
    </AnimatedPage>
  )
}

function CourierQueue({ session }) {
  const ordersQuery = useQuery(() => api.customer.orders({ courierId: session.id }), [session.id], [])
  const activeOrders = useMemo(() => ordersQuery.data.filter((order) => !isTerminal(order.status)), [ordersQuery.data])

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Courier"
        title="Queue"
        icon={PackageCheck}
        action={
          <Button variant="ghost" onClick={ordersQuery.reload}>
            <RefreshCw size={15} />
            Refresh
          </Button>
        }
      />
      <StatsRow
        items={[
          { label: 'Assigned', value: ordersQuery.data.length, icon: ClipboardList },
          { label: 'Active', value: activeOrders.length, icon: Truck },
          { label: 'Done', value: ordersQuery.data.filter((item) => item.status === 'ORDER_COMPLETED').length, icon: CheckCircle2 },
        ]}
      />
      <Alert>{ordersQuery.error}</Alert>
      <OrderTable
        emptyTitle="No assigned orders returned"
        loading={ordersQuery.loading}
        onRefresh={ordersQuery.reload}
        orders={activeOrders}
      />
    </AnimatedPage>
  )
}

function CourierStatus({ session }) {
  const ordersQuery = useQuery(() => api.customer.orders({ courierId: session.id }), [session.id], [])
  const activeOrders = ordersQuery.data.filter((order) => !isTerminal(order.status))
  const [updating, setUpdating] = useState('')

  async function setStatus(orderId, status) {
    setUpdating(`${orderId}:${status}`)
    try {
      await api.customer.updateStatus(orderId, status)
      await ordersQuery.reload()
    } finally {
      setUpdating('')
    }
  }

  return (
    <AnimatedPage>
      <PageHeader eyebrow="Courier" title="Status updates" icon={RouteIcon} />
      <Alert>{ordersQuery.error}</Alert>
      {ordersQuery.loading ? <SkeletonRows rows={4} /> : null}
      {!ordersQuery.loading && !activeOrders.length ? <EmptyState icon={Truck} title="No active deliveries returned" /> : null}
      <div className="work-grid">
        {activeOrders.map((order) => (
          <article className="panel order-work-card" key={order.id}>
            <div className="panel-head">
              <h2>{shortId(order.id)}</h2>
              <StatusPill status={order.status} />
            </div>
            <StatusStepper status={order.status} flow={courierFlow} />
            <div className="action-row">
              <Button
                busy={updating === `${order.id}:DELIVERY_PICKING`}
                onClick={() => setStatus(order.id, 'DELIVERY_PICKING')}
                variant="secondary"
              >
                Pick up
              </Button>
              <Button
                busy={updating === `${order.id}:DELIVERY_DELIVERING`}
                onClick={() => setStatus(order.id, 'DELIVERY_DELIVERING')}
                variant="secondary"
              >
                En route
              </Button>
              <Button
                busy={updating === `${order.id}:ORDER_COMPLETED`}
                onClick={() => setStatus(order.id, 'ORDER_COMPLETED')}
              >
                Delivered
              </Button>
            </div>
          </article>
        ))}
      </div>
    </AnimatedPage>
  )
}

function CourierHistory({ session }) {
  const ordersQuery = useQuery(() => api.customer.orders({ courierId: session.id }), [session.id], [])
  const history = ordersQuery.data.filter((order) => isTerminal(order.status))

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Courier"
        title="History"
        icon={History}
        action={
          <Button variant="ghost" onClick={ordersQuery.reload}>
            <RefreshCw size={15} />
            Refresh
          </Button>
        }
      />
      <Alert>{ordersQuery.error}</Alert>
      <OrderTable
        emptyTitle="No completed deliveries returned"
        loading={ordersQuery.loading}
        onRefresh={ordersQuery.reload}
        orders={history}
      />
    </AnimatedPage>
  )
}

function RestaurantMenu({ session }) {
  const menuQuery = useQuery(() => api.restaurant.menu(session.id), [session.id], [])
  const [form, setForm] = useState({ name: '', price: '', quantity: '1', description: '' })
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [message, setMessage] = useState('')

  function updateForm(event) {
    const { name, value } = event.target
    setForm((current) => ({ ...current, [name]: value }))
  }

  async function submit(event) {
    event.preventDefault()
    setBusy(true)
    setError('')
    setMessage('')

    try {
      await api.restaurant.uploadMenuItem({
        restaurantId: session.id,
        name: form.name,
        price: form.price,
        quantity: form.quantity,
        description: form.description,
      })
      setForm({ name: '', price: '', quantity: '1', description: '' })
      setMessage('Menu item uploaded')
      await menuQuery.reload()
    } catch (requestError) {
      setError(requestError.message || 'Upload failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Restaurant"
        title="Menu management"
        icon={Utensils}
        action={
          <Button variant="ghost" onClick={menuQuery.reload}>
            <RefreshCw size={15} />
            Refresh
          </Button>
        }
      />
      <div className="restaurant-menu-layout">
        <section className="panel">
          <div className="panel-head">
            <h2>Menu items</h2>
            <span>{menuQuery.data.length}</span>
          </div>
          <Alert>{menuQuery.error}</Alert>
          {menuQuery.loading ? <SkeletonRows rows={5} /> : null}
          {!menuQuery.loading && !menuQuery.data.length ? <EmptyState icon={Utensils} title="No menu items returned" /> : null}
          <div className="menu-list-table">
            {menuQuery.data.map((item) => (
              <div className="menu-row" key={item.id}>
                <span>
                  <strong>{item.name}</strong>
                  <small>{item.description || shortId(item.id)}</small>
                </span>
                <strong>{formatMoney(item.price)}</strong>
                <span>{item.quantity}</span>
              </div>
            ))}
          </div>
        </section>

        <aside className="panel">
          <div className="panel-head">
            <h2>Upload item</h2>
            <Plus size={18} />
          </div>
          <form className="form-stack" onSubmit={submit}>
            <Field label="Name">
              <input name="name" value={form.name} onChange={updateForm} required />
            </Field>
            <div className="form-grid-two">
              <Field label="Price">
                <input min="0.01" name="price" step="0.01" type="number" value={form.price} onChange={updateForm} required />
              </Field>
              <Field label="Quantity">
                <input min="1" name="quantity" step="1" type="number" value={form.quantity} onChange={updateForm} required />
              </Field>
            </div>
            <Field label="Description">
              <textarea name="description" rows="4" value={form.description} onChange={updateForm} />
            </Field>
            <Alert>{error}</Alert>
            <Alert tone="success">{message}</Alert>
            <Button busy={busy} type="submit">
              Upload
            </Button>
          </form>
        </aside>
      </div>
    </AnimatedPage>
  )
}

function RestaurantOrders({ session }) {
  const [status, setStatus] = useState('')
  const ordersQuery = useQuery(() => api.restaurant.orders(session.id, status), [session.id, status], [])

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Restaurant"
        title="Orders"
        icon={ClipboardList}
        action={
          <div className="toolbar">
            <select value={status} onChange={(event) => setStatus(event.target.value)}>
              <option value="">All statuses</option>
              {orderStatuses.map((item) => (
                <option key={item} value={item}>
                  {formatStatus(item)}
                </option>
              ))}
            </select>
            <Button variant="ghost" onClick={ordersQuery.reload}>
              <RefreshCw size={15} />
              Refresh
            </Button>
          </div>
        }
      />
      <Alert>{ordersQuery.error}</Alert>
      <OrderTable
        emptyTitle="No restaurant orders returned"
        loading={ordersQuery.loading}
        onRefresh={ordersQuery.reload}
        orders={ordersQuery.data}
      />
    </AnimatedPage>
  )
}

function KitchenBoard({ session }) {
  const ordersQuery = useQuery(() => api.restaurant.orders(session.id), [session.id], [])
  const [updating, setUpdating] = useState('')
  const activeOrders = ordersQuery.data.filter((order) => !['ORDER_COMPLETED', 'ORDER_FAILED'].includes(order.status))

  async function setStatus(orderId, status) {
    setUpdating(`${orderId}:${status}`)
    try {
      await api.customer.updateStatus(orderId, status)
      await ordersQuery.reload()
    } finally {
      setUpdating('')
    }
  }

  return (
    <AnimatedPage>
      <PageHeader
        eyebrow="Restaurant"
        title="Kitchen board"
        icon={ChefHat}
        action={
          <Button variant="ghost" onClick={ordersQuery.reload}>
            <RefreshCw size={15} />
            Refresh
          </Button>
        }
      />
      <Alert>{ordersQuery.error}</Alert>
      {ordersQuery.loading ? <SkeletonRows rows={4} /> : null}
      {!ordersQuery.loading && !activeOrders.length ? <EmptyState icon={ChefHat} title="No kitchen orders returned" /> : null}
      <div className="kitchen-grid">
        {activeOrders.map((order) => (
          <article className="panel order-work-card" key={order.id}>
            <div className="panel-head">
              <h2>{shortId(order.id)}</h2>
              <StatusPill status={order.status} />
            </div>
            <StatusStepper status={order.status} flow={kitchenFlow} />
            <div className="action-row">
              <Button
                busy={updating === `${order.id}:KITCHEN_ACCEPTED`}
                onClick={() => setStatus(order.id, 'KITCHEN_ACCEPTED')}
                variant="secondary"
              >
                Accept
              </Button>
              <Button
                busy={updating === `${order.id}:KITCHEN_PREPARING`}
                onClick={() => setStatus(order.id, 'KITCHEN_PREPARING')}
                variant="secondary"
              >
                Preparing
              </Button>
              <Button
                busy={updating === `${order.id}:DELIVERY_PENDING`}
                onClick={() => setStatus(order.id, 'DELIVERY_PENDING')}
              >
                Ready
              </Button>
              <Button
                busy={updating === `${order.id}:KITCHEN_DENIED`}
                onClick={() => setStatus(order.id, 'KITCHEN_DENIED')}
                variant="danger"
              >
                Deny
              </Button>
            </div>
          </article>
        ))}
      </div>
    </AnimatedPage>
  )
}

function InfoRow({ icon: Icon, label, value }) {
  return (
    <div className="info-row">
      <Icon size={17} />
      <span>{label}</span>
      <strong>{value || '—'}</strong>
    </div>
  )
}

function StatsRow({ items }) {
  return (
    <div className="stats-row">
      {items.map((item) => {
        const Icon = item.icon
        return (
          <div className="stat-box" key={item.label}>
            <Icon size={18} />
            <span>{item.label}</span>
            <strong>{item.value}</strong>
          </div>
        )
      })}
    </div>
  )
}

function ProfileScreen({ role, session }) {
  const config = roles[role]
  const Icon = config.icon

  return (
    <AnimatedPage>
      <PageHeader eyebrow={config.label} title="Profile" icon={User} />
      <div className="profile-grid">
        <section className="panel profile-card">
          <span className="profile-avatar">
            <Icon size={32} />
          </span>
          <h2>{session.name || config.label}</h2>
          <StatusPill status={session.isActive ? 'KITCHEN_ACCEPTED' : 'KITCHEN_DENIED'} />
        </section>
        <section className="panel">
          <div className="profile-list">
            <InfoRow icon={ShieldCheck} label="Role" value={config.label} />
            <InfoRow icon={Wallet} label="Wallet" value={session.walletAddress} />
            <InfoRow icon={User} label="Account ID" value={session.id} />
            {session.address ? <InfoRow icon={MapPin} label="Address" value={session.address} /> : null}
            {session.transportType ? <InfoRow icon={Bike} label="Transport" value={session.transportType} /> : null}
            <InfoRow icon={Clock3} label="Session until" value={session.expiration ? formatDate(session.expiration * 1000) : '—'} />
          </div>
        </section>
      </div>
    </AnimatedPage>
  )
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<PortalSelector />} />
      <Route path="/:role/auth" element={<AuthPage />} />
      <Route path="/:role/*" element={<RoleApp />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
