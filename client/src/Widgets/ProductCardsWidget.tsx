import React, { useEffect, useState } from 'react';
import { Col, Row } from 'antd';
import CardWidget from "./CardWidget";
import { UUID } from 'crypto';

import Purchase from '../Models/Purchase';


interface ProductCardsWidgetProps {
  productListProp: Purchase[];
  onDelete: (product_id: UUID) => void;
}

const ProductCardsWidget: React.FC<ProductCardsWidgetProps> = ({ productListProp, onDelete }) => {
  const [productList, setProductList] = useState<Purchase[]>(productListProp);

  useEffect(() => {
    setProductList(productListProp);
  }, [productListProp]);

  const handleDelete = (product_id: UUID) => {
      console.log('product card widget handleDelete', product_id);
      onDelete(product_id);
  };

  return (
    <div
      style={{
        flex: 1,
        overflowY: 'auto',
        border: '1px solid #f0f0f0',
        borderRadius: '8px',
        padding: '4px',
        maxHeight: 'calc(100vh - 300px)', // Ограничиваем высоту, например, с учётом других элементов
      }}
    >
      <Row gutter={[0, 0]} justify="start">
        {productList.map((purchase, _) => (
          <Col key={purchase.product.id} xs={24} sm={12} md={12} lg={12} xl={12}>
            <CardWidget purchaseProp={purchase} onDelete={handleDelete} />
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default ProductCardsWidget;
