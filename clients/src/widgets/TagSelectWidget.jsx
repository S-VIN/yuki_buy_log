/* eslint-disable react/prop-types */
/* eslint-disable react/display-name */
import { forwardRef, useImperativeHandle, useState } from 'react';
import { Select } from 'antd';

const options = ['food', 'electronics', 'clothes'];

const TagSelectWidget = forwardRef(({ onTagChange, style }, ref) => {
  const [selectedTags, setSelectedTags] = useState([]);

  useImperativeHandle(ref, () => ({
    resetTags: () => setSelectedTags([]),
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

export default TagSelectWidget;
