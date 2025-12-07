import { useNavigate } from 'react-router-dom'

export default function Landing() {
  const navigate = useNavigate()

  const portals = [
    { name: 'CUSTOMER', path: '/customer/auth', role: 'customer' },
    { name: 'COURIER', path: '/courier/auth', role: 'courier' },
    { name: 'RESTAURANT', path: '/restaurant/auth', role: 'restaurant' }
  ]

  return (
    <div className="landing-container">
      <div className="portal-grid">
        {portals.map((portal) => (
          <button
            key={portal.role}
            className="portal-button"
            onClick={() => navigate(portal.path, { state: { role: portal.role } })}
          >
            {portal.name}
          </button>
        ))}
      </div>
    </div>
  )
}
