import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Landing from './pages/Landing'
import Auth from './pages/Auth'
import Dashboard from './pages/Dashboard'

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/customer/auth" element={<Auth />} />
        <Route path="/courier/auth" element={<Auth />} />
        <Route path="/restaurant/auth" element={<Auth />} />
        <Route path="/:role/dashboard" element={<Dashboard />} />
      </Routes>
    </Router>
  )
}
