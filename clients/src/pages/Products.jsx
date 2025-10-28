import { Tag, Card, Row, Col, Typography, Input } from 'antd';
import { observer } from 'mobx-react-lite';
import { useState } from 'react';
import productStore from '../stores/ProductStore.jsx';

const { Text } = Typography;

const Products = observer(() => {
  const [searchText, setSearchText] = useState('');

  const filteredProducts = productStore.products.filter((product) =>
    product.name.toLowerCase().includes(searchText.toLowerCase())
  );

  if (productStore.products.length === 0) {
    return (
      <div
        style={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          height: '50vh',
          padding: 16,
        }}
      >
        <Text type="secondary" style={{ fontSize: 16 }}>
          No products added yet
        </Text>
      </div>
    );
  }

  return (
    <div style={{ paddingBottom: 8 }}>
      {/* Fixed header with search */}
      <div
        style={{
          position: 'sticky',
          top: 0,
          zIndex: 1,
          backgroundColor: '#fff',
          padding: '8px 8px 4px 8px',
        }}
      >
        <Card
          bodyStyle={{ padding: '8px 16px' }}
          style={{
            marginBottom: 4,
          }}
        >
          <Input
            placeholder="Search by product name..."
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            allowClear
          />
        </Card>
      </div>

      {/* Products list */}
      <div style={{ padding: 8 }}>
        {filteredProducts.length === 0 ? (
          <div
            style={{
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              padding: 32,
            }}
          >
            <Text type="secondary">No products match your search</Text>
          </div>
        ) : (
          filteredProducts.map((product) => (
            <Card
              key={product.id}
              style={{
                width: '100%',
                marginBottom: 8,
                padding: 0,
              }}
              bodyStyle={{ padding: '8px 16px' }}
            >
              <Row align="middle" gutter={[8, 4]}>
                <Col span={24}>
                  <Text strong>{product.name}</Text>
                </Col>
                <Col span={24}>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4 }}>
                    {product.brand && (
                      <Tag color="blue" size="small">
                        {product.brand}
                      </Tag>
                    )}
                    {product.volume && (
                      <Tag color="green" size="small">
                        {product.volume}
                      </Tag>
                    )}
                    {product.default_tags &&
                      product.default_tags.map((tag) => (
                        <Tag key={tag} color="orange" size="small">
                          {tag}
                        </Tag>
                      ))}
                  </div>
                </Col>
              </Row>
            </Card>
          ))
        )}
      </div>
    </div>
  );
});

export default Products;
