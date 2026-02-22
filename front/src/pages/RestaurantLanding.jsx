import { useNavigate } from 'react-router-dom'

export default function RestaurantLanding() {
  const navigate = useNavigate()

  return (
    <div className="role-landing restaurant-landing">
      <div className="role-landing-bg" />
      <div className="role-landing-content">
        <header className="role-header">
          <button className="back-button" onClick={() => navigate('/')}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M19 12H5"/>
              <path d="M12 19l-7-7 7-7"/>
            </svg>
            <span>BACK</span>
          </button>
          <div className="role-logo">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/>
              <path d="M7 2v20"/>
              <path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"/>
            </svg>
          </div>
        </header>

        <main className="role-main">
          <div className="role-title-section">
            <h1 className="role-title">
              <span className="title-line">RESTAURANT</span>
              <span className="title-line">PORTAL</span>
            </h1>
            <p className="role-subtitle">Manage your menu and orders efficiently</p>
          </div>

          <div className="features-grid">
            <div className="feature-card" style={{ '--delay': '0ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </div>
              <h3>Menu Management</h3>
              <p>Update your menu items and prices in real-time</p>
            </div>

            <div className="feature-card" style={{ '--delay': '100ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                  <line x1="16" y1="13" x2="8" y2="13"/>
                  <line x1="16" y1="17" x2="8" y2="17"/>
                  <polyline points="10,9 9,9 8,9"/>
                </svg>
              </div>
              <h3>Order Processing</h3>
              <p>Accept, reject, or manage incoming orders easily</p>
            </div>

            <div className="feature-card" style={{ '--delay': '200ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <line x1="18" y1="20" x2="18" y2="10"/>
                  <line x1="12" y1="20" x2="12" y2="4"/>
                  <line x1="6" y1="20" x2="6" y2="14"/>
                </svg>
              </div>
              <h3>Analytics</h3>
              <p>Track sales, popular items, and performance</p>
            </div>

            <div className="feature-card" style={{ '--delay': '300ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </div>
              <h3>Staff Management</h3>
              <p>Manage your team and their permissions</p>
            </div>
          </div>

          <div className="stats-row">
            <div className="stat-item" style={{ '--delay': '400ms' }}>
              <span className="stat-value mono">500+</span>
              <span className="stat-label">RESTAURANTS</span>
            </div>
            <div className="stat-item" style={{ '--delay': '500ms' }}>
              <span className="stat-value mono">1M+</span>
              <span className="stat-label">ORDERS/MONTH</span>
            </div>
            <div className="stat-item" style={{ '--delay': '600ms' }}>
              <span className="stat-value mono">15%</span>
              <span className="stat-label">AVG GROWTH</span>
            </div>
          </div>

          <button
            className="role-cta"
            onClick={() => navigate('/restaurant/auth')}
            style={{ '--delay': '700ms' }}
          >
            <span>PARTNER WITH US</span>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M5 12h14"/>
              <path d="M12 5l7 7-7 7"/>
            </svg>
          </button>
        </main>

        <footer className="role-footer">
          <div className="footer-line">
            <span className="footer-label">STATUS</span>
            <span className="status-indicator">
              <span className="status-dot online" />
              <span>ACCEPTING</span>
            </span>
          </div>
          <div className="footer-line">
            <span className="footer-label">REGION</span>
            <span className="mono">NATIONWIDE</span>
          </div>
        </footer>
      </div>
    </div>
  )
}
