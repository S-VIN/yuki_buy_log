/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';

const shops = ['Supermarket', 'ElectroMart', 'Cafe'];

const ShopSelectWidget = ({ value, onChange }) => {
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
