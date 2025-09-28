import { useEffect, useState } from 'react';
import { Tag, Card, Row, Col, Typography } from 'antd';
import dayjs from 'dayjs';
import { useNavigate } from 'react-router-dom';

const { Text } = Typography;

const Receipts = () => {
  const [receipts, setReceipts] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const data = JSON.parse(localStorage.getItem('receipts') || '[]');
    data.sort((a, b) => dayjs(b.date).diff(dayjs(a.date)));
    setReceipts(data);
  }, []);

  const getTags = (items) => {
    const set = new Set();
    items.forEach((i) => i.tags && i.tags.forEach((t) => set.add(t)));
    return Array.from(set);
  };

  const getSum = (items) => items.reduce((sum, i) => sum + i.price * i.quantity, 0);

  const groupReceiptsByDate = (receipts) => {
    const groups = {};
    receipts.forEach((receipt) => {
      const date = receipt.date;
      if (!groups[date]) {
        groups[date] = [];
      }
      groups[date].push(receipt);
    });
    return groups;
  };

  const groupedReceipts = groupReceiptsByDate(receipts);
  const sortedDates = Object.keys(groupedReceipts).sort((a, b) => dayjs(b).diff(dayjs(a)));

  return (
    <div style={{ padding: 8 }}>
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
          {groupedReceipts[date].map((item) => (
            <Card 
              key={item.id}
              style={{ 
                width: '100%', 
                marginBottom: 8, 
                cursor: 'pointer',
                padding: 0
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
          ))}
        </div>
      ))}
    </div>
  );
};

export default Receipts;
