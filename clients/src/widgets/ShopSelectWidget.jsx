/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';
import { useMemo } from 'react';
import { useData } from '../stores/DataContext.jsx';

const ShopSelectWidget = ({ value, onChange }) => {
  const { purchases } = useData();

  const shops = useMemo(() => {
    const shopSet = new Set();
    purchases.forEach((p) => p.store && shopSet.add(p.store));
    return Array.from(shopSet).sort();
  }, [purchases]);

  return (
    <AutoComplete
      style={{ width: '100%' }}
      value={value}
      options={shops.map((s) => ({ value: s }))}
      onChange={onChange}
      placeholder="shop"
    />
  );
};

export default ShopSelectWidget;
