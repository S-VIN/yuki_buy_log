import { Layout } from 'antd';
import { Route, Routes, useLocation } from 'react-router-dom';
import BottomNav from './components/BottomNav.jsx';
import Login from './pages/Login.jsx';
import AddReceipt from './pages/AddReceipt.jsx';
import Receipts from './pages/Receipts.jsx';
import ReceiptDetails from './pages/ReceiptDetails.jsx';
import Settings from './pages/Settings.jsx';

const App = () => {
  const location = useLocation();
  const hideNav = location.pathname === '/login';

  return (
    <Layout style={{ minHeight: '100vh', background: '#fff' }}>
      <Layout.Content style={{ paddingBottom: hideNav ? 0 : 56 }}>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={<AddReceipt />} />
          <Route path="/add" element={<AddReceipt />} />
          <Route path="/edit/:id" element={<AddReceipt />} />
          <Route path="/receipts" element={<Receipts />} />
          <Route path="/receipts/:id" element={<ReceiptDetails />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Layout.Content>
      {!hideNav && <BottomNav />}
    </Layout>
  );
};

export default App;
