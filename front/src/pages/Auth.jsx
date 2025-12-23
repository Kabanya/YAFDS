import { useEffect, useMemo, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

export default function Auth() {
  const location = useLocation()
  const navigate = useNavigate()
  const role = useMemo(() => {
    const [, pathRole] = location.pathname.split('/')
    const normalized = (pathRole || '').toLowerCase()
    return ['customer', 'courier', 'restaurant'].includes(normalized) ? normalized : 'customer'
  }, [location.pathname])

  const apiBaseByRole = useMemo(
    () => ({
      customer: import.meta.env.VITE_CUSTOMER_API_URL || 'http://localhost:8091',
      courier: import.meta.env.VITE_COURIER_API_URL || 'http://localhost:8090',
      restaurant: import.meta.env.VITE_RESTAURANT_API_URL || 'http://localhost:8092',
    }),
    [],
  )

  const apiBase = apiBaseByRole[role] || apiBaseByRole.customer

  const [isSignUp, setIsSignUp] = useState(false)
  const [name, setName] = useState('')
  const [walletAddress, setWalletAddress] = useState('')
  const [address, setAddress] = useState('')
  const [password, setPassword] = useState('')
  const [transportType, setTransportType] = useState('bicycle')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (location.state?.expired) {
      setError('Session expired, please sign in again.')
    }
  }, [location.state])

  const persistSession = (profile) => {
    localStorage.setItem('currentUser', JSON.stringify(profile))
    navigate(`/${role}/dashboard`, { state: { user: profile } })
  }

  const authenticate = async (wallet, pass) => {
    const response = await fetch(`${apiBase}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        wallet_address: wallet,
        password: pass,
      }),
    })

    const data = await response.json()
    if (!response.ok) throw new Error(data?.error_message || 'Invalid wallet address or password')

    const userId = data.id || data.Id

    return {
      id: userId,
      name: data.name || wallet,
      wallet_address: data.wallet_address || wallet,
      address: data.address || wallet,
      transport_type: data.transport_type,
      status: data.status,
      token: data.token,
      expiration: Number(data.expiration) || 0,
      role,
    }
  }

  const handleAuth = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      if (isSignUp) {
        const body = {
          name,
          wallet_address: walletAddress,
          password,
        }

        if (role === 'customer' || role === 'restaurant') {
          body.address = address
        }

        if (role === 'restaurant') {
          body.status = true
        }

        if (role === 'courier') {
          body.transport_type = transportType
        }

        const response = await fetch(`${apiBase}/register`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(body),
        })

        const data = await response.json()
        if (!response.ok) {
          throw new Error(data?.error_message || 'Failed to register')
        }

        const profile = await authenticate(walletAddress, password)
        persistSession(profile)
      } else {
        const profile = await authenticate(walletAddress, password)
        persistSession(profile)
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-container">
      <div className="auth-box">
        <h1 className="auth-title">{role.toUpperCase()}</h1>

        <form onSubmit={handleAuth} className="auth-form">
          <input
            type="text"
            placeholder="Wallet Address"
            value={walletAddress}
            onChange={(e) => setWalletAddress(e.target.value)}
            required
            className="auth-input"
          />

          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            className="auth-input"
          />

          {isSignUp && (
            <>
              <input
                type="text"
                placeholder={role === 'restaurant' ? 'Restaurant Name' : 'Name (optional)'}
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="auth-input"
              />
              {(role === 'customer' || role === 'restaurant') && (
                <input
                  type="text"
                  placeholder="Address (optional)"
                  value={address}
                  onChange={(e) => setAddress(e.target.value)}
                  className="auth-input"
                />
              )}
              {role === 'courier' && (
                <select
                  value={transportType}
                  onChange={(e) => setTransportType(e.target.value)}
                  className="auth-input"
                >
                  <option value="bicycle">Bicycle</option>
                  <option value="car">Car</option>
                  <option value="scooter">Scooter</option>
                  <option value="foot">Foot</option>
                </select>
              )}
            </>
          )}

          {error && <div className="auth-error">{error}</div>}

          <button type="submit" disabled={loading} className="auth-button">
            {loading ? 'LOADING...' : isSignUp ? 'SIGN UP' : 'SIGN IN'}
          </button>
        </form>

        <button
          onClick={() => setIsSignUp(!isSignUp)}
          className="auth-toggle"
        >
          {isSignUp ? 'Already have an account? Sign in' : "Don't have an account? Sign up"}
        </button>

        <button onClick={() => navigate('/')} className="auth-back">
          Back to portals
        </button>
      </div>
    </div>
  )
}
