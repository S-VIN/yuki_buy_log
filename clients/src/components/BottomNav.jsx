import { Menu } from 'antd';
import { PlusCircleOutlined, ProfileOutlined, SettingOutlined, ShoppingOutlined } from '@ant-design/icons';
import { useLocation, useNavigate } from 'react-router-dom';

const BottomNav = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const items = [
    { key: '/add', icon: <PlusCircleOutlined /> },
    { key: '/receipts', icon: <ProfileOutlined /> },
    { key: '/products', icon: <ShoppingOutlined /> },
    { key: '/settings', icon: <SettingOutlined /> },
  ];

  const selectedKey = location.pathname === '/' ? '/add' : location.pathname;

  return (
    <Menu
      mode="horizontal"
      selectedKeys={[selectedKey]}
      onClick={(e) => navigate(e.key)}
      items={items}
      style={{
        width: '100%',
        display: 'flex',
        justifyContent: 'space-around',
        height: 56,
        minHeight: 56,
        maxHeight: 56,
        flexShrink: 0,
        borderTop: '1px solid #f0f0f0',
        WebkitFlexShrink: 0,
        position: 'relative',
      }}
    />
  );
};

export default BottomNav;
