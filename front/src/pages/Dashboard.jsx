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
      </div>
    </div>
  )
}
