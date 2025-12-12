import { useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

export default function Auth() {
  const location = useLocation()
  const navigate = useNavigate()
  const role = location.state?.role || 'customer'

  const customerApiUrl = import.meta.env.VITE_CUSTOMER_API_URL || 'http://localhost:8081'

  const [isSignUp, setIsSignUp] = useState(false)
  const [name, setName] = useState('')
  const [walletAddress, setWalletAddress] = useState('')
  const [address, setAddress] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleAuth = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      if (isSignUp) {
        const response = await fetch(`${customerApiUrl}/register`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            name,
            wallet_address: walletAddress,
            address,
            password,
          }),
        })

        const data = await response.json()
        if (!response.ok) {
          throw new Error(data?.error_message || 'Failed to register')
        }

        const profile = {
          id: data.id,
          name: name || walletAddress,
          wallet_address: walletAddress,
          address: address || walletAddress,
          role,
        }

        localStorage.setItem('currentUser', JSON.stringify(profile))
        navigate(`/${role}/dashboard`, { state: { user: profile } })
      } else {
        const response = await fetch(`${customerApiUrl}/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            wallet_address: walletAddress,
            password,
          }),
        })

        const data = await response.json()
        if (!response.ok) throw new Error(data?.error_message || 'Invalid wallet address or password')

        const profile = {
          id: data.id,
          name: data.name,
          wallet_address: data.wallet_address,
          address: data.address,
          role,
        }

        localStorage.setItem('currentUser', JSON.stringify(profile))
        navigate(`/${role}/dashboard`, { state: { user: profile } })
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
                placeholder="Name (optional)"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="auth-input"
              />
              <input
                type="text"
                placeholder="Address (optional)"
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                className="auth-input"
              />
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
