import { Tag, Card, Row, Col, Typography } from 'antd';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';
import { observer } from 'mobx-react-lite';
import purchaseStore from '../stores/PurchaseStore.jsx';
import groupStore from '../stores/GroupStore.jsx';

const { Text } = Typography;

const Receipts = observer(() => {
  const navigate = useNavigate();

  const receipts = purchaseStore.receipts.sort((a, b) => dayjs(b.date).diff(dayjs(a.date)));

  const getTags = (items) => {
    const set = new Set();
    items.forEach((i) => i.tags && i.tags.forEach((t) => set.add(t)));
    return Array.from(set);
  };

  const getSum = (items) => items.reduce((sum, i) => sum + i.price * i.quantity, 0);

  const groupReceiptsByDate = (receipts) => {
    const groups = {};
    receipts.forEach((receipt) => {
      const date = dayjs(receipt.date).format('YYYY-MM-DD');
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(receipt);
    });
    return groups;
  };

  const groupedReceipts = groupReceiptsByDate(receipts);
  const sortedDates = Object.keys(groupedReceipts).sort((a, b) => dayjs(b).diff(dayjs(a)));

  // Функция для получения цвета полосы для чека
  const getMemberBarColor = (receipt) => {
    // Проверяем, находимся ли мы в мультиюзерной группе
    if (!groupStore.isInMultiUserGroup) {
      return null;
    }

    // Получаем user_id из первой покупки в чеке
    const userId = receipt.items[0]?.user_id;
    if (!userId) {
      return null;
    }

    // Если это чек текущего пользователя, не показываем полосу
    if (groupStore.isCurrentUserPurchase(userId)) {
      return null;
    }

    // Получаем member_number для этого user_id
    const memberNumber = groupStore.getMemberNumberByUserId(userId);
    if (!memberNumber) {
      return null;
    }

    return groupStore.getMemberHexColor(memberNumber);
  };

  if (receipts.length === 0) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '50vh',
        padding: 16
      }}>
        <Text type="secondary" style={{ fontSize: 16 }}>
          No receipts added yet
        </Text>
      </div>
    );
  }

  return (
    <div style={{
      height: '100vh',
      overflowY: 'auto',
      padding: 8,
      paddingBottom: 64
    }}>
      {sortedDates.map((date) => (
        <div key={date}>
              <div style={{
                display: 'flex',
                alignItems: 'center',
                paddingLeft: '16px',
                paddingRight: '16px',
                marginBottom: '4px',
                marginTop: '8px',
              }}>
                  <Text strong>{dayjs(date).format('DD-MM-YYYY')}</Text>
                  <div style={{ flex: 1, height: '1px', backgroundColor: '#d9d9d9', margin: '0 16px' }}></div>
                  <Text strong> {groupedReceipts[date].reduce((sum, receipt) => sum + getSum(receipt.items), 0)} ₽ </Text>
              </div>
          {groupedReceipts[date].map((item) => {
            const barColor = getMemberBarColor(item);
            return (
              <Card
                key={item.id}
                style={{
                  width: '100%',
                  marginBottom: 8,
                  cursor: 'pointer',
                  padding: 0,
                  borderLeft: barColor ? `4px solid ${barColor}` : undefined,
                }}
                bodyStyle={{ padding: '8px 16px' }}
                onClick={() => navigate(`/receipts/${item.id}`)}
              >
                <Row align="middle">
                  <Col span={8} style={{ textAlign: 'left' }}>
                    <Text>{item.shop}</Text>
                  </Col>
                  <Col span={10}>
                    <div>
                      {getTags(item.items).map((tag) => (
                        <Tag key={tag} color="orange" size="small">
                          {tag}
                        </Tag>
                      ))}
                    </div>
                  </Col>
                  <Col span={6} style={{ textAlign: 'right' }}>
                    <Text>{getSum(item.items)} ₽</Text>
                  </Col>
                </Row>
              </Card>
            );
          })}
        </div>
      ))}
    </div>
  );
});

export default Receipts;
