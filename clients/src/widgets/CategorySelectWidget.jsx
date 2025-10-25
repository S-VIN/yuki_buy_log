import PropTypes from 'prop-types';
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

CategorySelectWidget.propTypes = {
  value: PropTypes.string,
  onChange: PropTypes.func.isRequired,
};

export default CategorySelectWidget;
