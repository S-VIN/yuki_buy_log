import { Button } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import dayjs from 'dayjs';

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
      <h3>{dayjs(receipt.date).format('YYYY-MM-DD')} — {receipt.shop}</h3>
      <ProductCardsWidget productListProp={receipt.items} />
      <Button type="primary" block style={{ marginTop: 16 }} onClick={() => navigate(`/edit/${id}`)}>
        Edit
      </Button>
    </div>
  );
};

export default ReceiptDetails;
