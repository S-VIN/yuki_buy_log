/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';
import { useMemo } from 'react';
import { useData } from '../stores/DataContext.jsx';

const BrandSelectWidget = ({ value, onChange }) => {
  const { products } = useData();

  const brands = useMemo(() => {
    const brandSet = new Set();
    products.forEach((p) => p.brand && brandSet.add(p.brand));
    return Array.from(brandSet).sort();
  }, [products]);

  return (
    <AutoComplete
      placeholder="brand"
      options={brands.map((b) => ({ value: b }))}
      value={value}
      onChange={onChange}
    />
  );
};

export default BrandSelectWidget;
