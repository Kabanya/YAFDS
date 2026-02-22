import { useNavigate } from 'react-router-dom'

export default function Landing() {
  const navigate = useNavigate()

  const portals = [
    {
      name: 'CUSTOMER',
      path: '/customer/auth',
      role: 'customer',
      icon: (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
          <circle cx="9" cy="7" r="4"/>
          <path d="M3 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2"/>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          <path d="M21 21v-2a4 4 0 0 0-3-3.87"/>
        </svg>
      ),
      description: 'Order food from restaurants'
    },
    {
      name: 'COURIER',
      path: '/courier/auth',
      role: 'courier',
      icon: (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <path d="M5 12h14"/>
          <path d="M12 5l7 7-7 7"/>
        </svg>
      ),
      description: 'Deliver orders and earn'
    },
    {
      name: 'RESTAURANT',
      path: '/restaurant/auth',
      role: 'restaurant',
      icon: (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
          <path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/>
          <path d="M7 2v20"/>
          <path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"/>
        </svg>
      ),
      description: 'Manage menu and orders'
    }
  ]

  return (
    <div className="landing-container">
      <div className="landing-content">
        <header className="landing-header">
          <div className="landing-logo">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 2L2 7l10 5 10-5-10-5z"/>
              <path d="M2 17l10 5 10-5"/>
              <path d="M2 12l10 5 10-5"/>
            </svg>
          </div>
          <div className="landing-title">
            <h1>YAFDS</h1>
            <p className="landing-subtitle">Yet Another Food Delivery System</p>
          </div>
        </header>

        <div className="portal-grid">
          {portals.map((portal, index) => (
            <button
              key={portal.role}
              className="portal-button"
              onClick={() => navigate(portal.path, { state: { role: portal.role } })}
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <div className="portal-icon">{portal.icon}</div>
              <span className="portal-name">{portal.name}</span>
              <span className="portal-desc">{portal.description}</span>
              <div className="portal-status">
                <span className="status-dot" />
                <span className="status-text">Online</span>
              </div>
            </button>
          ))}
        </div>

        <footer className="landing-footer">
          <div className="footer-info">
            <span className="footer-label">SYSTEM STATUS</span>
            <span className="footer-value mono">ALL SERVICES UP</span>
          </div>
          <div className="footer-info">
            <span className="footer-label">VERSION</span>
            <span className="footer-value mono">1.0.0</span>
          </div>
        </footer>
      </div>
    </div>
  )
}
