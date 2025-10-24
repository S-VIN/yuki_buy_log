import { Layout } from 'antd';
import { Route, Routes, useLocation, Navigate } from 'react-router-dom';
import { useEffect } from 'react';
import BottomNav from './components/BottomNav.jsx';
import Login from './pages/Login.jsx';
import AddReceipt from './pages/AddReceipt.jsx';
import Receipts from './pages/Receipts.jsx';
import ReceiptDetails from './pages/ReceiptDetails.jsx';
import Settings from './pages/Settings.jsx';
import { useAuth } from './stores/AuthContext.jsx';
import { useData } from './stores/DataContext.jsx';

const App = () => {
  const location = useLocation();
  const { token } = useAuth();
  const { loadData } = useData();
  const hideNav = !token || location.pathname === '/login';

  useEffect(() => {
    if (token) {
      loadData();
    }
  }, [token]);

  return (
    <Layout style={{ minHeight: '100vh', background: '#fff' }}>
      <Layout.Content style={{ paddingBottom: hideNav ? 0 : 56 }}>
        <Routes>
          <Route path="/login" element={!token ? <Login /> : <Navigate to="/" replace />} />
          {token && (
            <>
              <Route path="/" element={<AddReceipt />} />
              <Route path="/add" element={<AddReceipt />} />
              <Route path="/receipts" element={<Receipts />} />
              <Route path="/receipts/:id" element={<ReceiptDetails />} />
              <Route path="/settings" element={<Settings />} />
            </>
          )}
          <Route path="*" element={<Navigate to={token ? '/' : '/login'} replace />} />
        </Routes>
      </Layout.Content>
      {!hideNav && token && <BottomNav />}
    </Layout>
  );
};

export default App;
