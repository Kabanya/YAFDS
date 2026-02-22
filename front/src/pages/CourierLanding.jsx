import { useNavigate } from 'react-router-dom'

export default function CourierLanding() {
  const navigate = useNavigate()

  return (
    <div className="role-landing courier-landing">
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
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M5 12h14"/>
              <path d="M12 5l7 7-7 7"/>
            </svg>
          </div>
        </header>

        <main className="role-main">
          <div className="role-title-section">
            <h1 className="role-title">
              <span className="title-line">COURIER</span>
              <span className="title-line">PORTAL</span>
            </h1>
            <p className="role-subtitle">Deliver orders and earn money on your schedule</p>
          </div>

          <div className="features-grid">
            <div className="feature-card" style={{ '--delay': '0ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <circle cx="12" cy="12" r="10"/>
                  <polyline points="12,6 12,12 16,14"/>
                </svg>
              </div>
              <h3>Flexible Hours</h3>
              <p>Work whenever you want, be your own boss</p>
            </div>

            <div className="feature-card" style={{ '--delay': '100ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
                </svg>
              </div>
              <h3>Weekly Payouts</h3>
              <p>Get paid every week with direct deposit</p>
            </div>

            <div className="feature-card" style={{ '--delay': '200ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                  <circle cx="12" cy="10" r="3"/>
                </svg>
              </div>
              <h3>Local Routes</h3>
              <p>Deliver in your neighborhood</p>
            </div>

            <div className="feature-card" style={{ '--delay': '300ms' }}>
              <div className="feature-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                  <polyline points="22,4 12,14.01 9,11.01"/>
                </svg>
              </div>
              <h3>Instant Earnings</h3>
              <p>See your earnings after each delivery</p>
            </div>
          </div>

          <div className="stats-row">
            <div className="stat-item" style={{ '--delay': '400ms' }}>
              <span className="stat-value mono">$25+</span>
              <span className="stat-label">AVG/HOUR</span>
            </div>
            <div className="stat-item" style={{ '--delay': '500ms' }}>
              <span className="stat-value mono">10K+</span>
              <span className="stat-label">COURIERS</span>
            </div>
            <div className="stat-item" style={{ '--delay': '600ms' }}>
              <span className="stat-value mono">24/7</span>
              <span className="stat-label">SUPPORT</span>
            </div>
          </div>

          <button
            className="role-cta"
            onClick={() => navigate('/courier/auth')}
            style={{ '--delay': '700ms' }}
          >
            <span>START DELIVERING</span>
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
              <span>HIRING</span>
            </span>
          </div>
          <div className="footer-line">
            <span className="footer-label">REGION</span>
            <span className="mono">ALL CITIES</span>
          </div>
        </footer>
      </div>
    </div>
  )
}
