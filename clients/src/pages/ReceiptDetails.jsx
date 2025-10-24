import { Card, Typography } from 'antd';
import { useParams } from 'react-router-dom';
import { useMemo } from 'react';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import dayjs from 'dayjs';
import { useData } from '../stores/DataContext.jsx';

const { Title, Text } = Typography;

const ReceiptDetails = () => {
  const { id } = useParams();
  const { purchases, products } = useData();

  const receipt = useMemo(() => {
    const receiptPurchases = purchases.filter((p) => String(p.receipt_id) === id);
    if (receiptPurchases.length === 0) return null;

    return {
      id: id,
      date: receiptPurchases[0].date,
      store: receiptPurchases[0].store,
      items: receiptPurchases.map((p) => {
        const product = products.find((prod) => prod.id === String(p.product_id));
        return {
          ...p,
          product: product || { id: p.product_id, name: 'Unknown Product' },
        };
      }),
    };
  }, [id, purchases, products]);

  if (!receipt) {
    return <div style={{ padding: 16 }}>Receipt not found</div>;
  }

  return (
    <div style={{ padding: 8 }}>
      <Card style={{ marginBottom: 16, borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title level={3} style={{ margin: 0, marginBottom: 4 }}>{receipt.store}</Title>
            <Text type="secondary">{dayjs(receipt.date).format('DD-MM-YYYY')}</Text>
          </div>
        </div>
      </Card>
      <ProductCardsWidget productListProp={receipt.items} />
    </div>
  );
};

export default ReceiptDetails;
