import { createContext, useContext, useEffect, useState } from 'react';
import PropTypes from 'prop-types';

const AuthContext = createContext({ token: null, login: () => {}, logout: () => {} });

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(null);

  useEffect(() => {
    const stored = localStorage.getItem('token');
    if (stored) {
      setToken(stored);
    }
  }, []);

  const login = (newToken) => {
    setToken(newToken);
    localStorage.setItem('token', newToken);
  };

  const logout = () => {
    setToken(null);
    localStorage.removeItem('token');
  };

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export const useAuth = () => useContext(AuthContext);

export default AuthContext;
