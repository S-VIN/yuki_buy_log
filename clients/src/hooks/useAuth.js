import { useContext } from 'react';
import AuthContext from '../stores/AuthContext.jsx';

export const useAuth = () => useContext(AuthContext);
