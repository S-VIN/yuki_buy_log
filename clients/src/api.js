const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const getAuthHeaders = () => {
  const token = localStorage.getItem('token');
  const headers = { 'Content-Type': 'application/json' };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return headers;
};

// Base HTTP methods
const doGet = async (path) => {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `GET ${path} failed`);
  }

  return response.json();
};

const doPost = async (path, data) => {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `POST ${path} failed`);
  }

  return response.json();
};

const doDelete = async (path, data) => {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `DELETE ${path} failed`);
  }

  // DELETE may return 204 No Content
  if (response.status === 204) {
    return true;
  }

  return response.json();
};

// Group API
export const fetchGroupMembers = () => doGet('/group');
export const leaveGroup = () => doDelete('/group', {});

// Invite API
export const fetchInvites = () => doGet('/invite');
export const sendInvite = (login) => doPost('/invite', { login });

// Product API
export const fetchProducts = () => doGet('/products');
export const createProduct = (product) => doPost('/products', product);

// Purchase API
export const fetchPurchases = () => doGet('/purchases');
export const createPurchase = (purchase) => doPost('/purchases', purchase);
export const deletePurchase = (purchaseId) => doDelete('/purchases', { id: purchaseId });

export default API_URL;
