import { Card, Typography } from 'antd';
import { useParams } from 'react-router-dom';
import { useMemo } from 'react';
import { observer } from 'mobx-react-lite';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import dayjs from 'dayjs';
import { useProductStore, usePurchaseStore } from '../stores/DataContext.jsx';

const { Title, Text } = Typography;

const ReceiptDetails = observer(() => {
  const { id } = useParams();
  const productStore = useProductStore();
  const purchaseStore = usePurchaseStore();

  const receipt = useMemo(() => {
    const receiptPurchases = purchaseStore.getPurchasesByReceiptId(id);
    if (receiptPurchases.length === 0) return null;

    return {
      id: id,
      date: receiptPurchases[0].date,
      store: receiptPurchases[0].store,
      items: receiptPurchases.map((p) => {
        const product = productStore.getProductById(String(p.product_id));
        return {
          ...p,
          product: product || { id: p.product_id, name: 'Unknown Product' },
        };
      }),
    };
  }, [id, productStore.products, purchaseStore.purchases]);

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
});

export default ReceiptDetails;
