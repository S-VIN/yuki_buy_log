import { useEffect, useState } from 'react';
import { List, Tag } from 'antd';
import dayjs from 'dayjs';
import { authFetch } from '../api.js';
import { useAuth } from '../stores/AuthContext.jsx';

const Receipts = () => {
  const [purchases, setPurchases] = useState([]);
  const { user } = useAuth();

  useEffect(() => {
    const load = async () => {
      try {
        const [prodRes, purRes] = await Promise.all([
          authFetch('/products'),
          authFetch('/purchases'),
        ]);
        const prodJson = await prodRes.json();
        const products = Object.fromEntries(prodJson.products.map((p) => [p.id, p]));
        const purJson = await purRes.json();
        const list = purJson.purchases.map((p) => ({ ...p, product: products[p.product_id] }));
        setPurchases(list);
      } catch {
        // ignore
      }
    };
    load();
  }, []);

  return (
    <div style={{ padding: 8 }}>
      <List
        dataSource={purchases}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              title={item.product ? item.product.name : `Product ${item.product_id}`}
              description={`${item.price}₽ x ${item.quantity} — ${item.store} — ${dayjs(item.date).format('YYYY-MM-DD')}`}
            />
            {item.login && item.login !== user && (
              <Tag color="purple">{item.login}</Tag>
            )}
          </List.Item>
        )}
      />
    </div>
  );
};

export default Receipts;
