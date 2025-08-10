import React, { useMemo } from "react";
import { List, Tag, Card, Row, Col, Typography } from "antd";
import moment from "moment";

const { Text } = Typography;

type Check = {
  id: number;
  shop: string;
  date: string; // ISO date
  sum_price: number;
  tags: string[];
};

const initialData: Check[] = [
  {
    id: 1,
    shop: "Supermarket A",
    date: "2025-07-10",
    sum_price: 150.5,
    tags: ["food", "groceries"],
  },
  {
    id: 2,
    shop: "ElectroMart",
    date: "2025-07-11",
    sum_price: 2300,
    tags: ["electronics"],
  },
  {
    id: 3,
    shop: "Cafe B",
    date: "2025-07-09",
    sum_price: 25.75,
    tags: ["food", "coffee"],
  },
];

const ViewChecksScreen: React.FC = () => {
  const sortedData = useMemo(() => {
    return [...initialData].sort((a, b) =>
      moment(b.date).diff(moment(a.date))
    );
  }, []);

  return (
    <div>
      <List
        dataSource={sortedData}
        renderItem={(item) => (
          <List.Item style={{ padding: 0 }}>
            <Card style={{ width: "100%" }}>
              <Row align="middle">
                {/* Магазин слева */}
                <Col span={6} style={{ textAlign: "left" }}>
                  <Text strong>{item.shop}</Text>
                </Col>

                {/* Теги по центру */}
                <Col span={12}>
                  <div >
                    {item.tags.map((tag) => (
                      <Tag key={tag} color="orange">
                        {tag}
                      </Tag>
                    ))}
                  </div>
                </Col>

                {/* Сумма и дата справа */}
                <Col span={6} style={{ textAlign: "right" }}>
                  <Text>{moment(item.date).format("YYYY-MM-DD")}</Text>
                  <br />
                  <Text>{item.sum_price.toLocaleString()} ₽</Text>
                </Col>
              </Row>
            </Card>
          </List.Item>
        )}
      />
    </div>
  );
};

export default ViewChecksScreen;

