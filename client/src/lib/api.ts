import { auth } from './auth.svelte';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

function getAuthHeaders(): Record<string, string> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (auth.token) {
    headers.Authorization = `Bearer ${auth.token}`;
  }
  return headers;
}

async function handleResponse(response: Response, method: string, path: string): Promise<Response> {
  if (response.status === 401) {
    auth.logout();
    throw new Error('Session expired. Please log in again.');
  }
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || `${method} ${path} failed`);
  }
  return response;
}

async function doGet(path: string) {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'GET',
    headers: getAuthHeaders(),
  });
  return (await handleResponse(response, 'GET', path)).json();
}

async function doPost(path: string, data?: unknown) {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  return (await handleResponse(response, 'POST', path)).json();
}

async function doDelete(path: string, data?: unknown) {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'DELETE',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  const handled = await handleResponse(response, 'DELETE', path);
  if (handled.status === 204) return true;
  return handled.json();
}

async function doPut(path: string, data?: unknown) {
  const response = await fetch(`${API_URL}${path}`, {
    method: 'PUT',
    headers: getAuthHeaders(),
    body: JSON.stringify(data),
  });
  return (await handleResponse(response, 'PUT', path)).json();
}

// Auth API (no token needed)
export async function apiLogin(login: string, password: string): Promise<string> {
  const response = await fetch(`${API_URL}/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password }),
  });
  if (!response.ok) throw new Error('Login failed');
  const data = await response.json();
  return data.token;
}

export async function apiRegister(login: string, password: string): Promise<string> {
  const response = await fetch(`${API_URL}/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ login, password }),
  });
  if (!response.ok) throw new Error('Registration failed');
  const data = await response.json();
  return data.token;
}

// Group API
export const fetchGroupMembers = () => doGet('/group');
export const leaveGroup = () => doDelete('/group', {});

// Invite API
export const fetchInvites = () => doGet('/invite');
export const sendInvite = (login: string) => doPost('/invite', { login });

// Product API
export const fetchProducts = () => doGet('/products');
export const createProduct = (product: unknown) => doPost('/products', product);
export const updateProduct = (product: unknown) => doPut('/products', product);

// Purchase API
export const fetchPurchases = () => doGet('/purchases');
export const createPurchase = (purchase: unknown) => doPost('/purchases', purchase);
export const deletePurchase = (purchaseId: string) => doDelete('/purchases', { id: purchaseId });

export default API_URL;