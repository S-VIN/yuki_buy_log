/* eslint-disable react/prop-types */
import { Select } from 'antd';

const DefaultTagsWidget = ({ value, onChange, options = [], placeholder = 'Default tags' }) => {
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

export default DefaultTagsWidget;