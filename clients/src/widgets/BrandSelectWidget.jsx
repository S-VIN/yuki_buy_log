/* eslint-disable react/prop-types */
import { AutoComplete } from 'antd';

const brands = ['BrandA', 'BrandB', 'BrandC'];

const BrandSelectWidget = ({ value, onChange }) => (
  <AutoComplete
    placeholder="brand"
    options={brands.map((b) => ({ value: b }))}
    value={value}
    onChange={onChange}
  />
);

export default BrandSelectWidget;
