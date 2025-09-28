/* eslint-disable react/prop-types */
import { Card, Row, Col, Tag, Button } from 'antd';
import { DeleteTwoTone, EditTwoTone } from '@ant-design/icons';

const ProductCardsWidget = ({ productListProp, onDelete, onEdit }) => (
  <div
    style={{
      flex: 1,
      overflowY: 'auto',
      border: '1px solid #f0f0f0',
      borderRadius: 8,
      padding: 4,
      maxHeight: 'calc(100vh - 300px)',
    }}
  >
    <Row gutter={[0, 0]} justify="start">
      {productListProp.map((purchase) => (
        <Col key={purchase.product.id} xs={24} sm={12} md={12} lg={12} xl={12}>
          <Card
            size="small"
            title={purchase.product.name}
            style={{ boxShadow: '0 2px 8px rgba(0,0,0,0.1)', padding: 0, margin: '0 5px', textAlign: 'left' }}
            extra={
              <div>
                {onEdit && (
                  <Button
                    type="text"
                    icon={<EditTwoTone />}
                    onClick={() => onEdit(purchase)}
                    style={{ marginRight: 4 }}
                  />
                )}
                {onDelete && (
                  <Button
                    type="text"
                    icon={<DeleteTwoTone />}
                    onClick={() => onDelete(purchase.product.id)}
                  />
                )}
              </div>
            }
          >
            <p style={{ margin: 0 }}>
              {purchase.price}₽ x {purchase.quantity} = {purchase.price * purchase.quantity}₽
            </p>
            <div>
              <Tag color="green">{purchase.product.brand}</Tag>
              <Tag color="blue">{purchase.product.volume}</Tag>
              <Tag color="yellow">{purchase.product.category}</Tag>
              {purchase.tags && purchase.tags.map((tag, idx) => (
                <Tag key={idx}>{tag}</Tag>
              ))}
            </div>
          </Card>
        </Col>
      ))}
    </Row>
  </div>
);

export default ProductCardsWidget;
