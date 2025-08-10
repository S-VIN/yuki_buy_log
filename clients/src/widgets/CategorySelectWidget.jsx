/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';

const categories = ['Dairy', 'Bakery', 'Beverages'];

const CategorySelectWidget = ({ value, onChange }) => (
  <AutoComplete
    placeholder="category"
    options={categories.map((c) => ({ value: c }))}
    value={value}
    onChange={onChange}
  />
);

export default CategorySelectWidget;
