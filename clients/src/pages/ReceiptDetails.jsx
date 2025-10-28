import { Card, Typography, Button, Tag, message, Modal } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import ProductCardsWidget from '../widgets/ProductCardsWidget.jsx';
import dayjs from 'dayjs';
import productStore from '../stores/ProductStore.jsx';
import purchaseStore from '../stores/PurchaseStore.jsx';
import groupStore from '../stores/GroupStore.jsx';

const { Title, Text } = Typography;

const ReceiptDetails = observer(() => {
  const { id } = useParams();
  const navigate = useNavigate();

  const receiptPurchases = purchaseStore.getPurchasesByReceiptId(id);

  const receipt = receiptPurchases.length === 0 ? null : {
    id: id,
    date: receiptPurchases[0].date,
    shop: receiptPurchases[0].store,
    userId: receiptPurchases[0].user_id,
    items: receiptPurchases.map((p) => {
      const product = productStore.getProductById(String(p.product_id));
      return {
        ...p,
        product: product || { id: p.product_id, name: 'Unknown Product' },
      };
    }),
  };

  // Получаем информацию об участнике, который создал чек
  const memberInfo = receipt?.userId ? groupStore.getMemberInfo(receipt.userId) : null;
  const memberColor = memberInfo ? groupStore.getMemberColor(memberInfo.memberNumber) : null;

  const handleEditReceipt = async () => {
    try {
      // Удаляем все покупки из чека
      for (const purchase of receiptPurchases) {
        await purchaseStore.removePurchase(purchase.id);
      }

      message.success('Receipt purchases deleted. Redirecting to edit...');

      // Перенаправляем на экран добавления чека с предзаполненными данными
      navigate('/add', { state: { receipt } });
    } catch (error) {
      message.error(`Failed to delete purchases: ${error.message}`);
      console.error('Edit receipt error:', error);
    }
  };

  const handleDeleteReceipt = () => {
    Modal.confirm({
      title: 'Delete receipt?',
      content: 'All purchases from this receipt will be deleted. This action cannot be undone.',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          // Удаляем все покупки из чека
          for (const purchase of receiptPurchases) {
            await purchaseStore.removePurchase(purchase.id);
          }

          message.success('Receipt deleted successfully');

          // Перенаправляем на список чеков
          navigate('/receipts');
        } catch (error) {
          message.error(`Failed to delete receipt: ${error.message}`);
          console.error('Delete receipt error:', error);
        }
      },
    });
  };

  if (!receipt) {
    return <div style={{ padding: 16 }}>Receipt not found</div>;
  }

  return (
    <div style={{ padding: 8 }}>
      <Card style={{ marginBottom: 16, borderRadius: 8, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
              <Title level={3} style={{ margin: 0 }}>{receipt.shop}</Title>
              {memberInfo && (
                <Tag color={memberColor}>{memberInfo.login}</Tag>
              )}
            </div>
            <Text type="secondary">{dayjs(receipt.date).format('DD-MM-YYYY')}</Text>
          </div>
          <div style={{ display: 'flex', gap: 8 }}>
            <Button
              type="primary"
              icon={<EditOutlined />}
              onClick={handleEditReceipt}
            />
            <Button
              danger
              icon={<DeleteOutlined />}
              onClick={handleDeleteReceipt}
            />
          </div>
        </div>
      </Card>
      <ProductCardsWidget productListProp={receipt.items} />
    </div>
  );
});

export default ReceiptDetails;
