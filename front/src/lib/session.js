const prefix = 'yafds:session'

export function sessionKey(role) {
  return `${prefix}:${role}`
}

export function normalizeSession(role, data) {
  return {
    role,
    id: data?.id || data?.ID || data?.Id || '',
    name: data?.name || data?.Name || '',
    walletAddress: data?.wallet_address || data?.walletAddress || data?.WalletAddress || '',
    address: data?.address || data?.Address || '',
    isActive: data?.is_active ?? data?.isActive ?? data?.IsActive ?? true,
    transportType: data?.transport_type || data?.transportType || data?.TransportType || '',
    token: data?.token || data?.Token || '',
    expiration: data?.expiration || data?.Expiration || null,
  }
}

export function readSession(role) {
  try {
    const raw = localStorage.getItem(sessionKey(role))
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function writeSession(role, data) {
  const session = normalizeSession(role, data)
  localStorage.setItem(sessionKey(role), JSON.stringify(session))
  return session
}

export function clearSession(role) {
  localStorage.removeItem(sessionKey(role))
}
