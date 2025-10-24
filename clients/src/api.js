const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const authFetch = (path, options = {}) => {
  const token = localStorage.getItem('token');
  const headers = { ...options.headers };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return fetch(`${API_URL}${path}`, { ...options, headers });
};

// Group API functions
export const fetchGroupMembers = async () => {
  const response = await authFetch('/group');
  if (!response.ok) {
    throw new Error('Failed to fetch group members');
  }
  return response.json();
};

// Invite API functions
export const fetchInvites = async () => {
  const response = await authFetch('/invite');
  if (!response.ok) {
    throw new Error('Failed to fetch invites');
  }
  return response.json();
};

export const sendInvite = async (login) => {
  const response = await authFetch('/invite', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login }),
  });
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to send invite');
  }
  return response.json();
};

// Product API functions
export const fetchProducts = async () => {
  const response = await authFetch('/products');
  if (!response.ok) {
    throw new Error('Failed to fetch products');
  }
  return response.json();
};

export const createProduct = async (product) => {
  const response = await authFetch('/products', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(product),
  });
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to create product');
  }
  return response.json();
};

// Purchase API functions
export const fetchPurchases = async () => {
  const response = await authFetch('/purchases');
  if (!response.ok) {
    throw new Error('Failed to fetch purchases');
  }
  return response.json();
};

export const createPurchase = async (purchase) => {
  const response = await authFetch('/purchases', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(purchase),
  });
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to create purchase');
  }
  return response.json();
};

export const deletePurchase = async (purchaseId) => {
  const response = await authFetch('/purchases', {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ id: purchaseId }),
  });
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Failed to delete purchase');
  }
  // DELETE returns 204 No Content, so no body to parse
  return true;
};

export default API_URL;
