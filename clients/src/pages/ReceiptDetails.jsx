import { Button, Card, Typography } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

const ReceiptDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const receipts = JSON.parse(localStorage.getItem('receipts') || '[]');
  const receipt = receipts.find((r) => String(r.id) === id);

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
            type="text" 
            icon={<EditOutlined />} 
            size="large"
            onClick={() => navigate(`/edit/${id}`)}
          />
        </div>
      </Card>
      <ProductCardsWidget productListProp={receipt.items} />
    </div>
  );
};

export default ReceiptDetails;
