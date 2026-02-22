import { useNavigate } from 'react-router-dom'

export default function CustomerLanding() {
  const navigate = useNavigate()

  return (
    <div className="role-landing customer-landing">
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
              <circle cx="9" cy="7" r="4"/>
              <path d="M3 21v-2a4 4 0 0 1 4-4h4a4 4 0 0 1 4 4v2"/>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              <path d="M21 21v-2a4 4 0 0 0-3-3.87"/>
            </svg>
          </div>
        </header>

        <main className="role-main">
          <div className="role-title-section">
            <h1 className="role-title">
              <span className="title-line">CUSTOMER</span>
              <span className="title-line">PORTAL</span>
            </h1>
            <p className="role-subtitle">Order food from your favorite restaurants</p>
          </div>

          <div className="features-grid">
            <div className="feature-card" style={{ '--delay': '0ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                  <polyline points="9,22 9,12 15,12 15,22"/>
                </svg>
              </div>
              <h3>Browse Restaurants</h3>
              <p>Explore nearby restaurants and their menus</p>
            </div>

            <div className="feature-card" style={{ '--delay': '100ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <circle cx="9" cy="21" r="1"/>
                  <circle cx="20" cy="21" r="1"/>
                  <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
                </svg>
              </div>
              <h3>Easy Ordering</h3>
              <p>Quick checkout with multiple payment options</p>
            </div>

            <div className="feature-card" style={{ '--delay': '200ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                </svg>
              </div>
              <h3>Secure Payments</h3>
              <p>Your transactions are always protected</p>
            </div>

            <div className="feature-card" style={{ '--delay': '300ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
              </div>
              <h3>Track Orders</h3>
              <p>Real-time updates on your delivery status</p>
            </div>
          </div>

          <button
            className="role-cta"
            onClick={() => navigate('/customer/auth')}
            style={{ '--delay': '400ms' }}
          >
            <span>GET STARTED</span>
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
              <span>ONLINE</span>
            </span>
          </div>
          <div className="footer-line">
            <span className="footer-label">REGION</span>
            <span className="mono">ALL LOCATIONS</span>
          </div>
        </footer>
      </div>
    </div>
  )
}
