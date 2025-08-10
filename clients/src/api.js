const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const authFetch = (path, options = {}) => {
  const token = localStorage.getItem('token');
  const headers = { ...options.headers };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return fetch(`${API_URL}${path}`, { ...options, headers });
};

export default API_URL;
