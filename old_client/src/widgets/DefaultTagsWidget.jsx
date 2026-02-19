import PropTypes from 'prop-types';
import { Select } from 'antd';

const DefaultTagsWidget = ({ value, onChange, options = [], placeholder = 'default tags' }) => {
  const handleChange = (selectedTags) => {
    onChange(selectedTags);
  };

  return (
    <Select
      mode="tags"
      placeholder={placeholder}
      style={{ width: '100%' }}
      onChange={handleChange}
      tokenSeparators={[',']}
      value={value}
      options={options.map((item) => ({ value: item, label: item }))}
    />
  );
};

DefaultTagsWidget.propTypes = {
  value: PropTypes.arrayOf(PropTypes.string),
  onChange: PropTypes.func.isRequired,
  options: PropTypes.arrayOf(PropTypes.string),
  placeholder: PropTypes.string,
};

export default DefaultTagsWidget;