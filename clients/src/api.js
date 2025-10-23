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

export default API_URL;
