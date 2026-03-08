import { auth } from './auth.svelte';
import { toastStore } from './toast.svelte';
import type { Product } from '../models/Product';
import type { Purchase, PurchaseId } from '../models/Purchase';
import type { Invite } from '../models/Invite';

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
    const msg = 'Session expired. Please log in again.';
    toastStore.showError(msg);
    throw new Error(msg);
  }
  if (!response.ok) {
    const error = await response.text();
    const msg = error || `${method} ${path} failed`;
    toastStore.showError(msg);
    throw new Error(msg);
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
export const fetchInvites = (): Promise<{ invites: Invite[] }> => doGet('/invite');
export const sendInvite = (login: string): Promise<{ message: string }> => doPost('/invite', { login });

// Product API
export const fetchProducts = (): Promise<{ products: Product[] }> => doGet('/products');
export const createProduct = (product: Omit<Product, 'id' | 'user_id'>): Promise<Product> => doPost('/products', product);
export const updateProduct = (product: Product): Promise<Product> => doPut('/products', product);

// Purchase API
export const fetchPurchases = (): Promise<{ purchases: Purchase[] }> => doGet('/purchases');
export const createPurchase = (purchase: Omit<Purchase, 'id'>): Promise<Purchase> => doPost('/purchases', purchase);
export const deletePurchase = (purchaseId: PurchaseId) => doDelete('/purchases', { id: purchaseId.toString() });

export default API_URL;