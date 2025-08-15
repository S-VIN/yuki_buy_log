import { createContext, useContext, useEffect, useState } from 'react';
import PropTypes from 'prop-types';

const AuthContext = createContext({ token: null, user: null, login: () => {}, logout: () => {} });

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(null);
  const [user, setUser] = useState(null);

  useEffect(() => {
    const storedToken = localStorage.getItem('token');
    const storedLogin = localStorage.getItem('login');
    if (storedToken) {
      setToken(storedToken);
    }
    if (storedLogin) {
      setUser(storedLogin);
    }
  }, []);

  const login = (newToken, loginName) => {
    setToken(newToken);
    setUser(loginName);
    localStorage.setItem('token', newToken);
    localStorage.setItem('login', loginName);
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
    localStorage.removeItem('login');
  };

  return (
    <AuthContext.Provider value={{ token, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export const useAuth = () => useContext(AuthContext);

export default AuthContext;
