/* eslint-disable react/prop-types */
/* eslint-disable react/display-name */
import { forwardRef, useImperativeHandle, useState } from 'react';
import { Input, Typography } from 'antd';
import { ShoppingCartOutlined } from '@ant-design/icons';

const { Text } = Typography;

const PriceQuantitySelectWidget = forwardRef(({ onPriceChanged, onQuantityChanged }, ref) => {
  const [price, setPrice] = useState(null);
  const [qty, setQty] = useState(1);

  useImperativeHandle(ref, () => ({
    reset: () => {
      setPrice(null);
      setQty(1);
    },
    setValues: (newPrice, newQty) => {
      setPrice(newPrice);
      setQty(newQty);
    },
  }));

  const handlePriceChange = (e) => {
    const value = parseInt(e.target.value, 10);
    if (!isNaN(value)) {
      setPrice(value);
      onPriceChanged(value);
    } else {
      setPrice(null);
    }
  };

  const handleQuantityChange = (e) => {
    const value = parseInt(e.target.value, 10);
    if (!isNaN(value)) {
      setQty(value);
      onQuantityChanged(value);
    } else {
      setQty(null);
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', flexDirection: 'row', gap: 8 }}>
        <Input
          type="text"
          inputMode="numeric"
          value={price ?? ''}
          onChange={handlePriceChange}
          placeholder="price"
          suffix={<label style={{ color: 'rgba(0,0,0,0.45)' }}>₽</label>}
        />
        <Input
          type="text"
          inputMode="numeric"
          value={qty ?? ''}
          onChange={handleQuantityChange}
          placeholder="count"
          suffix={<ShoppingCartOutlined style={{ color: 'rgba(0,0,0,0.45)' }} />}
        />
      </div>
      <Text type="secondary" style={{ height: 16, display: 'block', textAlign: 'left' }}>
        {!price || !qty ? ' ' : `total cost: ${price * qty}₽`}
      </Text>
    </div>
  );
});

export default PriceQuantitySelectWidget;
