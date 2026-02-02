// Stage configuration for each role
export const ORDER_STAGES = {
  customer: [
    { key: 'created', icon: '📝', label: 'Created' },
    { key: 'paid', icon: '💳', label: 'Paid' },
    { key: 'accepted', icon: '✅', label: 'Accepted' },
    { key: 'preparing', icon: '👨‍🍳', label: 'Preparing' },
    { key: 'pending_delivery', icon: '📦', label: 'Pending' },
    { key: 'picking', icon: '🚴', label: 'Picking' },
    { key: 'delivering', icon: '🚚', label: 'Delivering' },
    { key: 'completed', icon: '✨', label: 'Completed' }
  ],
  courier: [
    { key: 'paid', icon: '💳', label: 'Paid' },
    { key: 'accepted', icon: '✅', label: 'Accepted' },
    { key: 'preparing', icon: '👨‍🍳', label: 'Preparing' },
    { key: 'pending_delivery', icon: '📦', label: 'Ready' },
    { key: 'picking', icon: '🚴', label: 'Picking Up' },
    { key: 'delivering', icon: '🚚', label: 'On Route' },
    { key: 'completed', icon: '✨', label: 'Delivered' }
  ],
  restaurant: [
    { key: 'paid', icon: '💳', label: 'Order Paid' },
    { key: 'accepted', icon: '✅', label: 'Accepted' },
    { key: 'preparing', icon: '👨‍🍳', label: 'Preparing' },
    { key: 'pending_delivery', icon: '📦', label: 'Ready' },
    { key: 'completed', icon: '✨', label: 'Completed' }
  ]
}

// Backend status to stage key mapping
export const STATUS_STAGE_MAP = {
  // Positive flow
  CUSTOMER_CREATED: 'created',
  CUSTOMER_PAID: 'paid',
  KITCHEN_ACCEPTED: 'accepted',
  KITCHEN_PREPARING: 'preparing',
  DELIVERY_PENDING: 'pending_delivery',
  DELIVERY_PICKING: 'picking',
  DELIVERY_DELIVERING: 'delivering',
  ORDER_COMPLETED: 'completed',

  // Negative/alternative statuses
  CUSTOMER_CANCELLED: 'cancelled',
  KITCHEN_DENIED: 'denied',
  COURIER_REFUNDED: 'refunded',
  DELIVERY_DENIED: 'denied',
  DELIVERY_REFUNDED: 'refunded',
  ORDER_FAILED: 'failed'
}

// Failed/negative stage keys
const FAILED_STAGES = ['cancelled', 'denied', 'refunded', 'failed']

/**
 * Get the state of a stage based on current order status
 * @param {string} stageKey - The key of the stage to check
 * @param {string} currentStatus - The current order status (e.g., 'CUSTOMER_PAID')
 * @param {string} role - The current role ('customer', 'courier', 'restaurant')
 * @returns {'completed' | 'active' | 'pending' | 'failed'}
 */
export const getStageState = (stageKey, currentStatus, role) => {
  const currentStage = STATUS_STAGE_MAP[currentStatus] || 'created'
  const stages = ORDER_STAGES[role] || ORDER_STAGES.customer
  const stageOrder = stages.map(s => s.key)

  const currentIndex = stageOrder.indexOf(currentStage)
  const stageIndex = stageOrder.indexOf(stageKey)

  // Handle failed states
  if (FAILED_STAGES.includes(currentStage)) {
    if (stageKey === currentStage) return 'failed'
    if (stageIndex < currentIndex || currentIndex === -1) return 'completed'
    return 'pending'
  }

  // Handle edge cases
  if (currentIndex === -1) return 'pending'
  if (stageIndex === -1) return 'pending'

  // Normal flow
  if (stageIndex < currentIndex) return 'completed'
  if (stageIndex === currentIndex) return 'active'
  return 'pending'
}

/**
 * Get visible stages for a role based on current status
 * Shows all stages up to current + one pending, or all if completed/failed
 * @param {string} role - The current role
 * @param {string} currentStatus - The current order status
 * @returns {Array} Array of stage objects
 */
export const getVisibleStages = (role, currentStatus) => {
  const stages = ORDER_STAGES[role] || ORDER_STAGES.customer
  const currentStage = STATUS_STAGE_MAP[currentStatus] || 'created'
  const stageOrder = stages.map(s => s.key)
  const currentIndex = stageOrder.indexOf(currentStage)

  // Show all stages if completed or failed
  if (currentStage === 'completed' || FAILED_STAGES.includes(currentStage)) {
    return stages
  }

  // Show all stages up to and including current, plus one pending
  if (currentIndex !== -1) {
    return stages.filter((_, index) => index <= currentIndex + 1)
  }

  // Default: show first 3 stages
  return stages.slice(0, 3)
}

/**
 * Get a human-readable label for a status
 * @param {string} status - The backend status
 * @returns {string} Human-readable label
 */
export const getStatusLabel = (status) => {
  const labels = {
    CUSTOMER_CREATED: 'Created',
    CUSTOMER_PAID: 'Paid',
    KITCHEN_ACCEPTED: 'Accepted',
    KITCHEN_PREPARING: 'Preparing',
    DELIVERY_PENDING: 'Pending Delivery',
    DELIVERY_PICKING: 'Picking Up',
    DELIVERY_DELIVERING: 'Delivering',
    ORDER_COMPLETED: 'Completed',
    CUSTOMER_CANCELLED: 'Cancelled',
    KITCHEN_DENIED: 'Denied',
    COURIER_REFUNDED: 'Refunded',
    DELIVERY_DENIED: 'Delivery Denied',
    DELIVERY_REFUNDED: 'Refunded',
    ORDER_FAILED: 'Failed'
  }
  return labels[status] || status
}

/**
 * Get color type for a status
 * @param {string} status - The backend status
 * @returns {'success' | 'error' | 'accent' | 'default'}
 */
export const getStatusColor = (status) => {
  const failedStates = ['CUSTOMER_CANCELLED', 'KITCHEN_DENIED', 'COURIER_REFUNDED',
                       'DELIVERY_DENIED', 'DELIVERY_REFUNDED', 'ORDER_FAILED']

  if (failedStates.includes(status)) return 'error'
  if (status === 'ORDER_COMPLETED') return 'success'
  if (['DELIVERY_DELIVERING', 'KITCHEN_PREPARING', 'DELIVERY_PICKING'].includes(status)) return 'accent'
  return 'default'
}
