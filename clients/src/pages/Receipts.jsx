import { useEffect, useState } from 'react';
import { List, Tag, Card, Row, Col, Typography } from 'antd';
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

  return (
    <div style={{ padding: 8 }}>
      <List
        dataSource={receipts}
        renderItem={(item) => (
          <List.Item style={{ padding: 0 }} onClick={() => navigate(`/receipts/${item.id}`)}>
            <Card style={{ width: '100%' }}>
              <Row align="middle">
                <Col span={6} style={{ textAlign: 'left' }}>
                  <Text strong>{item.shop}</Text>
                </Col>
                <Col span={12}>
                  <div>
                    {getTags(item.items).map((tag) => (
                      <Tag key={tag} color="orange">
                        {tag}
                      </Tag>
                    ))}
                  </div>
                </Col>
                <Col span={6} style={{ textAlign: 'right' }}>
                  <Text>{dayjs(item.date).format('YYYY-MM-DD')}</Text>
                  <br />
                  <Text>{getSum(item.items)} ₽</Text>
                </Col>
              </Row>
            </Card>
          </List.Item>
        )}
      />
    </div>
  );
};

export default Receipts;
