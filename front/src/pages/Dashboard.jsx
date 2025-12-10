import { useEffect, useState } from 'react'
import { useParams, useNavigate, useLocation } from 'react-router-dom'

export default function Dashboard() {
  const { role } = useParams()
  const navigate = useNavigate()
  const location = useLocation()
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

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

    setUser(savedUser)
    setLoading(false)
  }

  const handleSignOut = () => {
    localStorage.removeItem('currentUser')
    navigate('/')
  }

  if (loading) {
    return (
      <div className="dashboard-container">
        <div className="dashboard-content">LOADING...</div>
      </div>
    )
  }

  return (
    <div className="dashboard-container">
      <div className="dashboard-content">
        <h1 className="dashboard-title">{role.toUpperCase()} DASHBOARD</h1>
        <p className="dashboard-wallet">{user?.wallet_address}</p>
        <button onClick={handleSignOut} className="dashboard-button">
          SIGN OUT
        </button>
      </div>
    </div>
  )
}
