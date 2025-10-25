import { forwardRef, useImperativeHandle, useState } from 'react';
import PropTypes from 'prop-types';
import { Select } from 'antd';

const options = ['food', 'electronics', 'clothes'];

const TagSelectWidget = forwardRef(({ onTagChange, style }, ref) => {
  const [selectedTags, setSelectedTags] = useState([]);

  useImperativeHandle(ref, () => ({
    resetTags: () => setSelectedTags([]),
    setTags: (tags) => setSelectedTags(tags),
  }));

  const handleTagSelect = (value) => {
    setSelectedTags(value);
    onTagChange(value);
  };

  return (
    <Select
      mode="tags"
      placeholder="tags"
      style={{ width: '100%', ...style }}
      onChange={handleTagSelect}
      tokenSeparators={[',']}
      value={selectedTags}
      options={options.map((item) => ({ value: item, label: item }))}
    />
  );
});

TagSelectWidget.displayName = 'TagSelectWidget';

TagSelectWidget.propTypes = {
  onTagChange: PropTypes.func.isRequired,
  style: PropTypes.object,
};

export default TagSelectWidget;
