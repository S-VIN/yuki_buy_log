import { Card, Typography, Button } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import { EditOutlined } from '@ant-design/icons';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import dayjs from 'dayjs';
import { useProductStore, usePurchaseStore } from '../stores/DataContext.jsx';

const { Title, Text } = Typography;

const ReceiptDetails = observer(() => {
  const { id } = useParams();
  const navigate = useNavigate();
  const productStore = useProductStore();
  const purchaseStore = usePurchaseStore();

  const receiptPurchases = purchaseStore.getPurchasesByReceiptId(id);

  const receipt = receiptPurchases.length === 0 ? null : {
    id: id,
    date: receiptPurchases[0].date,
    shop: receiptPurchases[0].store,
    items: receiptPurchases.map((p) => {
      const product = productStore.getProductById(String(p.product_id));
      return {
        ...p,
        product: product || { id: p.product_id, name: 'Unknown Product' },
      };
    }),
  };

  if (!receipt) {
    return <div style={{ padding: 16 }}>Receipt not found</div>;
  }

  return (
    <div style={{ padding: 8 }}>
      <Card style={{ marginBottom: 16, borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title level={3} style={{ margin: 0, marginBottom: 4 }}>{receipt.shop}</Title>
            <Text type="secondary">{dayjs(receipt.date).format('DD-MM-YYYY')}</Text>
          </div>
          <Button
            type="primary"
            icon={<EditOutlined />}
            onClick={() => navigate('/add', { state: { receipt } })}
          >
            Edit
          </Button>
        </div>
      </Card>
      <ProductCardsWidget productListProp={receipt.items} />
    </div>
  );
});

export default ReceiptDetails;
