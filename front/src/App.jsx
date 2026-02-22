import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Landing from './pages/Landing'
import CustomerLanding from './pages/CustomerLanding'
import CourierLanding from './pages/CourierLanding'
import RestaurantLanding from './pages/RestaurantLanding'
import Auth from './pages/Auth'
import Dashboard from './pages/Dashboard'
import ThemeToggle from './components/ThemeToggle'

export default function App() {
  return (
    <Router>
      <ThemeToggle />
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/customer" element={<CustomerLanding />} />
        <Route path="/courier" element={<CourierLanding />} />
        <Route path="/restaurant" element={<RestaurantLanding />} />
        <Route path="/customer/auth" element={<Auth />} />
        <Route path="/courier/auth" element={<Auth />} />
        <Route path="/restaurant/auth" element={<Auth />} />
        <Route path="/:role/dashboard" element={<Dashboard />} />
      </Routes>
    </Router>
  )
}
