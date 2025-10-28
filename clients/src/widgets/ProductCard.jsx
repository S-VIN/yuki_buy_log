import PropTypes from 'prop-types';
import { Card, Row, Col, Tag, Typography } from 'antd';

const { Text } = Typography;

const ProductCard = ({ product }) => (
  <Card
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
);

ProductCard.propTypes = {
  product: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    brand: PropTypes.string,
    volume: PropTypes.string,
    default_tags: PropTypes.arrayOf(PropTypes.string),
  }).isRequired,
};

export default ProductCard;
